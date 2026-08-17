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
	Role RoleType `json:"role" validate:"required,oneof=admin manager user"`
	
}	

type LoginRequest struct{
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type UpdateMemberRequest struct {
	Name           *string   `json:"name" validate:"required"`
	Email          *string   `json:"email" validate:"required,email"`
	Role           *RoleType `json:"role"`
	Address        *string   `json:"address"`
	PhoneNumber    *string   `json:"phone_number"`
	Country        *string   `json:"country"`
	Designation    *string   `json:"designation"`
	ProfileImage   *string   `json:"profile_image"`
	ProfileImageId *uint     `json:"profile_image_id"`
	Status         *StatusType `json:"status"`
	IsVerified     *bool     `json:"is_verified"`
}

type UpdateRequest struct {
	Name           *string `json:"name" validate:"required"`
	Address        *string `json:"address"`
	PhoneNumber    *string `json:"phone_number"`
	Country        *string `json:"country"`
	ProfileImage   *string `json:"profile_image"`
	ProfileImageId *uint   `json:"profile_image_id"`
	Designation    *string `json:"designation"`
}
    // Role     RoleType `json:"role" gorm:"type:varchar(50);not null"`
    // Status StatusType `json:"status" gorm:"type:varchar(50);not null"`
    // IsVerified bool   `json:"is_verified"`



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

type ResendOTPRequest struct {
	Email  string `json:"email" validate:"required,email"`
	Reason string `json:"reason" validate:"required,oneof=register reset_password"`
}