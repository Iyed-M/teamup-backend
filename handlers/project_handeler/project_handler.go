package project_handler

import (
	"github.com/Iyed-M/teamup-backend/internal/repository"
	"github.com/jackc/pgx/v5"
)

type projectHandler struct {
	conn *pgx.Conn // used for transactions
	repo *repository.Queries
}

func NewProjectHandler(repo *repository.Queries, conn *pgx.Conn) *projectHandler {
	return &projectHandler{repo: repo, conn: conn}
}
