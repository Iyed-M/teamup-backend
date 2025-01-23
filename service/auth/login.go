package auth

import (
	"net/http"

	"github.com/Iyed-M/teamup-backend/types"
	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
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

	var user types.User
	if err := a.db.Where("email = ?", req.Email).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return fiber.NewError(fiber.StatusUnauthorized, "Invalid credentials")
		}
		return fiber.NewError(fiber.StatusInternalServerError, "Database error")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return fiber.NewError(fiber.StatusUnauthorized, "Invalid credentials")
	}

	accessToken, err := a.generateToken(user.ID.String(), a.JWTAccessDuration)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to generate access token")
	}

	refreshToken, err := a.generateToken(user.ID.String(), a.JWTRefreshDuration)
	if err != nil {
		return fiber.NewError(http.StatusInternalServerError, "Failed to generate refresh token")
	}
	a.db.Model(&types.User{}).Where("id = ?", user.ID).Update("refresh_token", refreshToken)

	return c.JSON(loginResponse{
		Email:        user.Email,
		Username:     user.Username,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	})
}
