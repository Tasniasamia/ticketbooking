package email

import (
	"fmt"
	"net/smtp"
	
    "ticketBooking/internal/config"
	"github.com/resend/resend-go/v3"
)

type EmailService struct {
	config config.Config
	env    string
}

func NewEmailService(cfg config.Config) *EmailService {
	return &EmailService{
		config: cfg,
	}
}

func (s *EmailService) SendOTP(toEmail, otp, reason string) error {
	subject := "Your Verification OTP"
	body := fmt.Sprintf("Your OTP is: %s\n\nThis OTP is valid for 2 minutes.\nDo not share this code with anyone.\n\nReason: %s", otp, reason)

	if s.config.AppEnv == "production" {
		return s.sendWithResend(toEmail, subject, body)
	}
	return s.sendWithGmail(toEmail, subject, body)
}

func (s *EmailService) sendWithGmail(to, subject, body string) error {
	from := s.config.SMTPEmail
	password := s.config.SMTPPassword
	host := s.config.SMTPHost
	port := s.config.SMTPPort

	if from == "" || password == "" {
		return fmt.Errorf("SMTP credentials not set")
	}
	if host == "" {
		host = "smtp.gmail.com"
	}
	if port == "" {
		port = "587"
	}

	msg := []byte(
		"From: " + from + "\r\n" +
			"To: " + to + "\r\n" +
			"Subject: " + subject + "\r\n" +
			"MIME-Version: 1.0\r\n" +
			"Content-Type: text/plain; charset=\"utf-8\"\r\n" +
			"\r\n" + body,
	)

	auth := smtp.PlainAuth("", from, password, host)
	addr := host + ":" + port

	if err := smtp.SendMail(addr, auth, from, []string{to}, msg); err != nil {
		return fmt.Errorf("gmail smtp error: %w", err)
	}

	fmt.Println("OTP sent via Gmail SMTP to:", to)
	return nil
}

func (s *EmailService) sendWithResend(to, subject, body string) error {
	apiKey := s.config.ResendAPIKey
	from := s.config.ResendFromEmail

	if apiKey == "" || from == "" {
		return fmt.Errorf("Resend API key or from email not set")
	}

	client := resend.NewClient(apiKey)
	params := &resend.SendEmailRequest{
		From:    from,
		To:      []string{to},
		Subject: subject,
		Text:    body,
	}

	sent, err := client.Emails.Send(params)
	if err != nil {
		return fmt.Errorf("resend error: %w", err)
	}

	fmt.Println("OTP sent via Resend. Email ID:", sent.Id)
	return nil
}