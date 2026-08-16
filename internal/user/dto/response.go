package dto

import "time"

type Response struct {
	Id        uint   `json:"id"`
	Name      string `json:"name"`
	Email     string `json:"email"`
	CreatedAt string `json:"created_at"`
	IsVerified bool   `json:"is_verified"`
	Token     string `json:"token,omitempty"`
		Role           RoleType   `json:"role"`

}

// type MemberResponse struct{
// 	Id uint `json:"id"`
//     Name  string `json:"name" validate:"required"`
// 	Email string `json:"email" validate:"required,email"`
// 	Role     RoleType `json:"role" gorm:"type:varchar(50);not null"`
// 	Address  string `json:"address" gorm:"type:varchar(255)"`
// 	PhoneNumber string `json:"phone_number" gorm:"type:varchar(20)"`
// 	Country  string `json:"country" gorm:"type:varchar(100)"`
// 	Designation string `json:"designation" gorm:"type:varchar(100)"`
// 	ProfileImage string `json:"profile_image" gorm:"type:varchar(255)"`
// 	ProfileImageId string `json:"profile_image_id" gorm:"type:varchar(255)"`
// }

// type UserResponse struct {
// 	Id             uint       `json:"id"`
// 	Name           string     `json:"name" validate:"required"`
// 	Email          string     `json:"email" validate:"required,email"`
// 	Role           RoleType   `json:"role" gorm:"type:varchar(50);not null"`
// 	Address        string     `json:"address" gorm:"type:varchar(255)"`
// 	PhoneNumber    string     `json:"phone_number" gorm:"type:varchar(20)"`
// 	Country        string     `json:"country" gorm:"type:varchar(100)"`
// 	ProfileImage   string     `json:"profile_image" gorm:"type:varchar(255)"`
// 	ProfileImageId uint       `json:"profile_image_id"`
// 	Designation    string     `json:"designation" gorm:"type:varchar(100)"`
// 	IsVerified bool   `json:"is_verified"`
// 	CreatedAt      *time.Time `json:"created_at"`
// 	UpdatedAt      *time.Time `json:"updated_at"`
// 	Status StatusType `json:"status" gorm:"type:varchar(50);not null"`


// }

type UserResponse struct {
	Id             uint       `json:"id"`
	Name           string     `json:"name"`
	Email          string     `json:"email"`
	Role           RoleType   `json:"role"`
	Address        string     `json:"address"`
	PhoneNumber    string     `json:"phone_number"`
	Country        string     `json:"country"`
	ProfileImage   string     `json:"profile_image"`
	ProfileImageId uint       `json:"profile_image_id"`
	Designation    string     `json:"designation,omitempty"`
	IsVerified     bool       `json:"is_verified"`
	CreatedAt      *time.Time `json:"created_at"`
	UpdatedAt      *time.Time `json:"updated_at"`
	Status         StatusType `json:"status"`
}