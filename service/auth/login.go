package auth

import (
	"net/http"

	"github.com/Iyed-M/teamup-backend/types"
	"github.com/labstack/echo/v4"
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

func (a authService) Login(c echo.Context) error {
	var req loginRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request body")
	}

	if req.Email == "" || req.Password == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "Email and password are required")
	}

	var user types.User
	if err := a.db.Where("email = ?", req.Email).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return echo.NewHTTPError(http.StatusUnauthorized, "Invalid credentials")
		}
		return echo.NewHTTPError(http.StatusInternalServerError, "Database error")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, "Invalid credentials")
	}

	accessToken, err := a.generateToken(user.ID.String(), a.JWTAccessDuration)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to generate access token")
	}

	refreshToken, err := a.generateToken(user.ID.String(), a.JWTRefreshDuration)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to generate refresh token")
	}

	return c.JSON(http.StatusOK, loginResponse{
		Email:        user.Email,
		Username:     user.Username,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	})
}
