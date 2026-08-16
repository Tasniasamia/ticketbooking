package otp

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"time"

	"ticketBooking/internal/email"
)

type Service struct {
	repo  Repository
	email *email.EmailService
}

func NewOTPService(repo Repository, emailSvc *email.EmailService) *Service {
	return &Service{
		repo:  repo,
		email: emailSvc,
	}
}

// GenerateAndSend creates OTP + sends email
// Rate limit: same email + reason এর জন্য ২ মিনিটের মধ্যে নতুন request করা যাবে না
func (s *Service) GenerateAndSend(emailAddr, reason string) error {
	// Rate limit check
	hasActive, err := s.repo.HasActiveOTP(emailAddr, reason)
	if err != nil {
		return fmt.Errorf("failed to check active otp: %w", err)
	}
	if hasActive {
		return fmt.Errorf("please wait 2 minutes before requesting another OTP")
	}

	code, err := generateOTP(6)
	if err != nil {
		return err
	}

	otp := &OTP{
		Email:     emailAddr,
		Code:      code,
		Reason:    reason,
		ExpiresAt: time.Now().Add(OTPExpiryMinutes * time.Minute),
		Used:      false,
	}

	if err := s.repo.Create(otp); err != nil {
		return fmt.Errorf("failed to save otp: %w", err)
	}

	// Send email
	if err := s.email.SendOTP(emailAddr, code, reason); err != nil {
		return fmt.Errorf("failed to send otp email: %w", err)
	}

	return nil
}

// Verify checks OTP and marks it as used
func (s *Service) Verify(emailAddr, code, reason string) error {
	otp, err := s.repo.FindValidOTP(emailAddr, reason)
	if err != nil {
		return fmt.Errorf("invalid or expired OTP")
	}

	if otp.Code != code {
		return fmt.Errorf("invalid OTP")
	}

	// Mark as used
	if err := s.repo.MarkAsUsed(otp.ID); err != nil {
		return fmt.Errorf("failed to mark otp as used: %w", err)
	}

	return nil
}

func generateOTP(length int) (string, error) {
	const digits = "0123456789"
	result := make([]byte, length)

	for i := 0; i < length; i++ {
		num, err := rand.Int(rand.Reader, big.NewInt(int64(len(digits))))
		if err != nil {
			return "", err
		}
		result[i] = digits[num.Int64()]
	}
	return string(result), nil
}