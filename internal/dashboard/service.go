package dashboard

import "ticketBooking/internal/dashboard/dto"

type Service struct{ repo *Repository }

func NewService(repo *Repository) *Service { return &Service{repo} }

func (s *Service) Admin() (*dto.AdminDashboard, error)            { return s.repo.Admin() }
func (s *Service) Manager(id uint) (*dto.ManagerDashboard, error) { return s.repo.Manager(id) }
func (s *Service) User(id uint) (*dto.UserDashboard, error)       { return s.repo.User(id) }
