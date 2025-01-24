package auth_handler

import (
	"github.com/Iyed-M/teamup-backend/internal/repository"
	"github.com/Iyed-M/teamup-backend/service/jwt"
)

type AuthHandler struct {
	jwtService jwt.JwtService
	repo         *repository.Queries
}

func NewAuthHandler(jwtService jwt.JwtService, repo *repository.Queries) *AuthHandler {
	return &AuthHandler{
		jwtService: jwtService,
		repo:         repo,
	}
}
