package jwt

import (
	"time"

	"github.com/Iyed-M/teamup-backend/types"
	"github.com/golang-jwt/jwt"
)

type Claims struct {
	Type   types.JWTType `json:"type"`
	UserID string        `json:"user_id"`
	jwt.StandardClaims
}

type JwtService struct {
	Secret             []byte
	JWTAccessDuration  time.Duration
	JWTRefreshDuration time.Duration
}

func NewJwtService(jwtSecret []byte, jwtAccessDuration time.Duration, jwtRefreshDuration time.Duration) JwtService {
	return JwtService{
		Secret:             jwtSecret,
		JWTAccessDuration:  jwtAccessDuration,
		JWTRefreshDuration: jwtRefreshDuration,
	}
}
