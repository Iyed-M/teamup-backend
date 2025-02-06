package project_service

import (
	"context"

	"github.com/Iyed-M/teamup-backend/internal/repository"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

func JoinProject(ctx context.Context, conn *pgx.Conn, userId, projectId uuid.UUID) error {
	tx, err := conn.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	repo := repository.New(tx)
	if err := repo.JoinProject(ctx, repository.JoinProjectParams{Projectid: projectId, Userid: userId}); err != nil {
		return err
	}
	if err := repo.DeleteInvitation(ctx, repository.DeleteInvitationParams{Userid: userId, Projectid: projectId}); err != nil {
		return err
	}
	return nil
}
