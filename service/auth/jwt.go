package auth

import (
	"time"

	"github.com/golang-jwt/jwt"
)

type Claims struct {
	UserID string `json:"user_id"`
	jwt.StandardClaims
}

func (a authService) generateToken(userId string, duration time.Duration) (string, error) {
	claims := &Claims{
		UserID: userId,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(duration).Unix(),
			IssuedAt:  time.Now().Unix(),
			Issuer:    "teamup-backend",
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(a.JWTSecret))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}
