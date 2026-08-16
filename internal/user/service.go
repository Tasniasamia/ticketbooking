package user

import (
	"errors"
	"fmt"
	"ticketBooking/internal/auth"
	"ticketBooking/internal/otp"
	"ticketBooking/internal/user/dto"
)



type Service struct{
	repo Repository;
	jwt auth.JwtService
	otp  *otp.Service
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


func (s *Service) GetMe(userId uint) (*dto.UserResponse, error) {
user,err :=s.repo.GetUserById(userId);

if err != nil{
	return nil,err;
}

if user == nil{
	return nil,nil

}

if user.Role == dto.MANAGER {
	return &dto.UserResponse{
		Id:             user.ID,
		Name:           user.Name,
		Email:          user.Email,
		Role:           user.Role,
		Address:        user.Address,
		PhoneNumber:    user.PhoneNumber,
		Country:        user.Country,
		Designation:    user.Designation,
		ProfileImage:   user.ProfileImage,
		ProfileImageId: user.ProfileImageId,
		CreatedAt:      user.CreatedAt,
		UpdatedAt:      user.UpdatedAt,
	}, nil
}

return &dto.UserResponse{
	Id:             user.ID,
	Name:           user.Name,
	Email:          user.Email,
	Role:           user.Role,
	Address:        user.Address,
	PhoneNumber:    user.PhoneNumber,
	Country:        user.Country,
	ProfileImage:   user.ProfileImage,
	ProfileImageId: user.ProfileImageId,
	CreatedAt:      user.CreatedAt,
	UpdatedAt:      user.UpdatedAt,
}, nil

}


func (s *Service) UpdateUser(userId uint, req *dto.UpdateRequest) (*dto.Response, error) {
	user, err := s.repo.GetUserById(userId)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.New("User not found")
	}

	user.Name = req.Name
	user.Email = req.Email

	updatedUser, err := s.repo.UpdateUser(user)
	if err != nil {
		return nil, err
	}

	return &dto.Response{
		Id:        updatedUser.ID,
		Name:      updatedUser.Name,
		Email:     updatedUser.Email,
		CreatedAt: updatedUser.CreatedAt.Format("2006-01-02 15:04:05"),
	}, nil
}

func (s *Service) DeleteUser(userId uint) error {
	return s.repo.DeleteUser(userId)
}


func (s *Service) VerifyOTP(req *dto.VerifyOTPRequest) error {
	if err := s.otp.Verify(req.Email, req.OTP, req.Reason); err != nil {
		return err
	}

	// যদি register হয় তাহলে user কে verified করে দাও
	if req.Reason == otp.ReasonRegister {
		if err := s.repo.MarkAsVerified(req.Email); err != nil {
			return fmt.Errorf("failed to verify user: %w", err)
		}
	}

	return nil
}

// ForgotPassword → OTP পাঠায়
func (s *Service) ForgotPassword(req *dto.ForgotPasswordRequest) error {
	user, err := s.repo.GetUserByEmail(req.Email)
	if err != nil || user == nil {
		// Security: user না থাকলেও success message দেওয়া ভালো
		return nil
	}

	return s.otp.GenerateAndSend(req.Email, otp.ReasonResetPassword)
}

// ResetPassword → OTP verify করে নতুন password সেট করে
func (s *Service) ResetPassword(req *dto.ResetPasswordRequest) error {
	// First verify OTP
	if err := s.otp.Verify(req.Email, req.OTP, otp.ReasonResetPassword); err != nil {
		return err
	}

	user, err := s.repo.GetUserByEmail(req.Email)
	if err != nil || user == nil {
		return fmt.Errorf("user not found")
	}

	if err := user.hashPassword(req.NewPassword); err != nil {
		return err
	}

	return nil;
}