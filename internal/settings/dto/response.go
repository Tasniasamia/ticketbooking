package dto

type SocialMediaLink struct {
	Name string `json:"name"`
	Link string `json:"link"`
}

type Response struct {
	ID        uint   `json:"id"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`

	SiteName        string `json:"site_name"`
	SiteEmail       string `json:"site_email"`
	SitePhone       string `json:"site_phone"`
	SiteLogo        string `json:"site_logo"`
	SiteAddress     string `json:"site_address"`
	SiteDescription string `json:"site_description"`
	SiteFooter      string `json:"site_footer"`

	CurrencyCode   string `json:"currency_code"`
	CurrencySymbol string `json:"currency_symbol"`

	ClientSideURL string `json:"client_side_url"`
	ServerSideURL string `json:"server_side_url"`

	SocialMediaLink []SocialMediaLink `json:"social_media_link"`
	Partner         []string          `json:"partner"`

	SslCommerzeStoreID       string `json:"sslCommerze_store_id"`
	SslCommerzeStorePassword string `json:"sslCommerze_store_password"`
	SslCommerzeSuccessURL    string `json:"sslCommerze_success_url"`
	SslCommerzeFailedURL     string `json:"sslCommerze_failed_url"`
	SslCommerzeCancelURL     string `json:"sslCommerze_cancel_url"`
	SslCommerzeEnable        bool   `json:"sslCommerze_enable"`

	StripePublishableKey string `json:"stripe_publishable_key"`
	StripeSecretKey      string `json:"stripe_secret_key"`
	StripeSuccessURL     string `json:"stripe_success_url"`
	StripeFailedURL      string `json:"stripe_failed_url"`
	StripeCancelURL      string `json:"stripe_cancel_url"`
	StripeEnable         bool   `json:"stripe_enable"`
}