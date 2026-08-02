package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)


type Config struct {
Dsn string ;
Port string;
}



func LoadConfig()Config {

	err :=godotenv.Load();

	if err != nil{
    log.Fatal("Error Loading .env file");
	}

	return Config{
		Dsn: os.Getenv("DATABASE_URL"),
		Port: os.Getenv("PORT"),
	}
	
	
	
	
	
}




