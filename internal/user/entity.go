package user

import (
	"ticketBooking/internal/user/dto"
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)


type User struct {
	gorm.Model
	Name     string `json:"name" gorm:"type:varchar(100);not null"`
	Email    string `json:"email" gorm:"type:varchar(255);uniqueIndex;not null"`
	Password string `json:"password" gorm:"type:varchar(100);not null"`
	Role     dto.RoleType `json:"role" gorm:"type:varchar(50);not null"`
	Address  string `json:"address" gorm:"type:varchar(255)"`
	PhoneNumber string `json:"phone_number" gorm:"type:varchar(20)"`
	Country  string `json:"country" gorm:"type:varchar(100)"`
	Designation string `json:"designation" gorm:"type:varchar(100)"`
	ProfileImage string `json:"profile_image" gorm:"type:varchar(255)"`
	ProfileImageId uint `json:"profile_image_id" gorm:"type:bigint"`
	CreatedAt *time.Time `json:"created_at"`
	UpdatedAt *time.Time `json:"updated_at"`
	Status dto.StatusType `json:"status" gorm:"type:varchar(50);not null"`
	IsVerified bool   `json:"is_verified" gorm:"default:false"`

}


func (u *User) hashPassword(password string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.Password = string(hashedPassword)
	return nil
}

func (u *User) checkPassword(password string) error {
	return bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
}

func (u *User) buildUserResponse() *dto.UserResponse {

	response := &dto.UserResponse{
		Id:             u.ID,
		Name:           u.Name,
		Email:          u.Email,
		Role:           u.Role,
		Address:        u.Address,
		PhoneNumber:    u.PhoneNumber,
		Country:        u.Country,
		ProfileImage:   u.ProfileImage,
		ProfileImageId: u.ProfileImageId,
		IsVerified:     u.IsVerified,
		CreatedAt:      u.CreatedAt,
		UpdatedAt:      u.UpdatedAt,
		Status:         u.Status,
	}

	if u.Role == dto.MANAGER {
		response.Designation = u.Designation
	}

	return response
}