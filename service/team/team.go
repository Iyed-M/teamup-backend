package team_service

import (
	"context"

	"github.com/Iyed-M/teamup-backend/internal/repository"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

func CreateTeam(ctx context.Context, conn *pgx.Conn, userId uuid.UUID, teamName string, teamColor *string) (repository.Team, error) {
	var team repository.Team
	tx, err := conn.Begin(ctx)
	if err != nil {
		return team, err
	}
	defer tx.Rollback(ctx)

	repo := repository.New(tx)

	team, err = repo.NewTeam(ctx, repository.NewTeamParams{Name: teamName, Color: teamColor})
	if err != nil {
		repo.AddTeamPermissions(ctx, repository.AddTeamPermissionsParams{Teamid: team.ID, Userid: userId, Isowner: true})
		return team, err
	}

	if err = tx.Commit(ctx); err != nil {
		return team, err
	}
	return team, nil
}
