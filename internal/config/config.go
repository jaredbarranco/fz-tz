package config

import (
	"log"
	"os"
	"github.com/joho/godotenv"
)

type Config struct {
	AppName string
	AppEnv string
	GeoApiKey string
}

func LoadConfig() *Config {
	err:= godotenv.Load()
	if err != nil {
		log.Println("Error loading .env")
	}

	cfg := &Config{
		AppName: os.Getenv("APP"),
		AppEnv: os.Getenv("APP_ENV"),
		GeoApiKey: os.Getenv("GEO_API_KEY"),
	}
	return cfg
}

