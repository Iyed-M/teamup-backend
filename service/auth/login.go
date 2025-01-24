package auth

import (
	"net/http"

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

func (a authService) Login(c *fiber.Ctx) error {
	var req loginRequest
	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid request body")
	}

	if req.Email == "" || req.Password == "" {
		return fiber.NewError(fiber.StatusBadRequest, "Email and password are required")
	}

	user, err := a.db.GetUserByEmail(c.Context(), req.Email)
	if err != nil {
		return fiber.NewError(fiber.StatusUnauthorized, "Wrong email or password")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return fiber.NewError(fiber.StatusUnauthorized, "Invalid credentials")
	}

	accessToken, err := a.jwt.newAccessToken(user.ID)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to generate access token")
	}

	refreshToken, err := a.jwt.newRefreshToken(user.ID)
	if err != nil {
		return fiber.NewError(http.StatusInternalServerError, "Failed to generate refresh token")
	}
	a.db.UpdateRefreshToken(c.Context(), repository.UpdateRefreshTokenParams{
		RefreshToken: &refreshToken,
		UserID:       user.ID,
	})

	return c.JSON(loginResponse{
		Email:        user.Email,
		Username:     user.Username,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	})
}
