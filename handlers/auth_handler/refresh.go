package auth_handler

import (
	"github.com/Iyed-M/teamup-backend/types"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/google/uuid"
)

type refreshResponse struct {
	AccessToken string `json:"accessToken"`
}

func (h *AuthHandler) Refresh(c *fiber.Ctx) error {
	authHeader := c.Get("Authorization")
	refreshClaims, refreshToken, err := h.jwtService.ParseFromHeader(authHeader)
	if err != nil {
		return err
	}
	if refreshClaims.Type != types.JWTTypeRefresh {
		return fiber.NewError(fiber.StatusUnauthorized, "Invalid token type")
	}

	userID := refreshClaims.UserID
	id, err := uuid.Parse(userID)
	if err != nil {
		return err
	}
	tokenDB, err := h.repo.GetRefreshToken(c.Context(), id)
	if err != nil || tokenDB == nil {
		log.Info("err", err)
		return fiber.NewError(fiber.StatusUnauthorized, "Invalid Refresh token")
	}
	if refreshToken != *tokenDB {
		return fiber.NewError(fiber.StatusUnauthorized, "Invalid refresh token")
	}
	newAccessToken, err := h.jwtService.NewAccessToken(id)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to generate access token")
	}

	return c.Status(fiber.StatusOK).JSON(refreshResponse{
		AccessToken: newAccessToken,
	})
}
