package auth_handler

import (
	"github.com/Iyed-M/teamup-backend/internal/repository"
	jwt_service "github.com/Iyed-M/teamup-backend/service/jwt"
)

type AuthHandler struct {
	jwtService jwt_service.JwtService
	repo       *repository.Queries
}

func NewAuthHandler(jwtService jwt_service.JwtService, repo *repository.Queries) *AuthHandler {
	return &AuthHandler{
		jwtService: jwtService,
		repo:       repo,
	}
}
