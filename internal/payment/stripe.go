package payment

import (
	"fmt"

	"github.com/stripe/stripe-go/v82"
	"github.com/stripe/stripe-go/v82/checkout/session"
)

type stripeClient struct{ secretKey string }

func newStripeClient(secretKey string) *stripeClient { return &stripeClient{secretKey: secretKey} }

type stripeSessionParams struct {
	Amount                                   int64
	Currency, ProductName, SuccessURL, CancelURL, ClientRef string
	Quantity                                 int64
	CustomerEmail                            string
	Metadata                                 map[string]string
}

func (sc *stripeClient) CreateCheckoutSession(p stripeSessionParams) (*stripe.CheckoutSession, error) {
	stripe.Key = sc.secretKey
	params := &stripe.CheckoutSessionParams{
		Mode: stripe.String(string(stripe.CheckoutSessionModePayment)),
		LineItems: []*stripe.CheckoutSessionLineItemParams{{
			PriceData: &stripe.CheckoutSessionLineItemPriceDataParams{
				Currency: stripe.String(p.Currency),
				ProductData: &stripe.CheckoutSessionLineItemPriceDataProductDataParams{
					Name: stripe.String(p.ProductName),
				},
				UnitAmount: stripe.Int64(p.Amount),
			},
			Quantity: stripe.Int64(p.Quantity),
		}},
		SuccessURL:        stripe.String(p.SuccessURL),
		CancelURL:         stripe.String(p.CancelURL),
		ClientReferenceID: stripe.String(p.ClientRef),
	}
	if p.CustomerEmail != "" {
		params.CustomerEmail = stripe.String(p.CustomerEmail)
	}
	if len(p.Metadata) > 0 {
		params.Metadata = p.Metadata
	}
	sess, err := session.New(params)
	if err != nil {
		return nil, fmt.Errorf("stripe session create: %w", err)
	}
	return sess, nil
}
