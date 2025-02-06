package project_service

import (
	"context"

	"github.com/Iyed-M/teamup-backend/internal/repository"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

func CreateProject(ctx context.Context, conn *pgx.Conn, userId uuid.UUID, teamName string, teamColor string) (repository.Project, error) {
	var project repository.Project
	tx, err := conn.Begin(ctx)
	if err != nil {
		return project, err
	}
	defer tx.Rollback(ctx)

	repo := repository.New(tx)

	project, err = repo.NewProject(ctx, repository.NewProjectParams{Name: teamName, Color: teamColor})
	if err != nil {
		return project, err
	}
	if err := repo.AddUserToProject(ctx, repository.AddUserToProjectParams{Projectid: project.ID, Userid: userId, Isowner: true}); err != nil {
		return project, err
	}

	if err = tx.Commit(ctx); err != nil {
		return project, err
	}
	return project, nil
}
