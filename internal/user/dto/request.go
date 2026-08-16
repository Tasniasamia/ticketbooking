package dto;
type RoleType string

const (
	ADMIN   RoleType = "admin"
	MANAGER   RoleType = "manager"
	USER    RoleType = "user"

)

type StatusType string;

const (
	ACTIVE StatusType = "active"
	DEACTIVATED StatusType = "deactivated"
)


type CreateRequest struct{
	Name  string `json:"name" validate:"required"`
	Email string `json:"email" validate:"required,email"`
    Password string `json:"password" validate:"required,min=6"`
	
	
}	

type LoginRequest struct{
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type UpdateMemberRequest struct{
	Name  string `json:"name" validate:"required"`
	Email string `json:"email" validate:"required,email"`
	Role     RoleType `json:"role" gorm:"type:varchar(50);not null"`
	Address  string `json:"address" gorm:"type:varchar(255)"`
	PhoneNumber string `json:"phone_number" gorm:"type:varchar(20)"`
	Country  string `json:"country" gorm:"type:varchar(100)"`
	Designation string `json:"designation" gorm:"type:varchar(100)"`
	ProfileImage string `json:"profile_image" gorm:"type:varchar(255)"`
	ProfileImageId uint `json:"profile_image_id" gorm:"type:bigint"`
	Status StatusType `json:"status" gorm:"type:varchar(50);not null"`
    IsVerified bool   `json:"is_verified"`
}

type UpdateRequest struct{
	Name  string `json:"name" validate:"required"`
	Email string `json:"email" validate:"required,email"`
	Role     RoleType `json:"role" gorm:"type:varchar(50);not null"`
	Address  string `json:"address" gorm:"type:varchar(255)"`
	PhoneNumber string `json:"phone_number" gorm:"type:varchar(20)"`
	Country  string `json:"country" gorm:"type:varchar(100)"`
	ProfileImage string `json:"profile_image" gorm:"type:varchar(255)"`
	ProfileImageId uint `json:"profile_image_id" gorm:"type:bigint"`
	Status StatusType `json:"status" gorm:"type:varchar(50);not null"`
    IsVerified bool   `json:"is_verified"`
}



type VerifyOTPRequest struct {
	Email  string `json:"email" validate:"required,email"`
	OTP    string `json:"otp" validate:"required,len=6"`
	Reason string `json:"reason" validate:"required,oneof=register reset_password"`
}

type ForgotPasswordRequest struct {
	Email string `json:"email" validate:"required,email"`
}

type ResetPasswordRequest struct {
	Email       string `json:"email" validate:"required,email"`
	OTP         string `json:"otp" validate:"required,len=6"`
	NewPassword string `json:"new_password" validate:"required,min=6"`
}