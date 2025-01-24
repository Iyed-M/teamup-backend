package auth

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
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

	claims := &Claims{}
	token, err := jwt.ParseWithClaims(refreshToken, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fiber.NewError(fiber.StatusUnauthorized, "Invalid token signing method")
		}
		return []byte(a.JWTSecret), nil
	})
	if err != nil {
		return fiber.NewError(fiber.StatusUnauthorized, "unparsable refresh token")
	}

	if !token.Valid {
		return fiber.NewError(fiber.StatusUnauthorized, "refresh token not valid")
	}

	userID := claims.UserID
	id, err := uuid.Parse(userID)
	if err != nil {
		return err
	}
	tokenDB, err := a.db.GetRefreshToken(c.Context(), id)
	log.Info("tokendb=", tokenDB)
	log.Info("userID=", userID)
	if err != nil || tokenDB == nil {
		log.Info("err", err)
		return fiber.NewError(fiber.StatusUnauthorized, "Invalid Refresh token")
	}
	if refreshToken != *tokenDB {
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
