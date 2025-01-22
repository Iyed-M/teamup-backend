package auth

import (
	"errors"
	"net/http"

	"github.com/Iyed-M/teamup-backend/types"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
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

func (a authService) Signup(c echo.Context) error {
	var req signupRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request body")
	}

	// Check if email already exists
	var existingUser types.User
	err := a.db.Where("email = ?", req.Email).First(&existingUser).Error
	if err == nil {
		return echo.NewHTTPError(http.StatusConflict, "Email already exists")
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return echo.NewHTTPError(http.StatusInternalServerError, "Database error")
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to hash password")
	}

	// Create new user
	user := types.User{
		Email:    req.Email,
		Password: string(hashedPassword),
		Username: req.Username,
		ID:       uuid.New(),
	}

	accessToken, err := a.generateToken(user.ID.String(), a.JWTAccessDuration)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to generate access token")
	}

	refreshToken, err := a.generateToken(user.ID.String(), a.JWTRefreshDuration)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to generate refresh token")
	}

	user.RefreshToken = &refreshToken
	if err := a.db.Create(&user).Error; err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to create user")
	}
	return c.JSON(http.StatusCreated, SignupResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		Email:        user.Email,
		Username:     user.Username,
	})
}
