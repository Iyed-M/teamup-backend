package auth

import (
	"github.com/Iyed-M/teamup-backend/internal/repository"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func (a authService) Logout(c *fiber.Ctx) error {
	authHeader := c.Get("Authorization")
	claims, _, err := a.jwt.parseFromHeader(authHeader)
	if err != nil {
		return err
	}
	if claims.Type != JWTTypeAccess {
		return fiber.NewError(fiber.StatusUnauthorized, "Invalid token type")
	}

	id, err := uuid.Parse(claims.Id)
	if err != nil {
		return err
	}

	result := a.db.UpdateRefreshToken(c.Context(), repository.UpdateRefreshTokenParams{
		RefreshToken: nil,
		UserID:       id,
	})
	if result.Error != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to logout user")
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Successfully logged out",
	})
}
