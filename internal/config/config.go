package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Dsn                 string
	Port                string
	JwtSecret           string
	CloudinaryCloudName string
	CloudinaryAPIKey    string
	CloudinaryAPISecret string
	AppEnv          string
	SMTPHost        string
	SMTPPort        string
	SMTPEmail       string
	SMTPPassword    string
	ResendAPIKey    string
	ResendFromEmail string
	
}

func LoadConfig() Config {
	err := godotenv.Load()
	if err != nil {
		log.Println("Warning: .env file not found, using environment variables")
	}

	return Config{
		Dsn:                 os.Getenv("DATABASE_URL"),
		Port:                os.Getenv("PORT"),
		JwtSecret:           os.Getenv("JWTSECRET"),
		CloudinaryCloudName: os.Getenv("CLOUDINARY_CLOUD_NAME"),
		CloudinaryAPIKey:    os.Getenv("CLOUDINARY_API_KEY"),
		CloudinaryAPISecret: os.Getenv("CLOUDINARY_API_SECRET"),
		AppEnv:          os.Getenv("APP_ENV"),
		SMTPHost:        os.Getenv("SMTP_HOST"),
		SMTPPort:        os.Getenv("SMTP_PORT"),
		SMTPEmail:       os.Getenv("SMTP_EMAIL"),
		SMTPPassword:    os.Getenv("SMTP_PASSWORD"),
		ResendAPIKey:    os.Getenv("RESEND_API_KEY"),
		ResendFromEmail: os.Getenv("RESEND_FROM_EMAIL"),
	}
}
