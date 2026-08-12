package settings

import (
	"errors"
	"ticketBooking/internal/settings/dto"

	"gorm.io/gorm"
)

type Service interface {
	Get() (*dto.Response, error)
	Upsert(req *dto.UpsertRequest) (*dto.Response, error)
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func (s *service) Get() (*dto.Response, error) {
	setting, err := s.repo.Get()
	if errors.Is(err, gorm.ErrRecordNotFound) {
		// Return empty defaults instead of 404 so frontend always gets a shape
		empty := &Setting{}
		return empty.ToResponse(), nil
	}
	if err != nil {
		return nil, err
	}
	return setting.ToResponse(), nil
}

func (s *service) Upsert(req *dto.UpsertRequest) (*dto.Response, error) {
	existing, err := s.repo.Get()
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	var setting *Setting
	if errors.Is(err, gorm.ErrRecordNotFound) {
		setting = &Setting{}
	} else {
		setting = existing
	}

	applyUpsert(setting, req)

	if err := s.repo.Upsert(setting); err != nil {
		return nil, err
	}

	// re-fetch to get updated timestamps etc.
	fresh, err := s.repo.Get()
	if err != nil {
		return nil, err
	}
	return fresh.ToResponse(), nil
}

func applyUpsert(s *Setting, req *dto.UpsertRequest) {
	if req.SiteName != nil {
		s.SiteName = *req.SiteName
	}
	if req.SiteEmail != nil {
		s.SiteEmail = *req.SiteEmail
	}
	if req.SitePhone != nil {
		s.SitePhone = *req.SitePhone
	}
	if req.SiteLogo != nil {
		s.SiteLogo = *req.SiteLogo
	}
	if req.SiteAddress != nil {
		s.SiteAddress = *req.SiteAddress
	}
	if req.SiteDescription != nil {
		s.SiteDescription = *req.SiteDescription
	}
	if req.SiteFooter != nil {
		s.SiteFooter = *req.SiteFooter
	}
	if req.CurrencyCode != nil {
		s.CurrencyCode = *req.CurrencyCode
	}
	if req.CurrencySymbol != nil {
		s.CurrencySymbol = *req.CurrencySymbol
	}
	if req.ClientSideURL != nil {
		s.ClientSideURL = *req.ClientSideURL
	}
	if req.ServerSideURL != nil {
		s.ServerSideURL = *req.ServerSideURL
	}

	if req.SocialMediaLink != nil {
		links := make(SocialMediaLinks, 0, len(*req.SocialMediaLink))
		for _, item := range *req.SocialMediaLink {
			links = append(links, SocialMediaLink{
				Name: item.Name,
				Link: item.Link,
			})
		}
		s.SocialMediaLink = links
	}

	if req.Partner != nil {
		s.Partner = Partners(*req.Partner)
	}

	// SSLCommerz
	if req.SslCommerzeStoreID != nil {
		s.SslCommerzeStoreID = *req.SslCommerzeStoreID
	}
	if req.SslCommerzeStorePassword != nil {
		s.SslCommerzeStorePassword = *req.SslCommerzeStorePassword
	}
	if req.SslCommerzeSuccessURL != nil {
		s.SslCommerzeSuccessURL = *req.SslCommerzeSuccessURL
	}
	if req.SslCommerzeFailedURL != nil {
		s.SslCommerzeFailedURL = *req.SslCommerzeFailedURL
	}
	if req.SslCommerzeCancelURL != nil {
		s.SslCommerzeCancelURL = *req.SslCommerzeCancelURL
	}
	if req.SslCommerzeEnable != nil {
		s.SslCommerzeEnable = *req.SslCommerzeEnable
	}

	// Stripe
	if req.StripePublishableKey != nil {
		s.StripePublishableKey = *req.StripePublishableKey
	}
	if req.StripeSecretKey != nil {
		s.StripeSecretKey = *req.StripeSecretKey
	}
	if req.StripeSuccessURL != nil {
		s.StripeSuccessURL = *req.StripeSuccessURL
	}
	if req.StripeFailedURL != nil {
		s.StripeFailedURL = *req.StripeFailedURL
	}
	if req.StripeCancelURL != nil {
		s.StripeCancelURL = *req.StripeCancelURL
	}
	if req.StripeEnable != nil {
		s.StripeEnable = *req.StripeEnable
	}
}