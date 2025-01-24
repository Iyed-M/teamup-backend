package auth_handler

import (
	"context"

	"github.com/Iyed-M/teamup-backend/internal/repository"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"golang.org/x/crypto/bcrypt"
)

type signupRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
	Username string `json:"username" validate:"required"`
}

type SignupResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	Email        string `json:"email"`
	Username     string `json:"username"`
}

func (h *AuthHandler) Signup(c *fiber.Ctx) error {
	ctx := context.Background()
	var req signupRequest
	if err := c.BodyParser(&req); err != nil {
		log.Errorw("parseError", "err", err)
		return fiber.NewError(fiber.StatusBadRequest, "Invalid request body")
	}

	_, err := h.repo.GetUserByEmail(ctx, req.Email)
	if err == nil {
		return fiber.NewError(fiber.StatusConflict, "Email already exists")
	} else if err != pgx.ErrNoRows {
		log.Errorw("repo error", "err", err)
		return fiber.NewError(fiber.StatusInternalServerError, "Database error")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Errorw("Cant hash password", "err", err)
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to hash password")
	}

	userID := uuid.New()

	accessToken, err := h.jwtService.NewAccessToken(userID)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to generate access token")
	}

	refreshToken, err := h.jwtService.NewRefreshToken(userID)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to generate refresh token")
	}

	user, err := h.repo.CreateUser(ctx, repository.CreateUserParams{
		Email:        req.Email,
		Password:     string(hashedPassword),
		RefreshToken: &refreshToken,
		Username:     req.Username,
	})
	if err != nil {
		log.Errorw("cant create user", "err", err)
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to create user")
	}

	return c.JSON(SignupResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		Email:        user.Email,
		Username:     user.Username,
	})
}
