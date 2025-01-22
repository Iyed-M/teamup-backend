package auth

import (
	"log"
	"net/http"
	"strings"

	"github.com/Iyed-M/teamup-backend/types"
	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
)

type refreshResponse struct {
	AccessToken string `json:"accessToken"`
}

func (a authService) Refresh(c echo.Context) error {
	// Get the Authorization header
	authHeader := c.Request().Header.Get("Authorization")
	if authHeader == "" {
		return echo.NewHTTPError(http.StatusUnauthorized, "Missing refresh token")
	}

	// Check if the Authorization header has the Bearer scheme
	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		return echo.NewHTTPError(http.StatusUnauthorized, "Invalid authorization header format")
	}

	refreshToken := parts[1]

	token, err := jwt.Parse(refreshToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, echo.NewHTTPError(http.StatusUnauthorized, "Invalid token signing method")
		}
		return []byte(a.JWTSecret), nil
	})
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, "Invalid refresh token")
	}

	if !token.Valid {
		return echo.NewHTTPError(http.StatusUnauthorized, "Invalid refresh token")
	}

	claims, ok := token.Claims.(Claims)
	if !ok {
		return echo.NewHTTPError(http.StatusUnauthorized, "Invalid token claims")
	}

	userID := claims.UserID

	user := &types.User{}
	a.db.Where("id = ? ", refreshToken, userID).First(&user)
	log.Println(user)
	if user.RefreshToken != &refreshToken {
		return echo.NewHTTPError(401, "Invalid refresh token")
	}
	newAccessToken, err := a.generateToken(userID, a.JWTAccessDuration)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to generate access token")
	}

	return c.JSON(http.StatusOK, refreshResponse{
		AccessToken: newAccessToken,
	})
}
