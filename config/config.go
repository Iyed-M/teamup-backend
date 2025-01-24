package config

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/Iyed-M/teamup-backend/internal/repository"
	"github.com/jackc/pgx/v5"
	"github.com/joho/godotenv"
)

type EnvVars struct {
	DbURL              string
	Port               string
	JWTSecret          string
	JWTAccessDuration  time.Duration
	JWTRefreshDuration time.Duration
}

func ParseEnv() (*EnvVars, error) {
	err := godotenv.Load()
	if err != nil {
		return nil, fmt.Errorf("[Config] loading .env file :%v", err)
	}
	port := os.Getenv("PORT")
	if port == "" {
		return nil, fmt.Errorf("[Config] PORT env is not set")
	}
	repourl := os.Getenv("DB")
	if repourl == "" {
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
	return &EnvVars{
		Port:               port,
		DbURL:              repourl,
		JWTSecret:          jwtSecret,
		JWTAccessDuration:  time.Hour * time.Duration(access),
		JWTRefreshDuration: time.Hour * time.Duration(refresh),
	}, nil
}

func InitDB(repourl string, ctx context.Context) (*repository.Queries, *pgx.Conn, error) {
	conn, err := pgx.Connect(ctx, repourl)
	if err != nil {
		return nil, nil, err
	}

	repo := repository.New(conn)
	return repo, conn, nil
}
