package testimonial
import ("errors"; "ticketBooking/internal/httpresponse"; "ticketBooking/internal/testimonial/dto"; "ticketBooking/internal/utils/query"; "gorm.io/gorm")
type Service struct{ repo Repository }
func NewService(repo Repository) *Service { return &Service{repo: repo} }
func (s *Service) toResponse(t *Testimonial) *dto.TestimonialResponse {
	return &dto.TestimonialResponse{ID:t.ID,Name:t.Name,Role:t.Role,Text:t.Text,ImageURL:t.ImageURL,ImageID:t.ImageID,IsActive:t.IsActive,SortOrder:t.SortOrder,CreatedAt:t.CreatedAt,UpdatedAt:t.UpdatedAt}
}
func (s *Service) Create(req *dto.CreateTestimonialRequest) (*dto.TestimonialResponse, error) {
	t := &Testimonial{Name:req.Name,Role:req.Role,Text:req.Text,ImageURL:req.ImageURL,ImageID:req.ImageID,IsActive:true}
	if req.IsActive != nil { t.IsActive = *req.IsActive }
	if req.SortOrder != nil { t.SortOrder = *req.SortOrder }
	if err := s.repo.Create(t); err != nil { return nil, err }
	return s.toResponse(t), nil
}
func (s *Service) Update(id uint, req *dto.UpdateTestimonialRequest) (*dto.TestimonialResponse, error) {
	t, err := s.repo.GetByID(id)
	if err != nil { if errors.Is(err, gorm.ErrRecordNotFound) { return nil, errors.New("testimonial not found") }; return nil, err }
	if req.Name != nil { t.Name = *req.Name }
	if req.Role != nil { t.Role = *req.Role }
	if req.Text != nil { t.Text = *req.Text }
	if req.ImageURL != nil { t.ImageURL = *req.ImageURL }
	if req.ImageID != nil { t.ImageID = *req.ImageID }
	if req.IsActive != nil { t.IsActive = *req.IsActive }
	if req.SortOrder != nil { t.SortOrder = *req.SortOrder }
	if err := s.repo.Update(t); err != nil { return nil, err }
	return s.toResponse(t), nil
}
func (s *Service) Delete(id uint) error {
	if _, err := s.repo.GetByID(id); err != nil { if errors.Is(err, gorm.ErrRecordNotFound) { return errors.New("testimonial not found") }; return err }
	return s.repo.Delete(id)
}
func (s *Service) GetByID(id uint) (*dto.TestimonialResponse, error) {
	t, err := s.repo.GetByID(id)
	if err != nil { if errors.Is(err, gorm.ErrRecordNotFound) { return nil, errors.New("testimonial not found") }; return nil, err }
	return s.toResponse(t), nil
}
func (s *Service) ListAdmin(p query.Params) (*httpresponse.PaginatedData, error) {
	items, total, err := s.repo.GetAll(p); if err != nil { return nil, err }
	docs := make([]*dto.TestimonialResponse, 0, len(items))
	for _, t := range items { docs = append(docs, s.toResponse(t)) }
	return &httpresponse.PaginatedData{Docs: docs, PaginationMeta: httpresponse.BuildPaginationMeta(total, p.Page, p.Limit)}, nil
}
func (s *Service) ListPublic(p query.Params) (*httpresponse.PaginatedData, error) {
	items, total, err := s.repo.GetAllActive(p); if err != nil { return nil, err }
	docs := make([]*dto.TestimonialResponse, 0, len(items))
	for _, t := range items { docs = append(docs, s.toResponse(t)) }
	return &httpresponse.PaginatedData{Docs: docs, PaginationMeta: httpresponse.BuildPaginationMeta(total, p.Page, p.Limit)}, nil
}
