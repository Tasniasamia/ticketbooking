package settings

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"ticketBooking/internal/settings/dto"
	"time"
	"gorm.io/gorm"
		"gorm.io/datatypes"

)

type SocialMediaLink struct{ Name, Link string }
type SocialMediaLinks []SocialMediaLink
func (s SocialMediaLinks) Value() (driver.Value, error) {
	if s == nil { return "[]", nil }
	b, e := json.Marshal(s); if e != nil { return nil, e }; return string(b), nil
}
func (s *SocialMediaLinks) Scan(v interface{}) error {
	if v == nil { *s = SocialMediaLinks{}; return nil }
	b, ok := v.([]byte)
	if !ok { str, ok := v.(string); if !ok { return errors.New("scan fail") }; b = []byte(str) }
	return json.Unmarshal(b, s)
}

type PageSettings struct{
	gorm.Model
	Slug string `json:"slug" validate:"required"`
	Content datatypes.JSON `json:"content" gorm:"type:jsonb"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
}



type Setting struct {
	ID uint `gorm:"primaryKey"`; CreatedAt, UpdatedAt time.Time; DeletedAt gorm.DeletedAt `gorm:"index"`
	SiteName, SiteEmail, SitePhone, SiteLogo, SiteAddress, SiteDescription, SiteFooter string
	ClientSideURL, ServerSideURL string
	SocialMediaLink SocialMediaLinks `gorm:"type:jsonb;default:'[]'"`
	
	SslCommerzeStoreID, SslCommerzeStorePassword, SslCommerzeSuccessURL, SslCommerzeFailedURL, SslCommerzeCancelURL string
	StripePublishableKey, StripeSecretKey, StripeSuccessURL, StripeFailedURL, StripeCancelURL, StripeWebhookSecret string
}
func (Setting) TableName() string { return "settings" }
func (s *Setting) ToResponse() *dto.Response {
	social := make([]dto.SocialMediaLink, 0, len(s.SocialMediaLink))
	for _, i := range s.SocialMediaLink { social = append(social, dto.SocialMediaLink{Name: i.Name, Link: i.Link}) }
	return &dto.Response{
		ID: s.ID, SiteName: s.SiteName, SiteEmail: s.SiteEmail, SitePhone: s.SitePhone,
		SiteLogo: s.SiteLogo, SiteAddress: s.SiteAddress, SiteDescription: s.SiteDescription, SiteFooter: s.SiteFooter,
		ClientSideURL: s.ClientSideURL, ServerSideURL: s.ServerSideURL,
		SocialMediaLink: social,
		SslCommerzeStoreID: s.SslCommerzeStoreID, SslCommerzeStorePassword: s.SslCommerzeStorePassword,
		SslCommerzeSuccessURL: s.SslCommerzeSuccessURL, SslCommerzeFailedURL: s.SslCommerzeFailedURL,
		SslCommerzeCancelURL: s.SslCommerzeCancelURL, 
		StripePublishableKey: s.StripePublishableKey, StripeSecretKey: s.StripeSecretKey,
		StripeSuccessURL: s.StripeSuccessURL, StripeFailedURL: s.StripeFailedURL,
		StripeCancelURL: s.StripeCancelURL,
		CreatedAt: s.CreatedAt.Format("2006-01-02 15:04:05"), UpdatedAt: s.UpdatedAt.Format("2006-01-02 15:04:05"),
		StripeWebhookSecret: s.StripeWebhookSecret,
	}
}
