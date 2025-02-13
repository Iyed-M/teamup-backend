package jwt_service

import (
	"time"

	"github.com/Iyed-M/teamup-backend/types"
	"github.com/gofiber/fiber/v2/log"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

func (j JwtService) NewAccessToken(userId uuid.UUID) (string, error) {
	return generateToken(userId.String(), j.JWTAccessDuration, types.JWTTypeAccess, j.Secret)
}

func (j JwtService) NewRefreshToken(userId uuid.UUID) (string, error) {
	return generateToken(userId.String(), j.JWTRefreshDuration, types.JWTTypeRefresh, j.Secret)
}

func generateToken(userId string, duration time.Duration, tokenType types.JWTType, JWTSecret []byte) (string, error) {
	claims := &Claims{
		UserID: userId,
		Type:   tokenType,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: &jwt.NumericDate{time.Now().Add(duration)},
			IssuedAt:  &jwt.NumericDate{time.Now()},
			Issuer:    "teamup-backend",
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(JWTSecret)
	if err != nil {
		log.Errorw("Cant create token ", "err", err)
		return "", err
	}
	return tokenString, nil
}
