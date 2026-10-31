package faq
import (
	"errors"
	"ticketBooking/internal/faq/dto"
	"ticketBooking/internal/httpresponse"
	"ticketBooking/internal/utils/query"
	"gorm.io/gorm"
)
type Service struct{ repo Repository }
func NewService(repo Repository) *Service { return &Service{repo: repo} }
func (s *Service) toResponse(f *FAQ) *dto.FAQResponse {
	return &dto.FAQResponse{ID: f.ID, Question: f.Question, Answer: f.Answer, IsActive: f.IsActive, SortOrder: f.SortOrder, CreatedAt: f.CreatedAt, UpdatedAt: f.UpdatedAt}
}
func (s *Service) Create(req *dto.CreateFAQRequest) (*dto.FAQResponse, error) {
	f := &FAQ{Question: req.Question, Answer: req.Answer, IsActive: true}
	if req.IsActive != nil { f.IsActive = *req.IsActive }
	if req.SortOrder != nil { f.SortOrder = *req.SortOrder }
	if err := s.repo.Create(f); err != nil { return nil, err }
	return s.toResponse(f), nil
}
func (s *Service) Update(id uint, req *dto.UpdateFAQRequest) (*dto.FAQResponse, error) {
	f, err := s.repo.GetByID(id)
	if err != nil { if errors.Is(err, gorm.ErrRecordNotFound) { return nil, errors.New("faq not found") }; return nil, err }
	if req.Question != nil { f.Question = *req.Question }
	if req.Answer != nil { f.Answer = *req.Answer }
	if req.IsActive != nil { f.IsActive = *req.IsActive }
	if req.SortOrder != nil { f.SortOrder = *req.SortOrder }
	if err := s.repo.Update(f); err != nil { return nil, err }
	return s.toResponse(f), nil
}
func (s *Service) Delete(id uint) error {
	if _, err := s.repo.GetByID(id); err != nil { if errors.Is(err, gorm.ErrRecordNotFound) { return errors.New("faq not found") }; return err }
	return s.repo.Delete(id)
}
func (s *Service) GetByID(id uint) (*dto.FAQResponse, error) {
	f, err := s.repo.GetByID(id)
	if err != nil { if errors.Is(err, gorm.ErrRecordNotFound) { return nil, errors.New("faq not found") }; return nil, err }
	return s.toResponse(f), nil
}
func (s *Service) ListAdmin(p query.Params) (*httpresponse.PaginatedData, error) {
	items, total, err := s.repo.GetAll(p); if err != nil { return nil, err }
	docs := make([]*dto.FAQResponse, 0, len(items))
	for _, f := range items { docs = append(docs, s.toResponse(f)) }
	return &httpresponse.PaginatedData{Docs: docs, PaginationMeta: httpresponse.BuildPaginationMeta(total, p.Page, p.Limit)}, nil
}
func (s *Service) ListPublic(p query.Params) (*httpresponse.PaginatedData, error) {
	items, total, err := s.repo.GetAllActive(p); if err != nil { return nil, err }
	docs := make([]*dto.FAQResponse, 0, len(items))
	for _, f := range items { docs = append(docs, s.toResponse(f)) }
	return &httpresponse.PaginatedData{Docs: docs, PaginationMeta: httpresponse.BuildPaginationMeta(total, p.Page, p.Limit)}, nil
}
