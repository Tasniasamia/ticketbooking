package subscribe
import ("errors"; "strings"; "ticketBooking/internal/httpresponse"; "ticketBooking/internal/subscribe/dto"; "ticketBooking/internal/utils/query"; "gorm.io/gorm")
type Service struct{ repo Repository }
func NewService(repo Repository) *Service { return &Service{repo: repo} }
func (s *Service) toResponse(sub *Subscribe) *dto.SubscribeResponse {
	return &dto.SubscribeResponse{ID:sub.ID,Email:sub.Email,IsActive:sub.IsActive,CreatedAt:sub.CreatedAt,UpdatedAt:sub.UpdatedAt}
}
func (s *Service) Create(req *dto.CreateSubscribeRequest) (*dto.SubscribeResponse, error) {
	email := strings.TrimSpace(strings.ToLower(req.Email))
	existing, err := s.repo.GetByEmail(email)
	if err == nil && existing != nil { return nil, errors.New("this email is already subscribed") }
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) { return nil, err }
	sub := &Subscribe{Email: email, IsActive: true}
	if err := s.repo.Create(sub); err != nil { return nil, err }
	return s.toResponse(sub), nil
}
func (s *Service) ListAdmin(p query.Params) (*httpresponse.PaginatedData, error) {
	items, total, err := s.repo.GetAll(p); if err != nil { return nil, err }
	docs := make([]*dto.SubscribeResponse, 0, len(items))
	for _, sub := range items { docs = append(docs, s.toResponse(sub)) }
	return &httpresponse.PaginatedData{Docs: docs, PaginationMeta: httpresponse.BuildPaginationMeta(total, p.Page, p.Limit)}, nil
}
