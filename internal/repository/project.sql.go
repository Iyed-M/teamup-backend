// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: project.sql

package repository

import (
	"context"

	"github.com/google/uuid"
)

const addUserToProject = `-- name: AddUserToProject :exec
INSERT INTO user_projects ( project_id, user_id,is_owner) VALUES ( $1, $2, $3)
`

type AddUserToProjectParams struct {
	Projectid uuid.UUID `json:"projectid"`
	Userid    uuid.UUID `json:"userid"`
	Isowner   bool      `json:"isowner"`
}

func (q *Queries) AddUserToProject(ctx context.Context, arg AddUserToProjectParams) error {
	_, err := q.db.Exec(ctx, addUserToProject, arg.Projectid, arg.Userid, arg.Isowner)
	return err
}

const getProjectByID = `-- name: GetProjectByID :one
SELECT id, name, created_at, color FROM projects WHERE id = $1 LIMIT 1
`

func (q *Queries) GetProjectByID(ctx context.Context, id uuid.UUID) (Project, error) {
	row := q.db.QueryRow(ctx, getProjectByID, id)
	var i Project
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.CreatedAt,
		&i.Color,
	)
	return i, err
}

const inviteToProject = `-- name: InviteToProject :exec
INSERT INTO project_invitations (project_id, sender_id, receiver_id) VALUES($1,$2,$3)
`

type InviteToProjectParams struct {
	Projectid  uuid.UUID `json:"projectid"`
	Senderid   uuid.UUID `json:"senderid"`
	Receiverid uuid.UUID `json:"receiverid"`
}

func (q *Queries) InviteToProject(ctx context.Context, arg InviteToProjectParams) error {
	_, err := q.db.Exec(ctx, inviteToProject, arg.Projectid, arg.Senderid, arg.Receiverid)
	return err
}

const joinProject = `-- name: JoinProject :exec
INSERT INTO user_projects (project_id, user_id, is_owner) VALUES ($1, $2, false)
`

type JoinProjectParams struct {
	Projectid uuid.UUID `json:"projectid"`
	Userid    uuid.UUID `json:"userid"`
}

func (q *Queries) JoinProject(ctx context.Context, arg JoinProjectParams) error {
	_, err := q.db.Exec(ctx, joinProject, arg.Projectid, arg.Userid)
	return err
}

const listProjects = `-- name: ListProjects :many
SELECT id, name, created_at, color FROM projects
	WHERE id IN (
		SELECT project_id
		FROM user_projects
		WHERE user_projects.user_id = $1
	)
`

func (q *Queries) ListProjects(ctx context.Context, userID uuid.UUID) ([]Project, error) {
	rows, err := q.db.Query(ctx, listProjects, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []Project{}
	for rows.Next() {
		var i Project
		if err := rows.Scan(
			&i.ID,
			&i.Name,
			&i.CreatedAt,
			&i.Color,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const newProject = `-- name: NewProject :one
INSERT INTO projects (
		name,
		color
) VALUES (
		$1,
    $2
) RETURNING id, name, created_at, color
`

type NewProjectParams struct {
	Name  string `json:"name"`
	Color string `json:"color"`
}

func (q *Queries) NewProject(ctx context.Context, arg NewProjectParams) (Project, error) {
	row := q.db.QueryRow(ctx, newProject, arg.Name, arg.Color)
	var i Project
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.CreatedAt,
		&i.Color,
	)
	return i, err
}

const resondToProjectInvitation = `-- name: ResondToProjectInvitation :exec
UPDATE project_invitations SET status = $1 where id = $2
`

type ResondToProjectInvitationParams struct {
	Status       InvitationStatus `json:"status"`
	Invitationid uuid.UUID        `json:"invitationid"`
}

func (q *Queries) ResondToProjectInvitation(ctx context.Context, arg ResondToProjectInvitationParams) error {
	_, err := q.db.Exec(ctx, resondToProjectInvitation, arg.Status, arg.Invitationid)
	return err
}
