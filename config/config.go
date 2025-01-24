package config

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	DbURL              string
	Port               string
	JWTSecret          string
	JWTAccessDuration  time.Duration
	JWTRefreshDuration time.Duration
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
	if jwtSecret == "" {
		return nil, fmt.Errorf("[Config] JWT_SECRET env is not set")
	}

	access, err := strconv.Atoi(os.Getenv("JWT_ACCESS_DURATION_HOURS"))
	if err != nil {
		return nil, fmt.Errorf("[Config] JWT_ACCESS_DURATION_HOURS env  err=%v", err)
	}

	refresh, err := strconv.Atoi(os.Getenv("JWT_REFRESH_DURATION_HOURS"))
	if err != nil {
		return nil, fmt.Errorf("[Config] JWT_REFRESH_DURATION_HOURS env  err=%v", err)
	}
	return &Config{
		Port:               port,
		DbURL:              dburl,
		JWTSecret:          jwtSecret,
		JWTAccessDuration:  time.Hour * time.Duration(access),
		JWTRefreshDuration: time.Hour * time.Duration(refresh),
	}, nil
}
