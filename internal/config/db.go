package config;

import (
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func ConnectDatabase(cfg Config) *gorm.DB {
	dsn := cfg.Dsn

	db, err := gorm.Open(postgres.New(postgres.Config{
		DSN:                  dsn,
		PreferSimpleProtocol: true, // Neon pooler fix
	}), &gorm.Config{
		TranslateError: true,
	})
	if err != nil {
		panic("Failed to connect to database")
	}

	log.Println("Database connection successful")
	return db
}