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
	"ticketBooking/internal/httpresponse"
	"ticketBooking/internal/payment/dto"

	"github.com/google/uuid"
	"github.com/stripe/stripe-go/v82"
	"github.com/stripe/stripe-go/v82/checkout/session"
	"github.com/stripe/stripe-go/v82/webhook"
	"ticketBooking/internal/utils/query"
)

var (
	ErrGatewayDisabled   = errors.New("selected payment gateway is not enabled")
	ErrInvalidMethod     = errors.New("invalid payment method")
	ErrSettingsNotFound  = errors.New("site settings not configured")
	ErrMissingGatewayKey = errors.New("payment gateway credentials missing")
)

// PaymentsListData — list + pagination + payment aggregation (default currency).
type PaymentsListData struct {
	*httpresponse.PaginatedData
	Payment dto.PaymentSummary `json:"payment"`
}

type Service interface {
	CreateCheckout(userID uint, userName, userEmail string, req dto.CreateCheckoutRequest) (*dto.CheckoutResponse, error)
	HandleStripeWebhook(payload []byte, signature string) error
	HandleSSLCommerzIPN(ipn dto.SSLCommerzIPN) error
	GetByID(id uint) (*dto.PaymentResponse, error)
	GetByTransactionID(tranID string) (*dto.PaymentResponse, error)
	GetByUserID(userID uint) ([]*dto.PaymentResponse, error)
	VerifyStripeSession(sessionID string) (*dto.PaymentResponse, error)
	VerifySSLCommerzSession(tranID, status string) (*dto.PaymentResponse, error)
	GetUserPayments(params query.Params, userID uint,lang string) (*PaymentsListData, error)
	GetManagerPayments(params query.Params, managerID uint,lang string) (*PaymentsListData, error)
	GetAllPayments(params query.Params,lang string) (*PaymentsListData, error)
	// Payment methods
	CreatePaymentMethod(req dto.CreatePaymentMethodRequest) (*dto.PaymentMethodResponse, error)
	UpdatePaymentMethod(id uint, req dto.UpdatePaymentMethodRequest) (*dto.PaymentMethodResponse, error)
	DeletePaymentMethod(id uint) error
	GetPaymentMethod(id uint) (*dto.PaymentMethodResponse, error)
ListPaymentMethods(p query.Params, enabledOnly bool) (*httpresponse.PaginatedData, error)
ListPaymentMethodsAdmin(p query.Params) (*httpresponse.PaginatedData, error)
}

type service struct {
	paymentRepo Repository
	bookingRepo booking.Repository
	eventRepo   event.Repository
	currencySvc currency.Service
}

func NewService(pr Repository, br booking.Repository, er event.Repository, cs currency.Service) Service {
	return &service{paymentRepo: pr, bookingRepo: br, eventRepo: er, currencySvc: cs}
}

func generateTransactionID() string {
	return "TXN_" + strings.ReplaceAll(uuid.New().String(), "-", "")
}

// getDecryptedCredentials loads payment method by code, checks enable flag,
// decrypts credentials and returns both the map and the method row.
func (s *service) getDecryptedCredentials(code string) (map[string]string, *PaymentMethod, error) {
	pm, err := s.paymentRepo.GetMethodByCode(code)
	if err != nil {
		return nil, nil, ErrGatewayDisabled
	}
	if !pm.Enable {
		return nil, nil, ErrGatewayDisabled
	}

	creds, err := DecryptCredentials(string(pm.Credentials))
	if err != nil {
		return nil, nil, fmt.Errorf("failed to decrypt credentials: %w", err)
	}
	return creds, &pm, nil
}

func (s *service) CreateCheckout(userID uint, userName, userEmail string, req dto.CreateCheckoutRequest) (*dto.CheckoutResponse, error) {
	method := PaymentMethodCode(strings.ToLower(strings.TrimSpace(req.PaymentMethod)))
	if method != MethodStripe && method != MethodSSLCommerz {
		return nil, ErrInvalidMethod
	}

	// Load + decrypt credentials (also checks enable)
	creds, _, err := s.getDecryptedCredentials(string(method))
	if err != nil {
		return nil, err
	}

	eventData, err := s.eventRepo.GetByID(req.EventID)
	if err != nil {
		return nil, fmt.Errorf("event not found: %w", err)
	}

	// Payment currency = event's price currency; fallback to site default
	currencyCode := strings.ToUpper(strings.TrimSpace(eventData.Currency))
	if currencyCode == "" {
		currencyCode = s.currencySvc.GetDefaultCode()
	}
	if currencyCode == "" {
		currencyCode = "BDT"
	}

	if err := s.eventRepo.DecrementTickets(req.EventID, req.Quantity); err != nil {
		return nil, booking.ErrNotEnoughTickets
	}

	totalPrice := float64(req.Quantity) * float64(eventData.Price)
	reason := "Event ticket booking"
	if t := eventData.Title.Get("en"); t != "" {
		reason = fmt.Sprintf("Event ticket: %s (x%d)", t, req.Quantity)
	}

	b := &booking.Booking{
		UserID:      userID,
		EventID:     req.EventID,
		Quantity:    req.Quantity,
		TotalPrice:  totalPrice,
		Status:      bookingdto.BookingPending,
		BookingCode: "GT-" + uuid.New().String(),
	}
	if err := s.bookingRepo.Create(b); err != nil {
		_ = s.eventRepo.IncrementTickets(req.EventID, req.Quantity)
		return nil, err
	}

	tranID := generateTransactionID()
	p := &Payment{
		BookingID:     b.ID,
		UserID:        userID,
		EventID:       req.EventID,
		Amount:        totalPrice,
		Currency:      currencyCode,
		Reason:        reason,
		PaymentMethod: method,
		Status:        StatusPending,
		TransactionID: tranID,
	}
	if err := s.paymentRepo.Create(p); err != nil {
		_ = s.bookingRepo.UpdateStatus(b.ID, bookingdto.BookingCancelled)
		_ = s.eventRepo.IncrementTickets(req.EventID, req.Quantity)
		return nil, err
	}

	var checkoutURL, sessionID string

	switch method {
	case MethodStripe:
		secretKey := creds["stripe_secret_key"]
		successURL := creds["stripe_success_url"]
		cancelURL := creds["stripe_cancel_url"]
		clientSideURL := creds["client_side_url"]

		if secretKey == "" {
			s.rollback(b, req.EventID, req.Quantity)
			return nil, ErrMissingGatewayKey
		}

		if successURL == "" {
			successURL = clientSideURL + "/payment/success"
		}
		if !strings.Contains(successURL, "{CHECKOUT_SESSION_ID}") {
			if strings.Contains(successURL, "?") {
				successURL += "&session_id={CHECKOUT_SESSION_ID}"
			} else {
				successURL += "?session_id={CHECKOUT_SESSION_ID}"
			}
		}
		if cancelURL == "" {
			cancelURL = clientSideURL + "/payment/cancel"
		}

		sc := newStripeClient(secretKey)
		sess, err := sc.CreateCheckoutSession(stripeSessionParams{
			Amount:        int64(eventData.Price * 100),
			Currency:      strings.ToLower(currencyCode),
			ProductName:   reason,
			Quantity:      int64(req.Quantity),
			SuccessURL:    successURL,
			CancelURL:     cancelURL,
			ClientRef:     tranID,
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
		if strings.ToUpper(currencyCode) != "BDT" {
			s.rollback(b, req.EventID, req.Quantity)
			converted, err := s.currencySvc.Convert(totalPrice, strings.ToUpper(currencyCode), "BDT")
			if err != nil {
				return nil, errors.New("Failed to  BDT currency")
			}
			totalPrice = converted
			currencyCode = "BDT"

		}

		storeID := creds["sslcommerz_store_id"]
		storePassword := creds["sslcommerz_store_password"]
		successURL := creds["sslcommerz_success_url"]
		failURL := creds["sslcommerz_failed_url"]
		cancelURL := creds["sslcommerz_cancel_url"]
		clientSideURL := creds["client_side_url"]

		if storeID == "" || storePassword == "" {
			s.rollback(b, req.EventID, req.Quantity)
			return nil, ErrMissingGatewayKey
		}

		if successURL == "" {
			successURL = clientSideURL + "/payment/success"
		}
		if failURL == "" {
			failURL = clientSideURL + "/payment/failed"
		}
		if cancelURL == "" {
			cancelURL = clientSideURL + "/payment/cancel"
		}

		// append tran_id if missing
		if !strings.Contains(successURL, "tran_id=") {
			if strings.Contains(successURL, "?") {
				successURL += "&tran_id=" + tranID
			} else {
				successURL += "?tran_id=" + tranID
			}
		}
		if !strings.Contains(failURL, "tran_id=") {
			if strings.Contains(failURL, "?") {
				failURL += "&tran_id=" + tranID
			} else {
				failURL += "?tran_id=" + tranID
			}
		}
		if !strings.Contains(cancelURL, "tran_id=") {
			if strings.Contains(cancelURL, "?") {
				cancelURL += "&tran_id=" + tranID
			} else {
				cancelURL += "?tran_id=" + tranID
			}
		}

		ipnURL := strings.TrimRight(clientSideURL, "/") + "/api/v1/payments/webhook/sslcommerz"
		client := newSSLCommerzClient(storeID, storePassword, true)

		initRes, err := client.Initiate(sslInitParams{
			TotalAmount:     fmt.Sprintf("%.2f", totalPrice),
			Currency:        currencyCode,
			TranID:          tranID,
			SuccessURL:      successURL,
			FailURL:         failURL,
			CancelURL:       cancelURL,
			IPNURL:          ipnURL,
			CusName:         req.CustomerName,
			CusEmail:        req.CustomerEmail,
			CusPhone:        req.CustomerPhoneNumber,
			CusAdd1:         req.CustomerAddress,
			CusCity:         req.Country,
			CusCountry:      req.Country,
			CusPostcode:     req.Postcode,
			ProductName:     reason,
			ProductCategory: "ticket",
			ProductProfile:  "general",
			ValueA:          strconv.FormatUint(uint64(b.ID), 10),
			ValueB:          strconv.FormatUint(uint64(p.ID), 10),
			ValueC:          strconv.FormatUint(uint64(userID), 10),
			ValueD:          strconv.FormatUint(uint64(req.EventID), 10),
		})
		if err != nil {
			fmt.Printf("SSLCommerz Initiate Error: %+v\n", err)
			s.rollback(b, req.EventID, req.Quantity)
			return nil, err
		}

		checkoutURL, sessionID = initRes.GatewayPageURL, initRes.SessionKey
		p.GatewaySessionID, p.CheckoutURL = initRes.SessionKey, checkoutURL
	}

	_ = s.paymentRepo.Update(p)

	return &dto.CheckoutResponse{
		PaymentID:     p.ID,
		BookingID:     b.ID,
		BookingCode:   b.BookingCode,
		Amount:        totalPrice,
		Currency:      currencyCode,
		Reason:        reason,
		PaymentMethod: string(method),
		TransactionID: tranID,
		CheckoutURL:   checkoutURL,
		SessionID:     sessionID,
		Status:        string(StatusPending),
		CreatedAt:     time.Now().Format("2006-01-02 15:04:05"),
	}, nil
}

func (s *service) rollback(b *booking.Booking, eventID uint, qty int) {
	_ = s.bookingRepo.UpdateStatus(b.ID, bookingdto.BookingCancelled)
	_ = s.eventRepo.IncrementTickets(eventID, qty)
}

func (s *service) HandleStripeWebhook(payload []byte, signature string) error {
	var event stripe.Event
	var err error

	creds, _, err := s.getDecryptedCredentials("stripe")
	if err != nil {
		return err
	}

	webhookSecret := creds["stripe_webhook_secret"]

	if webhookSecret != "" {
		event, err = webhook.ConstructEventWithOptions(
			payload,
			signature,
			webhookSecret,
			webhook.ConstructEventOptions{
				IgnoreAPIVersionMismatch: true,
			},
		)
		if err != nil {
			return fmt.Errorf("webhook signature verification failed: %w", err)
		}
	} else {
		if err := json.Unmarshal(payload, &event); err != nil {
			return err
		}
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
		if tranID := stripeMetaTranID(ch.Metadata); tranID != "" {
			return s.markFailedByTranID(tranID, StatusCancelled)
		}
		return nil
	}

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

	creds, _, err := s.getDecryptedCredentials("sslcommerz")
	if err != nil {
		return err
	}

	status := strings.ToUpper(strings.TrimSpace(ipn.Status))
	if ipn.ValID != "" && (status == "VALID" || status == "VALIDATED") {
		client := newSSLCommerzClient(creds["sslcommerz_store_id"], creds["sslcommerz_store_password"], true)
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

func (s *service) VerifyStripeSession(sessionID string) (*dto.PaymentResponse, error) {
	creds, _, err := s.getDecryptedCredentials("stripe")
	if err != nil {
		return nil, err
	}

	secretKey := creds["stripe_secret_key"]
	if secretKey == "" {
		return nil, ErrMissingGatewayKey
	}

	stripe.Key = secretKey
	sess, err := session.Get(sessionID, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch stripe session: %w", err)
	}

	if sess.PaymentStatus == stripe.CheckoutSessionPaymentStatusPaid ||
		sess.Status == stripe.CheckoutSessionStatusComplete {
		if err := s.handleStripeSessionSuccess(sess); err != nil {
			return nil, err
		}
	}

	p, err := s.paymentRepo.GetBySessionID(sessionID)
	if err != nil {
		if sess.ClientReferenceID != "" {
			p, err = s.paymentRepo.GetByTransactionID(sess.ClientReferenceID)
		}
		if err != nil {
			return nil, err
		}
	}

	code := ""
	if b, e := s.bookingRepo.GetByID(p.BookingID); e == nil {
		code = b.BookingCode
	}
	return p.ToResponse(code), nil
}

func (s *service) VerifySSLCommerzSession(tranID, status string) (*dto.PaymentResponse, error) {
	p, err := s.paymentRepo.GetByTransactionID(tranID)
	if err != nil {
		return nil, err
	}

	if p.Status == StatusSuccess {
		code := ""
		if b, e := s.bookingRepo.GetByID(p.BookingID); e == nil {
			code = b.BookingCode
		}
		return p.ToResponse(code), nil
	}

	finalStatus := strings.ToUpper(strings.TrimSpace(status))

	switch finalStatus {
	case "SUCCESS":
		now := time.Now()
		p.Status = StatusSuccess
		p.PaidAt = &now
		_ = s.paymentRepo.Update(p)
		_ = s.bookingRepo.UpdateStatus(p.BookingID, bookingdto.BookingConfirmed)
	case "FAILED":
		_ = s.markFailedByTranID(tranID, StatusFailed)
	case "CANCELLED", "CANCELED":
		_ = s.markFailedByTranID(tranID, StatusCancelled)
	}

	p, _ = s.paymentRepo.GetByTransactionID(tranID)
	code := ""
	if b, e := s.bookingRepo.GetByID(p.BookingID); e == nil {
		code = b.BookingCode
	}
	return p.ToResponse(code), nil
}

func (s *service) GetAllPayments(params query.Params,lang string) (*PaymentsListData, error) {
	events, total, err := s.paymentRepo.GetAllPayments(params,lang)
	if err != nil {
		return nil, err
	}

	docs := make([]*dto.PaymentResponse, 0, len(events))
	for _, e := range events {
		docs = append(docs, e.ToRawResponse())
	}

	rows, err := s.paymentRepo.AggregateAllPayments(params)
	if err != nil {
		return nil, err
	}
	summary := s.buildPaymentSummary(rows)

	meta := httpresponse.BuildPaginationMeta(total, params.Page, params.Limit)
	return &PaymentsListData{
		PaginatedData: &httpresponse.PaginatedData{
			Docs:           docs,
			PaginationMeta: meta,
		},
		Payment: summary,
	}, nil
}

func (s *service) GetManagerPayments(params query.Params, managerID uint, lang string) (*PaymentsListData, error) {
	events, total, err := s.paymentRepo.GetManagerPayments(params, managerID, lang)
	if err != nil {
		return nil, err
	}

	docs := make([]*dto.PaymentResponse, 0, len(events))
	for _, e := range events {
		docs = append(docs, e.ToRawResponse())
	}

	rows, err := s.paymentRepo.AggregateManagerPayments(params, managerID)
	if err != nil {
		return nil, err
	}
	summary := s.buildPaymentSummary(rows)

	meta := httpresponse.BuildPaginationMeta(total, params.Page, params.Limit)
	return &PaymentsListData{
		PaginatedData: &httpresponse.PaginatedData{
			Docs:           docs,
			PaginationMeta: meta,
		},
		Payment: summary,
	}, nil
}

func (s *service) GetUserPayments(params query.Params, userID uint, lang string) (*PaymentsListData, error) {
	events, total, err := s.paymentRepo.GetUserPayments(params, userID, lang)
	if err != nil {
		return nil, err
	}

	docs := make([]*dto.PaymentResponse, 0, len(events))
	for _, e := range events {
		docs = append(docs, e.ToRawResponse())
	}

	rows, err := s.paymentRepo.AggregateUserPayments(params, userID)
	if err != nil {
		return nil, err
	}
	summary := s.buildPaymentSummary(rows)

	meta := httpresponse.BuildPaginationMeta(total, params.Page, params.Limit)
	return &PaymentsListData{
		PaginatedData: &httpresponse.PaginatedData{
			Docs:           docs,
			PaginationMeta: meta,
		},
		Payment: summary,
	}, nil
}

// toDefaultCurrency converts amount from `from` into site default currency via currency.Service.Convert.
func (s *service) toDefaultCurrency(amount float64, from string) float64 {
	def := s.currencySvc.GetDefaultCode()
	if def == "" {
		def = "BDT"
	}
	from = strings.ToUpper(strings.TrimSpace(from))
	if from == "" {
		from = def
	}
	if strings.EqualFold(from, def) {
		return amount
	}
	converted, err := s.currencySvc.Convert(amount, from, def)
	if err != nil {
		return amount
	}
	return converted
}

func (s *service) buildPaymentSummary(rows []PaymentAmountRow) dto.PaymentSummary {
	def := s.currencySvc.GetDefaultCode()
	if def == "" {
		def = "BDT"
	}

	var pending, paid float64
	for _, row := range rows {
		amt := s.toDefaultCurrency(row.Total, row.Currency)
		switch PaymentStatus(strings.ToLower(strings.TrimSpace(row.Status))) {
		case StatusPending:
			pending += amt
		case StatusSuccess:
			paid += amt
		}
	}

	return dto.PaymentSummary{
		Pending:  pending,
		Paid:     paid,
		Total:    pending + paid,
		Currency: def,
	}
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

	encrypted, err := EncryptCredentials(req.Credentials)
	if err != nil {
		return nil, fmt.Errorf("failed to encrypt credentials: %w", err)
	}

	m := &PaymentMethod{
		Name:        strings.TrimSpace(req.Name),
		Code:        code,
		LogoURL:     strings.TrimSpace(req.LogoURL),
		LogoID:      req.LogoID,
		Enable:      enable,
		Credentials: string(encrypted),
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
	if req.LogoID != nil {
		m.LogoID = *req.LogoID
	}
	if req.Enable != nil {
		m.Enable = *req.Enable
	}
	if req.Credentials != nil {
		encrypted, err := EncryptCredentials(*req.Credentials)
		if err != nil {
			return nil, fmt.Errorf("failed to encrypt credentials: %w", err)
		}
		m.Credentials = string(encrypted)
	}

	if err := s.paymentRepo.UpdateMethod(&m); err != nil {
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

func (s *service) ListPaymentMethods(p query.Params, enabledOnly bool) (*httpresponse.PaginatedData, error) {

	list, total, err := s.paymentRepo.ListMethods(p,enabledOnly)
	if err != nil {
		return nil, err
	}

	docs := make([]*dto.PaymentMethodResponse, 0, len(list))
	for _, e := range list {
		docs = append(docs, e.ToResponse())
	}

	meta := httpresponse.BuildPaginationMeta(total, p.Page, p.Limit)
	return &httpresponse.PaginatedData{
		Docs:           docs,
		PaginationMeta: meta,
	}, nil
}

func (s *service) ListPaymentMethodsAdmin(p query.Params) (*httpresponse.PaginatedData, error) {

	list, total, err := s.paymentRepo.ListMethodsAdmin(p)
	if err != nil {
		return nil, err
	}

	docs := make([]*dto.PaymentMethodResponse, 0, len(list))
	for _, e := range list {
		docs = append(docs, e.ToResponse())
	}

	meta := httpresponse.BuildPaginationMeta(total, p.Page, p.Limit)
	return &httpresponse.PaginatedData{
		Docs:           docs,
		PaginationMeta: meta,
	}, nil
}