package user;

import (
	"ticketBooking/internal/user/dto"
)



type Service struct{
	repo Repository;
}


func NewUserService(repo Repository) *Service {
	return &Service{repo: repo};
}

func (s *Service) CreateUser(req *dto.CreateRequest) (*dto.Response, error) {
	user := &User{
		Name:  req.Name,
		Email: req.Email,
	}

	if err := user.HashPassword(req.Password); err != nil {
		return nil, err
	}

	if err := s.repo.CreateUser(user); err != nil {
		return nil, err
	}

	return &dto.Response{
		Id:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		CreatedAt: user.CreatedAt.Format("2006-01-02 15:04:05"),
	}, nil
}