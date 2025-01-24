package auth_handler

import (
	"github.com/Iyed-M/teamup-backend/internal/repository"
	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
)

type loginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
}

type loginResponse struct {
	Email        string `json:"email"`
	Username     string `json:"username"`
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}

func (h *AuthHandler) Login(c *fiber.Ctx) error {
	var req loginRequest
	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid request body")
	}

	if req.Email == "" || req.Password == "" {
		return fiber.NewError(fiber.StatusBadRequest, "Email and password are required")
	}

	user, err := h.repo.GetUserByEmail(c.Context(), req.Email)
	if err != nil {
		return fiber.NewError(fiber.StatusUnauthorized, "Wrong email or password")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return fiber.NewError(fiber.StatusUnauthorized, "Invalid credentials")
	}

	accessToken, err := h.jwtService.NewAccessToken(user.ID)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to generate access token")
	}

	refreshToken, err := h.jwtService.NewRefreshToken(user.ID)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to generate refresh token")
	}

	if err := h.repo.UpdateRefreshToken(c.Context(), repository.UpdateRefreshTokenParams{
		RefreshToken: &refreshToken,
		UserID:       user.ID,
	}); err != nil {
		return err
	}

	return c.JSON(loginResponse{
		Email:        user.Email,
		Username:     user.Username,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	})
}
