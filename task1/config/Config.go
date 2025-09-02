package config

import (
	"github.com/joho/godotenv"
	"log"
	"os"
)

type Config struct {
	INFURA_PK string
	Ak1       string
	APk2      string
	Ak2       string
}

func LoadConfig() Config {
	// Load from env
	if err := godotenv.Load("../../.env"); err != nil {
		log.Fatalln("Error loading .env file")
	}
	return Config{
		INFURA_PK: getEnv("INFURA_PK", ""),
		Ak1:       getEnv("Ak1", ""),
		Ak2:       getEnv("Ak2", ""),
		APk2:      getEnv("APk2", ""),
	}
}

func getEnv(key string, defaultValue string) string {
	if v := os.Getenv(key); v == "" {
		return defaultValue
	} else {
		return v
	}
}
