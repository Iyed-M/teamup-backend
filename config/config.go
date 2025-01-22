package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DbURL     string
	Port      string
	JWTSecret string
}

func NewConfig() (*Config, error) {
	err := godotenv.Load()
	if err != nil {
		return nil, fmt.Errorf("[Config] loading .env file :%v", err)
	}
	port := os.Getenv("PORT")
	if port == "" {
		return nil, fmt.Errorf("[Config] PORT env is not set")
	}
	dburl := os.Getenv("DB")
	if dburl == "" {
		return nil, fmt.Errorf("[Config] DB env is not set")
	}

	jwtSecret := os.Getenv("JWT_SECRET")
	if dburl == "" {
		return nil, fmt.Errorf("[Config] JWT_SECRET env is not set")
	}
	return &Config{
		Port:      port,
		DbURL:     dburl,
		JWTSecret: jwtSecret,
	}, nil
}
