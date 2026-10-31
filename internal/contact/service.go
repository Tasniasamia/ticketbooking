package contact
import ("errors"; "fmt"; "time"; "ticketBooking/internal/contact/dto"; "ticketBooking/internal/email"; "ticketBooking/internal/httpresponse"; "ticketBooking/internal/utils/query"; "gorm.io/gorm")
type Service struct { repo Repository; email *email.EmailService }
func NewService(repo Repository, emailSvc *email.EmailService) *Service { return &Service{repo: repo, email: emailSvc} }
func (s *Service) toResponse(c *Contact) *dto.ContactResponse {
	return &dto.ContactResponse{ID:c.ID,Name:c.Name,Email:c.Email,Subject:c.Subject,Message:c.Message,Status:string(c.Status),ReplyMessage:c.ReplyMessage,RepliedAt:c.RepliedAt,CreatedAt:c.CreatedAt,UpdatedAt:c.UpdatedAt}
}
func (s *Service) Create(req *dto.CreateContactRequest) (*dto.ContactResponse, error) {
	c := &Contact{Name:req.Name,Email:req.Email,Subject:req.Subject,Message:req.Message,Status:StatusPending}
	if err := s.repo.Create(c); err != nil { return nil, err }
	return s.toResponse(c), nil
}
func (s *Service) Reply(id uint, req *dto.ReplyContactRequest) (*dto.ContactResponse, error) {
	c, err := s.repo.GetByID(id)
	if err != nil { if errors.Is(err, gorm.ErrRecordNotFound) { return nil, errors.New("contact message not found") }; return nil, err }
	now := time.Now(); c.ReplyMessage = req.ReplyMessage; c.Status = StatusReplied; c.RepliedAt = &now
	if err := s.repo.Update(c); err != nil { return nil, err }
	if s.email != nil {
		_ = s.email.Send(c.Email, fmt.Sprintf("Re: %s", c.Subject),
			fmt.Sprintf("Hello %s,\n\n%s\n\n---\nYour original message:\n%s\n\nRegards,\nTicketBooking Support", c.Name, req.ReplyMessage, c.Message))
	}
	return s.toResponse(c), nil
}
func (s *Service) GetByID(id uint) (*dto.ContactResponse, error) {
	c, err := s.repo.GetByID(id)
	if err != nil { if errors.Is(err, gorm.ErrRecordNotFound) { return nil, errors.New("contact message not found") }; return nil, err }
	return s.toResponse(c), nil
}
func (s *Service) ListAdmin(p query.Params) (*httpresponse.PaginatedData, error) {
	items, total, err := s.repo.GetAll(p); if err != nil { return nil, err }
	docs := make([]*dto.ContactResponse, 0, len(items))
	for _, c := range items { docs = append(docs, s.toResponse(c)) }
	return &httpresponse.PaginatedData{Docs: docs, PaginationMeta: httpresponse.BuildPaginationMeta(total, p.Page, p.Limit)}, nil
}
