package user;

import (
	"ticketBooking/internal/user/dto"
	    "ticketBooking/internal/auth"

)



type Service struct{
	repo Repository;
	jwt auth.JwtService
}


func NewUserService(repo Repository, jwt auth.JwtService) *Service {
	return &Service{repo: repo, jwt: jwt};
}

func (s *Service) CreateUser(req *dto.CreateRequest) (*dto.Response, error) {
	user := &User{
		Name:  req.Name,
		Email: req.Email,
	}

	if err := user.hashPassword(req.Password); err != nil {
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


func (s *Service) LoginUser(req *dto.LoginRequest) (*dto.Response, error) {
user,err :=s.repo.GetUserByEmail(req.Email);

if err != nil{
	return nil,err;
}

if user == nil{
	return nil,nil
}

err =user.checkPassword(req.Password);
if err != nil{
	return nil,err;
}




token, err := s.jwt.GenerateToken(user.ID,user.Name,user.Email)
if err != nil {
	return nil, err
}

return &dto.Response{
	Id:        user.ID,
	Name:      user.Name,
	Email:     user.Email,
	CreatedAt: user.CreatedAt.Format("2006-01-02 15:04:05"),
	Token:     token,
}, nil


}