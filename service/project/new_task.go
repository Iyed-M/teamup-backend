package project_service

import (
	"context"

	"github.com/Iyed-M/teamup-backend/internal/repository"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

func NewTask(ctx context.Context, conn *pgx.Conn, task repository.CreateTaskParams, userId uuid.UUID) error {
	tx, err := conn.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)
	repo := repository.New(tx)
	taskId, err := repo.CreateTask(ctx, task)
	if err != nil {
		return err
	}
	if err := repo.AssignTask(ctx, repository.AssignTaskParams{Userid: userId, Taskid: taskId}); err != nil {
		return err
	}
	return nil
}
