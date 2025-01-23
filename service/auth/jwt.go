package auth

import (
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
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

func parseFromHeader(authHeader string, jwtSecret []byte) (userID string, err error) {
	if authHeader == "" {
		return "", fiber.NewError(fiber.StatusUnauthorized, "Missing authorization token")
	}

	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		return "", fiber.NewError(fiber.StatusUnauthorized, "Invalid authorization header format")
	}

	tokenString := parts[1]

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fiber.NewError(fiber.StatusUnauthorized, "Invalid token signing method")
		}
		return jwtSecret, nil
	})

	if err != nil || !token.Valid {
		return "", fiber.NewError(fiber.StatusUnauthorized, "Invalid token")
	}

	claims, ok := token.Claims.(Claims)
	if !ok {
		return "", fiber.NewError(fiber.StatusUnauthorized, "Invalid token claims")
	}
	return claims.UserID, err
}
