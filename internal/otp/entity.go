package otp;

import (
	"time"

	"gorm.io/gorm"
)

const (
	ReasonRegister      = "register"
	ReasonResetPassword = "reset_password"
	OTPExpiryMinutes    = 2
)

type OTP struct {
	gorm.Model
	Email     string    `gorm:"type:varchar(255);index;not null"`
	Code      string    `gorm:"type:varchar(10);not null"`
	Reason    string    `gorm:"type:varchar(50);not null"` // register / reset_password
	ExpiresAt time.Time `gorm:"not null"`
	Used      bool      `gorm:"default:false"`
}