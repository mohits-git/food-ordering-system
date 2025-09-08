package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	SQLITE_DSN   string
	JWT_SECRET   string
	JWT_ISSUER   string
	JWT_AUDIENCE string
}

func LoadConfig() Config {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	var config Config

	config.SQLITE_DSN = os.Getenv("SQLITE_DSN")
	if config.SQLITE_DSN == "" {
		config.SQLITE_DSN = "file:food-ordering-system.db?cache=shared&mode=rwc"
	}

	config.JWT_SECRET = os.Getenv("JWT_SECRET")
	if config.JWT_SECRET == "" {
		config.JWT_SECRET = "jwt_secret_pass"
	}

	config.JWT_ISSUER = os.Getenv("JWT_ISSUER")
	if config.JWT_ISSUER == "" {
		config.JWT_ISSUER = "jwt_issuer_name"
	}

	config.JWT_AUDIENCE = os.Getenv("JWT_AUDIENCE")
	if config.JWT_AUDIENCE == "" {
		config.JWT_AUDIENCE = "jwt_audience"
	}

	return config
}
