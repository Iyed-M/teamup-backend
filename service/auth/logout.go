package auth

import (
	"github.com/Iyed-M/teamup-backend/types"
	"github.com/gofiber/fiber/v2"
)

func (a authService) Logout(c *fiber.Ctx) error {
	authHeader := c.Get("Authorization")
	userID, err := parseFromHeader(authHeader, a.JWTSecret)
	if err != nil {
		return err
	}

	result := a.db.Model(&types.User{}).Where("id = ?", userID).Update("refresh_token", nil)
	if result.Error != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to logout user")
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Successfully logged out",
	})
}
