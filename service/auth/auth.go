package auth

import (
	"time"

	"github.com/Iyed-M/teamup-backend/internal/repository"
)

type authService struct {
	JWTSecret          []byte
	JWTAccessDuration  time.Duration
	JWTRefreshDuration time.Duration
	db                 *repository.Queries
}

func NewAuthService(JWTSecret []byte, JWTAccessDuration, JWTRefreshDuration time.Duration, db *repository.Queries) authService {
	return authService{
		JWTSecret:          JWTSecret,
		JWTAccessDuration:  JWTAccessDuration,
		JWTRefreshDuration: JWTRefreshDuration,
		db:                 db,
	}
}
