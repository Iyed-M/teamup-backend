package auth

import (
	"context"
	"net/http"

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

func (a authService) Signup(c *fiber.Ctx) error {
	ctx := context.Background()
	var req signupRequest
	if err := c.BodyParser(&req); err != nil {
		log.Errorw("parseError", "err", err)
		return fiber.NewError(fiber.StatusBadRequest, "Invalid request body")
	}

	_, err := a.db.GetUserByEmail(ctx, req.Email)
	if err == nil {
		return fiber.NewError(http.StatusConflict, "Email already exists")
	} else if err != pgx.ErrNoRows {
		log.Errorw("db error", "err", err)
		return fiber.NewError(http.StatusInternalServerError, "Database error")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Errorw("Cant hash password", "err", err)
		return fiber.NewError(http.StatusInternalServerError, "Failed to hash password")
	}

	userID := uuid.New()

	accessToken, err := a.jwt.newAccessToken(userID)
	if err != nil {
		return fiber.NewError(http.StatusInternalServerError, "Failed to generate access token")
	}

	refreshToken, err := a.jwt.newAccessToken(userID)
	if err != nil {
		return fiber.NewError(http.StatusInternalServerError, "Failed to generate refresh token")
	}

	user, err := a.db.CreateUser(ctx, repository.CreateUserParams{
		Email:        req.Email,
		Password:     string(hashedPassword),
		RefreshToken: &refreshToken,
		Username:     req.Username,
	})
	if err != nil {
		log.Errorw("cant create user", "err", err)
		return fiber.NewError(http.StatusInternalServerError, "Failed to create user")
	}

	return c.JSON(SignupResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		Email:        user.Email,
		Username:     user.Username,
	})
}
