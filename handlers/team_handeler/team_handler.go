package team_handler

import (
	"github.com/Iyed-M/teamup-backend/internal/repository"
	"github.com/jackc/pgx/v5"
)

type teamHandler struct {
	conn *pgx.Conn // used for transactions
	repo *repository.Queries
}

func NewTeamHandler(repo *repository.Queries, conn *pgx.Conn) *teamHandler {
	return &teamHandler{repo: repo, conn: conn}
}
