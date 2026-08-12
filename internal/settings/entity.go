package settings

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"ticketBooking/internal/settings/dto"
	"time"

	"gorm.io/gorm"
)

// SocialMediaLink represents one social media entry
type SocialMediaLink struct {
	Name string `json:"name"`
	Link string `json:"link"`
}

// SocialMediaLinks is a slice stored as jsonb
type SocialMediaLinks []SocialMediaLink

func (s SocialMediaLinks) Value() (driver.Value, error) {
	if s == nil {
		return "[]", nil
	}
	b, err := json.Marshal(s)
	if err != nil {
		return nil, err
	}
	return string(b), nil
}

func (s *SocialMediaLinks) Scan(value interface{}) error {
	if value == nil {
		*s = SocialMediaLinks{}
		return nil
	}
	bytes, ok := value.([]byte)
	if !ok {
		str, ok := value.(string)
		if !ok {
			return errors.New("failed to scan SocialMediaLinks")
		}
		bytes = []byte(str)
	}
	return json.Unmarshal(bytes, s)
}

// Partners is a string slice stored as jsonb
type Partners []string

func (p Partners) Value() (driver.Value, error) {
	if p == nil {
		return "[]", nil
	}
	b, err := json.Marshal(p)
	if err != nil {
		return nil, err
	}
	return string(b), nil
}

func (p *Partners) Scan(value interface{}) error {
	if value == nil {
		*p = Partners{}
		return nil
	}
	bytes, ok := value.([]byte)
	if !ok {
		str, ok := value.(string)
		if !ok {
			return errors.New("failed to scan Partners")
		}
		bytes = []byte(str)
	}
	return json.Unmarshal(bytes, p)
}

// Setting is a singleton site-wide configuration record.
type Setting struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	// Site info
	SiteName        string `gorm:"size:255" json:"site_name"`
	SiteEmail       string `gorm:"size:255" json:"site_email"`
	SitePhone       string `gorm:"size:50" json:"site_phone"`
	SiteLogo        string `gorm:"type:text" json:"site_logo"`
	SiteAddress     string `gorm:"type:text" json:"site_address"`
	SiteDescription string `gorm:"type:text" json:"site_description"`
	SiteFooter      string `gorm:"type:text" json:"site_footer"`

	// Currency display defaults
	CurrencyCode   string `gorm:"size:10" json:"currency_code"`
	CurrencySymbol string `gorm:"size:10" json:"currency_symbol"`

	// URLs
	ClientSideURL string `gorm:"type:text" json:"client_side_url"`
	ServerSideURL string `gorm:"type:text" json:"server_side_url"`

	// Social & partners
	SocialMediaLink SocialMediaLinks `gorm:"type:jsonb;default:'[]'" json:"social_media_link"`
	Partner         Partners         `gorm:"type:jsonb;default:'[]'" json:"partner"`

	// SSLCommerz
	SslCommerzeStoreID       string `gorm:"size:255" json:"sslCommerze_store_id"`
	SslCommerzeStorePassword string `gorm:"size:255" json:"sslCommerze_store_password"`
	SslCommerzeSuccessURL    string `gorm:"type:text" json:"sslCommerze_success_url"`
	SslCommerzeFailedURL     string `gorm:"type:text" json:"sslCommerze_failed_url"`
	SslCommerzeCancelURL     string `gorm:"type:text" json:"sslCommerze_cancel_url"`
	SslCommerzeEnable        bool   `gorm:"default:false" json:"sslCommerze_enable"`

	// Stripe
	StripePublishableKey string `gorm:"size:255" json:"stripe_publishable_key"`
	StripeSecretKey      string `gorm:"size:255" json:"stripe_secret_key"`
	StripeSuccessURL     string `gorm:"type:text" json:"stripe_success_url"`
	StripeFailedURL      string `gorm:"type:text" json:"stripe_failed_url"`
	StripeCancelURL      string `gorm:"type:text" json:"stripe_cancel_url"`
	StripeEnable         bool   `gorm:"default:false" json:"stripe_enable"`
}

func (Setting) TableName() string {
	return "settings"
}

func (s *Setting) ToResponse() *dto.Response {
	social := make([]dto.SocialMediaLink, 0, len(s.SocialMediaLink))
	for _, item := range s.SocialMediaLink {
		social = append(social, dto.SocialMediaLink{
			Name: item.Name,
			Link: item.Link,
		})
	}

	partners := make([]string, 0, len(s.Partner))
	partners = append(partners, s.Partner...)

	return &dto.Response{
		ID:                       s.ID,
		SiteName:                 s.SiteName,
		SiteEmail:                s.SiteEmail,
		SitePhone:                s.SitePhone,
		SiteLogo:                 s.SiteLogo,
		SiteAddress:              s.SiteAddress,
		SiteDescription:          s.SiteDescription,
		SiteFooter:               s.SiteFooter,
		CurrencyCode:             s.CurrencyCode,
		CurrencySymbol:           s.CurrencySymbol,
		ClientSideURL:            s.ClientSideURL,
		ServerSideURL:            s.ServerSideURL,
		SocialMediaLink:          social,
		Partner:                  partners,
		SslCommerzeStoreID:       s.SslCommerzeStoreID,
		SslCommerzeStorePassword: s.SslCommerzeStorePassword,
		SslCommerzeSuccessURL:    s.SslCommerzeSuccessURL,
		SslCommerzeFailedURL:     s.SslCommerzeFailedURL,
		SslCommerzeCancelURL:     s.SslCommerzeCancelURL,
		SslCommerzeEnable:        s.SslCommerzeEnable,
		StripePublishableKey:     s.StripePublishableKey,
		StripeSecretKey:          s.StripeSecretKey,
		StripeSuccessURL:         s.StripeSuccessURL,
		StripeFailedURL:          s.StripeFailedURL,
		StripeCancelURL:          s.StripeCancelURL,
		StripeEnable:             s.StripeEnable,
		CreatedAt:                s.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:                s.UpdatedAt.Format("2006-01-02 15:04:05"),
	}
}