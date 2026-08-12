package payment

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"ticketBooking/internal/booking"
	bookingdto "ticketBooking/internal/booking/dto"
	"ticketBooking/internal/currency"
	"ticketBooking/internal/event"
	"ticketBooking/internal/payment/dto"
	"ticketBooking/internal/settings"

	"github.com/google/uuid"
	"github.com/stripe/stripe-go/v82"
	"github.com/stripe/stripe-go/v82/webhook"
)

var (
	ErrGatewayDisabled   = errors.New("selected payment gateway is not enabled")
	ErrInvalidMethod     = errors.New("invalid payment method")
	ErrSettingsNotFound  = errors.New("site settings not configured")
	ErrMissingGatewayKey = errors.New("payment gateway credentials missing in settings")
)

type Service interface {
	CreateCheckout(userID uint, userName, userEmail string, req dto.CreateCheckoutRequest) (*dto.CheckoutResponse, error)
	HandleStripeWebhook(payload []byte, signature string) error
	HandleSSLCommerzIPN(ipn dto.SSLCommerzIPN) error
	GetByID(id uint) (*dto.PaymentResponse, error)
	GetByTransactionID(tranID string) (*dto.PaymentResponse, error)
	GetByUserID(userID uint) ([]*dto.PaymentResponse, error)

	// Payment methods
	CreatePaymentMethod(req dto.CreatePaymentMethodRequest) (*dto.PaymentMethodResponse, error)
	UpdatePaymentMethod(id uint, req dto.UpdatePaymentMethodRequest) (*dto.PaymentMethodResponse, error)
	DeletePaymentMethod(id uint) error
	GetPaymentMethod(id uint) (*dto.PaymentMethodResponse, error)
	ListPaymentMethods(enabledOnly bool) ([]*dto.PaymentMethodResponse, error)
}

type service struct {
	paymentRepo   Repository
	bookingRepo   booking.Repository
	eventRepo     event.Repository
	settingsSvc   settings.Service
	currencySvc   currency.Service
	webhookSecret string
}

func NewService(pr Repository, br booking.Repository, er event.Repository, ss settings.Service, cs currency.Service) Service {
	return &service{paymentRepo: pr, bookingRepo: br, eventRepo: er, settingsSvc: ss, currencySvc: cs,}
}

func generateTransactionID() string {
	return "TXN_" + strings.ReplaceAll(uuid.New().String(), "-", "")
}

func (s *service) CreateCheckout(userID uint, userName, userEmail string, req dto.CreateCheckoutRequest) (*dto.CheckoutResponse, error) {
	method := PaymentMethodCode(strings.ToLower(strings.TrimSpace(req.PaymentMethod)))
	if method != MethodStripe && method != MethodSSLCommerz {
		return nil, ErrInvalidMethod
	}

	// Prefer payment_methods table enable flag when the row exists.
	if pm, err := s.paymentRepo.GetMethodByCode(string(method)); err == nil {
		if !pm.Enable {
			return nil, ErrGatewayDisabled
		}
	}

	setting, err := s.settingsSvc.GetRaw()
	if err != nil {
		return nil, ErrSettingsNotFound
	}
	eventData, err := s.eventRepo.GetByID(req.EventID)
	if err != nil {
		return nil, fmt.Errorf("event not found: %w", err)
	}
	currencyCode := s.currencySvc.GetDefaultCode()
	if currencyCode == "" {
		currencyCode = "BDT"
	}
	// if setting.CurrencyCode != "" {
	// 	currencyCode = setting.CurrencyCode
	// }
	if err := s.eventRepo.DecrementTickets(req.EventID, req.Quantity); err != nil {
		return nil, booking.ErrNotEnoughTickets
	}
	totalPrice := float64(req.Quantity) * float64(eventData.Price)
	reason := "Event ticket booking"
	if t := eventData.Title.Get("en"); t != "" {
		reason = fmt.Sprintf("Event ticket: %s (x%d)", t, req.Quantity)
	}
	b := &booking.Booking{
		UserID: userID, EventID: req.EventID, Quantity: req.Quantity,
		TotalPrice: totalPrice, Status: bookingdto.BookingPending,
		BookingCode: "GT-" + uuid.New().String(),
	}
	if err := s.bookingRepo.Create(b); err != nil {
		_ = s.eventRepo.IncrementTickets(req.EventID, req.Quantity)
		return nil, err
	}
	tranID := generateTransactionID()
	p := &Payment{
		BookingID: b.ID, UserID: userID, EventID: req.EventID,
		Amount: totalPrice, Currency: currencyCode, Reason: reason,
		PaymentMethod: method, Status: StatusPending, TransactionID: tranID,
	}
	if err := s.paymentRepo.Create(p); err != nil {
		_ = s.bookingRepo.UpdateStatus(b.ID, bookingdto.BookingCancelled)
		_ = s.eventRepo.IncrementTickets(req.EventID, req.Quantity)
		return nil, err
	}

		

	var checkoutURL, sessionID string
	switch method {
	case MethodStripe:
		// if !setting.StripeEnable {
		// 	s.rollback(b, req.EventID, req.Quantity)
		// 	return nil, ErrGatewayDisabled
		// }
		if setting.StripeSecretKey == "" {
			s.rollback(b, req.EventID, req.Quantity)
			return nil, ErrMissingGatewayKey
		}
		successURL := setting.StripeSuccessURL
		cancelURL := setting.StripeCancelURL
		if successURL == "" {
			successURL = setting.ClientSideURL + "/payment/success?session_id={CHECKOUT_SESSION_ID}"
		}
		if cancelURL == "" {
			cancelURL = setting.ClientSideURL + "/payment/cancel"
		}
		sc := newStripeClient(setting.StripeSecretKey)
		sess, err := sc.CreateCheckoutSession(stripeSessionParams{
			Amount: int64(eventData.Price * 100), Currency: strings.ToLower(currencyCode),
			ProductName: reason, Quantity: int64(req.Quantity),
			SuccessURL: successURL, CancelURL: cancelURL, ClientRef: tranID,
			CustomerEmail: req.CustomerEmail,
			Metadata: map[string]string{
				"payment_id":     strconv.FormatUint(uint64(p.ID), 10),
				"booking_id":     strconv.FormatUint(uint64(b.ID), 10),
				"transaction_id": tranID,
			},
		})
		if err != nil {
			s.rollback(b, req.EventID, req.Quantity)
			return nil, err
		}
		checkoutURL, sessionID = sess.URL, sess.ID
		p.GatewaySessionID, p.CheckoutURL = sess.ID, sess.URL

	case MethodSSLCommerz:
		// if !setting.SslCommerzeEnable {
		// 	s.rollback(b, req.EventID, req.Quantity)
		// 	return nil, ErrGatewayDisabled
		// }
		if setting.SslCommerzeStoreID == "" || setting.SslCommerzeStorePassword == "" {
			s.rollback(b, req.EventID, req.Quantity)
			return nil, ErrMissingGatewayKey
		}
		successURL := setting.SslCommerzeSuccessURL
		failURL := setting.SslCommerzeFailedURL
		cancelURL := setting.SslCommerzeCancelURL
		if successURL == "" {
			successURL = setting.ClientSideURL + "/payment/success"
		}
		if failURL == "" {
			failURL = setting.ClientSideURL + "/payment/failed"
		}
		if cancelURL == "" {
			cancelURL = setting.ClientSideURL + "/payment/cancel"
		}
		ipnURL := strings.TrimRight(setting.ServerSideURL, "/") + "/api/v1/payments/webhook/sslcommerz"
		client := newSSLCommerzClient(setting.SslCommerzeStoreID, setting.SslCommerzeStorePassword, true)
		initRes, err := client.Initiate(sslInitParams{
			TotalAmount: fmt.Sprintf("%.2f", totalPrice), Currency: currencyCode, TranID: tranID,
			SuccessURL: successURL, FailURL: failURL, CancelURL: cancelURL, IPNURL: ipnURL,
			CusName: req.CustomerName, CusEmail: req.CustomerEmail, CusPhone: req.CustomerPhoneNumber,
			CusAdd1: req.CustomerAddress, CusCity: req.Country, CusCountry: req.Country, CusPostcode: req.Postcode,
			ProductName: reason, ProductCategory: "ticket", ProductProfile: "general",
			ValueA: strconv.FormatUint(uint64(b.ID), 10),
			ValueB: strconv.FormatUint(uint64(p.ID), 10),
			ValueC: strconv.FormatUint(uint64(userID), 10),
			ValueD: strconv.FormatUint(uint64(req.EventID), 10),
		})
		if err != nil {
			s.rollback(b, req.EventID, req.Quantity)
			return nil, err
		}
		checkoutURL, sessionID = initRes.GatewayPageURL, initRes.SessionKey
		p.GatewaySessionID, p.CheckoutURL = initRes.SessionKey, checkoutURL
	}

	_ = s.paymentRepo.Update(p)
	return &dto.CheckoutResponse{
		PaymentID: p.ID, BookingID: b.ID, BookingCode: b.BookingCode,
		Amount: totalPrice, Currency: currencyCode, Reason: reason,
		PaymentMethod: string(method), TransactionID: tranID,
		CheckoutURL: checkoutURL, SessionID: sessionID, Status: string(StatusPending),
		CreatedAt: time.Now().Format("2006-01-02 15:04:05"),
	}, nil
}

func (s *service) rollback(b *booking.Booking, eventID uint, qty int) {
	_ = s.bookingRepo.UpdateStatus(b.ID, bookingdto.BookingCancelled)
	_ = s.eventRepo.IncrementTickets(eventID, qty)
}

// HandleStripeWebhook processes events forwarded by Stripe CLI / dashboard webhooks.
// Supported (local + production friendly):
//   - checkout.session.completed
//   - checkout.session.expired
//   - checkout.session.async_payment_failed
//   - checkout.session.async_payment_succeeded
//   - payment_intent.succeeded
//   - payment_intent.payment_failed
//   - charge.succeeded
//   - charge.failed
//   - charge.refunded
func (s *service) HandleStripeWebhook(payload []byte, signature string) error {
	var event stripe.Event
	var err error
	if s.webhookSecret != "" {
		event, err = webhook.ConstructEvent(payload, signature, s.webhookSecret)
		if err != nil {
			return err
		}
	} else if err := json.Unmarshal(payload, &event); err != nil {
		return err
	}

	switch event.Type {
	case "checkout.session.completed", "checkout.session.async_payment_succeeded":
		var sess stripe.CheckoutSession
		if err := json.Unmarshal(event.Data.Raw, &sess); err != nil {
			return err
		}
		return s.handleStripeSessionSuccess(&sess)

	case "checkout.session.expired":
		var sess stripe.CheckoutSession
		if err := json.Unmarshal(event.Data.Raw, &sess); err != nil {
			return err
		}
		return s.handleStripeSessionFail(&sess, StatusCancelled)

	case "checkout.session.async_payment_failed":
		var sess stripe.CheckoutSession
		if err := json.Unmarshal(event.Data.Raw, &sess); err != nil {
			return err
		}
		return s.handleStripeSessionFail(&sess, StatusFailed)

	case "payment_intent.succeeded":
		var pi stripe.PaymentIntent
		if err := json.Unmarshal(event.Data.Raw, &pi); err != nil {
			return err
		}
		if tranID := stripeMetaTranID(pi.Metadata); tranID != "" {
			return s.markSuccessByTranID(tranID)
		}
		// Fallback: metadata may carry payment_id
		if idStr, ok := pi.Metadata["payment_id"]; ok && idStr != "" {
			if id, e := strconv.ParseUint(idStr, 10, 64); e == nil {
				if p, e2 := s.paymentRepo.GetByID(uint(id)); e2 == nil {
					return s.markSuccessByTranID(p.TransactionID)
				}
			}
		}
		return nil

	case "payment_intent.payment_failed":
		var pi stripe.PaymentIntent
		if err := json.Unmarshal(event.Data.Raw, &pi); err != nil {
			return err
		}
		if tranID := stripeMetaTranID(pi.Metadata); tranID != "" {
			return s.markFailedByTranID(tranID, StatusFailed)
		}
		return nil

	case "charge.succeeded":
		var ch stripe.Charge
		if err := json.Unmarshal(event.Data.Raw, &ch); err != nil {
			return err
		}
		if tranID := stripeMetaTranID(ch.Metadata); tranID != "" {
			return s.markSuccessByTranID(tranID)
		}
		return nil

	case "charge.failed":
		var ch stripe.Charge
		if err := json.Unmarshal(event.Data.Raw, &ch); err != nil {
			return err
		}
		if tranID := stripeMetaTranID(ch.Metadata); tranID != "" {
			return s.markFailedByTranID(tranID, StatusFailed)
		}
		return nil

	case "charge.refunded":
		var ch stripe.Charge
		if err := json.Unmarshal(event.Data.Raw, &ch); err != nil {
			return err
		}
		// Treat full/partial refund as cancelled for booking purposes
		if tranID := stripeMetaTranID(ch.Metadata); tranID != "" {
			return s.markFailedByTranID(tranID, StatusCancelled)
		}
		return nil
	}
	// Unknown event types are acknowledged (200) so Stripe does not retry forever.
	return nil
}

func stripeMetaTranID(meta map[string]string) string {
	if meta == nil {
		return ""
	}
	if v := meta["transaction_id"]; v != "" {
		return v
	}
	return ""
}

func (s *service) handleStripeSessionSuccess(sess *stripe.CheckoutSession) error {
	tranID := sess.ClientReferenceID
	if tranID == "" && sess.Metadata != nil {
		tranID = sess.Metadata["transaction_id"]
	}
	if tranID != "" {
		return s.markSuccessByTranID(tranID)
	}
	return s.markSuccessBySession(sess.ID)
}

func (s *service) handleStripeSessionFail(sess *stripe.CheckoutSession, st PaymentStatus) error {
	tranID := sess.ClientReferenceID
	if tranID == "" && sess.Metadata != nil {
		tranID = sess.Metadata["transaction_id"]
	}
	if tranID != "" {
		return s.markFailedByTranID(tranID, st)
	}
	return s.markFailedBySession(sess.ID, st)
}

func (s *service) markSuccessByTranID(tranID string) error {
	p, err := s.paymentRepo.GetByTransactionID(tranID)
	if err != nil {
		return err
	}
	if p.Status == StatusSuccess {
		return nil
	}
	now := time.Now()
	p.Status = StatusSuccess
	p.PaidAt = &now
	if err := s.paymentRepo.Update(p); err != nil {
		return err
	}
	return s.bookingRepo.UpdateStatus(p.BookingID, bookingdto.BookingConfirmed)
}

func (s *service) markSuccessBySession(sessionID string) error {
	p, err := s.paymentRepo.GetBySessionID(sessionID)
	if err != nil {
		return err
	}
	return s.markSuccessByTranID(p.TransactionID)
}

func (s *service) markFailedByTranID(tranID string, st PaymentStatus) error {
	p, err := s.paymentRepo.GetByTransactionID(tranID)
	if err != nil {
		return err
	}
	if p.Status == StatusSuccess {
		return nil
	}
	p.Status = st
	_ = s.paymentRepo.Update(p)
	_ = s.bookingRepo.UpdateStatus(p.BookingID, bookingdto.BookingCancelled)
	if b, e := s.bookingRepo.GetByID(p.BookingID); e == nil {
		_ = s.eventRepo.IncrementTickets(b.EventID, b.Quantity)
	}
	return nil
}

func (s *service) markFailedBySession(sessionID string, st PaymentStatus) error {
	p, err := s.paymentRepo.GetBySessionID(sessionID)
	if err != nil {
		return err
	}
	return s.markFailedByTranID(p.TransactionID, st)
}

func (s *service) HandleSSLCommerzIPN(ipn dto.SSLCommerzIPN) error {
	p, err := s.paymentRepo.GetByTransactionID(ipn.TranID)
	if err != nil {
		return err
	}
	if p.Status == StatusSuccess {
		return nil
	}
	setting, err := s.settingsSvc.GetRaw()
	if err != nil {
		return err
	}
	status := strings.ToUpper(strings.TrimSpace(ipn.Status))
	if ipn.ValID != "" && (status == "VALID" || status == "VALIDATED") {
		client := newSSLCommerzClient(setting.SslCommerzeStoreID, setting.SslCommerzeStorePassword, true)
		if val, e := client.Validate(ipn.ValID); e == nil && val != nil {
			status = strings.ToUpper(val.Status)
		}
	}
	raw, _ := json.Marshal(ipn)
	p.GatewayResponse = string(raw)
	switch status {
	case "VALID", "VALIDATED":
		now := time.Now()
		p.Status = StatusSuccess
		p.PaidAt = &now
		_ = s.paymentRepo.Update(p)
		return s.bookingRepo.UpdateStatus(p.BookingID, bookingdto.BookingConfirmed)
	case "FAILED":
		return s.markFailedByTranID(ipn.TranID, StatusFailed)
	case "CANCELLED", "CANCELED":
		return s.markFailedByTranID(ipn.TranID, StatusCancelled)
	}
	return s.paymentRepo.Update(p)
}

func (s *service) GetByID(id uint) (*dto.PaymentResponse, error) {
	p, err := s.paymentRepo.GetByID(id)
	if err != nil {
		return nil, err
	}
	code := ""
	if b, e := s.bookingRepo.GetByID(p.BookingID); e == nil {
		code = b.BookingCode
	}
	return p.ToResponse(code), nil
}

func (s *service) GetByTransactionID(tranID string) (*dto.PaymentResponse, error) {
	p, err := s.paymentRepo.GetByTransactionID(tranID)
	if err != nil {
		return nil, err
	}
	code := ""
	if b, e := s.bookingRepo.GetByID(p.BookingID); e == nil {
		code = b.BookingCode
	}
	return p.ToResponse(code), nil
}

func (s *service) GetByUserID(userID uint) ([]*dto.PaymentResponse, error) {
	list, err := s.paymentRepo.GetByUserID(userID)
	if err != nil {
		return nil, err
	}
	out := make([]*dto.PaymentResponse, 0, len(list))
	for i := range list {
		code := ""
		if b, e := s.bookingRepo.GetByID(list[i].BookingID); e == nil {
			code = b.BookingCode
		}
		out = append(out, list[i].ToResponse(code))
	}
	return out, nil
}












// ---- Payment Method CRUD ----

func (s *service) CreatePaymentMethod(req dto.CreatePaymentMethodRequest) (*dto.PaymentMethodResponse, error) {
	code := strings.ToLower(strings.TrimSpace(req.Code))
	if code != string(MethodStripe) && code != string(MethodSSLCommerz) {
		return nil, ErrInvalidMethod
	}
	if _, err := s.paymentRepo.GetMethodByCode(code); err == nil {
		return nil, ErrPaymentMethodExists
	}
	enable := true
	if req.Enable != nil {
		enable = *req.Enable
	}
	m := &PaymentMethod{
		Name:   strings.TrimSpace(req.Name),
		Code:   code,
		LogoURL:   strings.TrimSpace(req.LogoURL),
		LogoID:    req.LogoID,
		Enable: enable,
	}
	if err := s.paymentRepo.CreateMethod(m); err != nil {
		return nil, err
	}
	return m.ToResponse(), nil
}

func (s *service) UpdatePaymentMethod(id uint, req dto.UpdatePaymentMethodRequest) (*dto.PaymentMethodResponse, error) {
	m, err := s.paymentRepo.GetMethodByID(id)
	if err != nil {
		return nil, err
	}
	if req.Name != nil {
		m.Name = strings.TrimSpace(*req.Name)
	}
	if req.Code != nil {
		code := strings.ToLower(strings.TrimSpace(*req.Code))
		if code != string(MethodStripe) && code != string(MethodSSLCommerz) {
			return nil, ErrInvalidMethod
		}
		// uniqueness check if code changes
		if code != m.Code {
			if existing, e := s.paymentRepo.GetMethodByCode(code); e == nil && existing.ID != m.ID {
				return nil, ErrPaymentMethodExists
			}
			m.Code = code
		}
	}
	if req.LogoURL != nil {
		m.LogoURL = strings.TrimSpace(*req.LogoURL)
	}
	if req.Enable != nil {
		m.Enable = *req.Enable
	}
	if err := s.paymentRepo.UpdateMethod(m); err != nil {
		return nil, err
	}
	return m.ToResponse(), nil
}

func (s *service) DeletePaymentMethod(id uint) error {
	return s.paymentRepo.DeleteMethod(id)
}

func (s *service) GetPaymentMethod(id uint) (*dto.PaymentMethodResponse, error) {
	m, err := s.paymentRepo.GetMethodByID(id)
	if err != nil {
		return nil, err
	}
	return m.ToResponse(), nil
}

func (s *service) ListPaymentMethods(enabledOnly bool) ([]*dto.PaymentMethodResponse, error) {
	list, err := s.paymentRepo.ListMethods(enabledOnly)
	if err != nil {
		return nil, err
	}
	out := make([]*dto.PaymentMethodResponse, 0, len(list))
	for i := range list {
		out = append(out, list[i].ToResponse())
	}
	return out, nil
}
