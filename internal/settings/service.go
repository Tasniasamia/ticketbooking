package settings

import (
	"errors"
	"ticketBooking/internal/settings/dto"
	"gorm.io/gorm"
)

type Service interface {
	Get() (*dto.Response, error)
	GetRaw() (*Setting, error)
	Upsert(req *dto.UpsertRequest) (*dto.Response, error)
}
type service struct{ repo Repository }
func NewService(repo Repository) Service { return &service{repo: repo} }
func (s *service) Get() (*dto.Response, error) {
	st, err := s.repo.Get()
	if errors.Is(err, gorm.ErrRecordNotFound) { return (&Setting{}).ToResponse(), nil }
	if err != nil { return nil, err }
	return st.ToResponse(), nil
}
func (s *service) GetRaw() (*Setting, error) { return s.repo.Get() }
func (s *service) Upsert(req *dto.UpsertRequest) (*dto.Response, error) {
	ex, err := s.repo.Get()
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) { return nil, err }
	var st *Setting
	if errors.Is(err, gorm.ErrRecordNotFound) { st = &Setting{} } else { st = ex }
	if req.SiteName != nil { st.SiteName = *req.SiteName }
	if req.SiteEmail != nil { st.SiteEmail = *req.SiteEmail }
	if req.ClientSideURL != nil { st.ClientSideURL = *req.ClientSideURL }
	if req.ServerSideURL != nil { st.ServerSideURL = *req.ServerSideURL }

	if req.SslCommerzeStoreID != nil { st.SslCommerzeStoreID = *req.SslCommerzeStoreID }
	if req.SslCommerzeStorePassword != nil { st.SslCommerzeStorePassword = *req.SslCommerzeStorePassword }
	if req.SslCommerzeSuccessURL != nil { st.SslCommerzeSuccessURL = *req.SslCommerzeSuccessURL }
	if req.SslCommerzeFailedURL != nil { st.SslCommerzeFailedURL = *req.SslCommerzeFailedURL }
	if req.SslCommerzeCancelURL != nil { st.SslCommerzeCancelURL = *req.SslCommerzeCancelURL }
	if req.StripePublishableKey != nil { st.StripePublishableKey = *req.StripePublishableKey }
	if req.StripeSecretKey != nil { st.StripeSecretKey = *req.StripeSecretKey }
	if req.StripeSuccessURL != nil { st.StripeSuccessURL = *req.StripeSuccessURL }
	if req.StripeFailedURL != nil { st.StripeFailedURL = *req.StripeFailedURL }
	if req.StripeCancelURL != nil { st.StripeCancelURL = *req.StripeCancelURL }
	if err := s.repo.Upsert(st); err != nil { return nil, err }
	fresh, err := s.repo.Get()
	if err != nil { return nil, err }
	return fresh.ToResponse(), nil
}
