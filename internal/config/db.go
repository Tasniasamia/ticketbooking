package config

import (
	"log"
	// "ticketBooking/internal/user"


	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func ConnectDatabase(cfg Config)*gorm.DB {
	dsn := cfg.Dsn

	if dsn == "" {
	cfg.Dsn="postgresql://neondb_owner:npg_r7nbCawdxM2h@ep-wandering-voice-axi8y3qf-pooler.c-4.us-east-2.aws.neon.tech/neondb?sslmode=require&channel_binding=require"
	}

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		TranslateError: true,
	})
	if err != nil {
		panic("Failed to connect to database")
	} else {
		log.Println("Database connection successful")
	}

	//  db.AutoMigrate(&user.User{})
return db;
}