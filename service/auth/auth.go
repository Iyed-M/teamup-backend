package auth

import (
	"time"

	"github.com/Iyed-M/teamup-backend/internal/repository"
)

type authService struct {
	db  *repository.Queries
	jwt jwtService
}

func NewAuthService(JWTSecret []byte, JWTAccessDuration, JWTRefreshDuration time.Duration, db *repository.Queries) authService {
	return authService{
		jwt: newJwtService(JWTSecret, JWTAccessDuration, JWTRefreshDuration),
		db:  db,
	}
}
