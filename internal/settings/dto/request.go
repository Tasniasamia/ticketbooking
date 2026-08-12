package dto

// UpsertRequest is used for both create and update (singleton settings).
// All fields are optional so partial updates work.
type UpsertRequest struct {
	SiteName        *string `json:"site_name" validate:"omitempty,min=1,max=255"`
	SiteEmail       *string `json:"site_email" validate:"omitempty,email"`
	SitePhone       *string `json:"site_phone" validate:"omitempty,max=50"`
	SiteLogo        *string `json:"site_logo"`
	SiteAddress     *string `json:"site_address"`
	SiteDescription *string `json:"site_description"`
	SiteFooter      *string `json:"site_footer"`

	CurrencyCode   *string `json:"currency_code" validate:"omitempty,min=2,max=10"`
	CurrencySymbol *string `json:"currency_symbol" validate:"omitempty,max=10"`

	ClientSideURL *string `json:"client_side_url"`
	ServerSideURL *string `json:"server_side_url"`

	SocialMediaLink *[]SocialMediaLink `json:"social_media_link"`
	Partner         *[]string          `json:"partner"`

	// SSLCommerz
	SslCommerzeStoreID       *string `json:"sslCommerze_store_id"`
	SslCommerzeStorePassword *string `json:"sslCommerze_store_password"`
	SslCommerzeSuccessURL    *string `json:"sslCommerze_success_url"`
	SslCommerzeFailedURL     *string `json:"sslCommerze_failed_url"`
	SslCommerzeCancelURL     *string `json:"sslCommerze_cancel_url"`
	SslCommerzeEnable        *bool   `json:"sslCommerze_enable"`

	// Stripe
	StripePublishableKey *string `json:"stripe_publishable_key"`
	StripeSecretKey      *string `json:"stripe_secret_key"`
	StripeSuccessURL     *string `json:"stripe_success_url"`
	StripeFailedURL      *string `json:"stripe_failed_url"`
	StripeCancelURL      *string `json:"stripe_cancel_url"`
	StripeEnable         *bool   `json:"stripe_enable"`
}