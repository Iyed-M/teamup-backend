package auth

import (
	"log"
	"strings"

	"github.com/Iyed-M/teamup-backend/types"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt"
)

type refreshResponse struct {
	AccessToken string `json:"accessToken"`
}

func (a authService) Refresh(c *fiber.Ctx) error {
	// Get the Authorization header
	authHeader := c.Get("Authorization")
	if authHeader == "" {
		return fiber.NewError(fiber.StatusUnauthorized, "Missing refresh token")
	}

	// Check if the Authorization header has the Bearer scheme
	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		return fiber.NewError(fiber.StatusUnauthorized, "Invalid authorization header format")
	}

	refreshToken := parts[1]

	token, err := jwt.Parse(refreshToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fiber.NewError(fiber.StatusUnauthorized, "Invalid token signing method")
		}
		return []byte(a.JWTSecret), nil
	})
	if err != nil {
		return fiber.NewError(fiber.StatusUnauthorized, "Invalid refresh token")
	}

	if !token.Valid {
		return fiber.NewError(fiber.StatusUnauthorized, "Invalid refresh token")
	}

	claims, ok := token.Claims.(Claims)
	if !ok {
		return fiber.NewError(fiber.StatusUnauthorized, "Invalid token claims")
	}

	userID := claims.UserID

	user := &types.User{}
	a.db.Where("id = ? ", refreshToken, userID).First(&user)
	log.Println(user)
	if user.RefreshToken != &refreshToken {
		return fiber.NewError(fiber.StatusUnauthorized, "Invalid refresh token")
	}
	newAccessToken, err := a.generateToken(userID, a.JWTAccessDuration)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to generate access token")
	}

	return c.Status(fiber.StatusOK).JSON(refreshResponse{
		AccessToken: newAccessToken,
	})
}
