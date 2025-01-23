package auth

import (
	"time"

	"gorm.io/gorm"
)

type authService struct {
	JWTSecret          []byte
	JWTAccessDuration  time.Duration
	JWTRefreshDuration time.Duration
	db                 *gorm.DB
}

func NewAuthService(JWTSecret []byte, JWTAccessDuration, JWTRefreshDuration time.Duration, db *gorm.DB) authService {
	return authService{
		JWTSecret:          JWTSecret,
		JWTAccessDuration:  JWTAccessDuration,
		JWTRefreshDuration: JWTRefreshDuration,
		db:                 db,
	}
}
