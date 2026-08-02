package config

import (
	"log"
	// "ticketBooking/internal/user"


	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func ConnectDatabase(cfg Config)*gorm.DB {
	dsn := cfg.Dsn



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