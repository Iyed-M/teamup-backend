package auth

import (
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
)

type JWTType string

var (
	JWTTypeRefresh JWTType = "refresh"
	JWTTypeAccess  JWTType = "access"
)

type Claims struct {
	Type   JWTType `json:"type"`
	UserID string  `json:"user_id"`
	jwt.StandardClaims
}

type jwtService struct {
	Secret             []byte
	JWTAccessDuration  time.Duration
	JWTRefreshDuration time.Duration
}

func newJwtService(jwtSecret []byte, jwtAccessDuration time.Duration, jwtRefreshDuration time.Duration) jwtService {
	return jwtService{
		Secret:             jwtSecret,
		JWTAccessDuration:  jwtAccessDuration,
		JWTRefreshDuration: jwtRefreshDuration,
	}
}

func (j jwtService) newAccessToken(userId uuid.UUID) (string, error) {
	return generateToken(userId.String(), j.JWTAccessDuration, JWTTypeAccess, j.Secret)
}

func (j jwtService) newRefreshToken(userId uuid.UUID) (string, error) {
	return generateToken(userId.String(), j.JWTRefreshDuration, JWTTypeRefresh, j.Secret)
}

func generateToken(userId string, duration time.Duration, tokenType JWTType, JWTSecret []byte) (string, error) {
	claims := &Claims{
		UserID: userId,
		Type:   tokenType,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(duration).Unix(),
			IssuedAt:  time.Now().Unix(),
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

func (j jwtService) parseFromHeader(authHeader string) (jwtClaims *Claims, token string, err error) {
	if authHeader == "" {
		return nil, "", fiber.NewError(fiber.StatusUnauthorized, "Missing authorization token")
	}

	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		return nil, "", fiber.NewError(fiber.StatusUnauthorized, "Invalid authorization header format")
	}

	tokenString := parts[1]

	claims := &Claims{}

	parsedToken, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fiber.NewError(fiber.StatusUnauthorized, "Invalid token signing method")
		}
		return []byte(j.Secret), nil
	})
	if err != nil {
		return nil, "", fiber.NewError(fiber.StatusUnauthorized, "unparsable token")
	}
	if !parsedToken.Valid {
		return nil, "", fiber.NewError(fiber.StatusUnauthorized, "refresh token not valid")
	}
	return claims, tokenString, err
}
