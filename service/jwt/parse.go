package jwt

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	_jwt "github.com/golang-jwt/jwt"
)

func (j JwtService) ParseFromHeader(authHeader string) (jwtClaims *Claims, token string, err error) {
	if authHeader == "" {
		return nil, "", fiber.NewError(fiber.StatusUnauthorized, "Missing authorization token")
	}

	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		return nil, "", fiber.NewError(fiber.StatusUnauthorized, "Invalid authorization header format")
	}

	tokenString := parts[1]

	claims := &Claims{}

	parsedToken, err := _jwt.ParseWithClaims(tokenString, claims, func(token *_jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*_jwt.SigningMethodHMAC); !ok {
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
