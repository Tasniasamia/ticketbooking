package otp

import (
	"time"

	"gorm.io/gorm"
)

type Repository interface {
	Create(otp *OTP) error
	FindValidOTP(email, reason string) (*OTP, error)
	MarkAsUsed(id uint) error
	HasActiveOTP(email, reason string) (bool, error)
}

type repository struct {
	db *gorm.DB
}

func NewOTPRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

func (r *repository) Create(otp *OTP) error {
	return r.db.Create(otp).Error
}

func (r *repository) FindValidOTP(email, reason string) (*OTP, error) {
	var otp OTP
	err := r.db.Where("email = ? AND reason = ? AND used = false AND expires_at > ?", email, reason, time.Now()).
		Order("created_at DESC").
		First(&otp).Error

	if err != nil {
		return nil, err
	}
	return &otp, nil
}

func (r *repository) MarkAsUsed(id uint) error {
	return r.db.Model(&OTP{}).Where("id = ?", id).Update("used", true).Error
}

func (r *repository) HasActiveOTP(email, reason string) (bool, error) {
	var count int64
	err := r.db.Model(&OTP{}).
		Where("email = ? AND reason = ? AND used = false AND expires_at > ?", email, reason, time.Now()).
		Count(&count).Error

	return count > 0, err
}