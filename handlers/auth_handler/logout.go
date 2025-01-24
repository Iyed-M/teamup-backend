package auth_handler

import (
	"github.com/Iyed-M/teamup-backend/internal/repository"
	"github.com/Iyed-M/teamup-backend/types"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func (h *AuthHandler) Logout(c *fiber.Ctx) error {
	authHeader := c.Get("Authorization")
	claims, _, err := h.jwtService.ParseFromHeader(authHeader)
	if err != nil {
		return err
	}
	if claims.Type != types.JWTTypeAccess {
		return fiber.NewError(fiber.StatusUnauthorized, "Invalid token type")
	}

	id, err := uuid.Parse(claims.Id)
	if err != nil {
		return err
	}

	if err := h.repo.UpdateRefreshToken(c.Context(), repository.UpdateRefreshTokenParams{
		RefreshToken: nil,
		UserID:       id,
	}); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to logout user")
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Successfully logged out",
	})
}
