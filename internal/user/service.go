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


func NewUserService(repo Repository, jwt auth.JwtService, otp *otp.Service) *Service {
	return &Service{repo: repo, jwt: jwt, otp: otp};
}

func (s *Service) CreateUser(req *dto.CreateRequest) (*dto.Response, error) {
	// Check if email already exists
	existing, _ := s.repo.GetUserByEmail(req.Email)
	if existing != nil {
		return nil, fmt.Errorf("email already registered")
	}

	user := &User{
		Name:       req.Name,
		Email:      req.Email,
		IsVerified: false,
	}

	if err := user.hashPassword(req.Password); err != nil {
		return nil, err
	}

	if err := s.repo.CreateUser(user); err != nil {
		return nil, err
	}

	// Send OTP for registration
	if err := s.otp.GenerateAndSend(user.Email,otp.ReasonRegister); err != nil {
		return nil, fmt.Errorf("user created but failed to send OTP: %w", err)
	}

	return &dto.Response{
		Id:         user.ID,
		Name:       user.Name,
		Email:      user.Email,
		IsVerified: user.IsVerified,
		CreatedAt:  user.CreatedAt.Format("2006-01-02 15:04:05"),
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



fmt.Println("Login user is ",user);
token, err := s.jwt.GenerateToken(user.ID,user.Name,user.Email,user.Role,user.IsVerified)
if err != nil {
	return nil, err
}

return &dto.Response{
	Id:        user.ID,
	Name:      user.Name,
	Email:     user.Email,
	Role:      user.Role,
	IsVerified: user.IsVerified,
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


func (s *Service) UpdateUser(userId uint, req *dto.UpdateRequest) (*dto.UserResponse, error) {
	user, err := s.repo.GetUserById(userId)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.New("user not found")
	}

	updates := make(map[string]interface{})

	if req.Name != nil {
		updates["name"] = *req.Name
	}
	if req.Address != nil {
		updates["address"] = *req.Address
	}
	if req.PhoneNumber != nil {
		updates["phone_number"] = *req.PhoneNumber
	}
	if req.Country != nil {
		updates["country"] = *req.Country
	}
	if req.ProfileImage != nil {
		updates["profile_image"] = *req.ProfileImage
	}
	if req.ProfileImageId != nil {
		updates["profile_image_id"] = *req.ProfileImageId
	}
	if req.Designation != nil && user.Role == dto.MANAGER {
		updates["designation"] = *req.Designation
	}

	if len(updates) == 0 {
		return user.buildUserResponse(), nil // কিছুই আপডেট করার নেই
	}

	updatedUser, err := s.repo.UpdateUserFields(userId, updates)
	if err != nil {
		return nil, err
	}

	return updatedUser.buildUserResponse(), nil
}

func (s *Service) UpdateMember(userId uint, req *dto.UpdateMemberRequest) (*dto.UserResponse, error) {
	user, err := s.repo.GetUserById(userId)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.New("user not found")
	}

	updates := make(map[string]interface{})

	if req.Name != nil {
		updates["name"] = *req.Name
	}
	if req.Email != nil {
		updates["email"] = *req.Email
	}
	if req.Role != nil {
		updates["role"] = *req.Role
	}
	if req.Address != nil {
		updates["address"] = *req.Address
	}
	if req.PhoneNumber != nil {
		updates["phone_number"] = *req.PhoneNumber
	}
	if req.Country != nil {
		updates["country"] = *req.Country
	}
	if req.ProfileImage != nil {
		updates["profile_image"] = *req.ProfileImage
	}
	if req.ProfileImageId != nil {
		updates["profile_image_id"] = *req.ProfileImageId
	}
	if req.Status != nil {
		updates["status"] = *req.Status
	}
	if req.IsVerified != nil {
		updates["is_verified"] = *req.IsVerified
	}

	// Designation logic
	if req.Role != nil {
		if *req.Role == dto.MANAGER {
			if req.Designation != nil {
				updates["designation"] = *req.Designation
			}
		} else {
			updates["designation"] = "" // manager না হলে খালি করে দাও
		}
	} else if req.Designation != nil && user.Role == dto.MANAGER {
		updates["designation"] = *req.Designation
	}

	if len(updates) == 0 {
		return user.buildUserResponse(), nil
	}

	updatedUser, err := s.repo.UpdateUserFields(userId, updates)
	if err != nil {
		return nil, err
	}

	return updatedUser.buildUserResponse(), nil
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

// ResendOTP - expired হওয়ার পর আবার OTP পাঠানোর জন্য
func (s *Service) ResendOTP(req *dto.ResendOTPRequest) error {
	// Register reason হলে user exist করে কিনা চেক করা ভালো
	if req.Reason == otp.ReasonRegister {
		user, err := s.repo.GetUserByEmail(req.Email)
		if err != nil || user == nil {
			return fmt.Errorf("user not found")
		}
		if user.IsVerified {
			return fmt.Errorf("email is already verified")
		}
	}

	// Reset password reason হলে user exist করে কিনা চেক
	if req.Reason == otp.ReasonResetPassword {
		user, err := s.repo.GetUserByEmail(req.Email)
		if err != nil || user == nil {
			// Security কারণে error না দিয়ে silent successও দেওয়া যায়
			return fmt.Errorf("user not found")
		}
	}

	return s.otp.GenerateAndSend(req.Email, req.Reason)
}