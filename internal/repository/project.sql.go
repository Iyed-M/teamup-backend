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

const deleteInvitation = `-- name: DeleteInvitation :exec
DELETE FROM project_invitations where receiver_id = $1 and project_id = $2
`

type DeleteInvitationParams struct {
	Userid    uuid.UUID `json:"userid"`
	Projectid uuid.UUID `json:"projectid"`
}

func (q *Queries) DeleteInvitation(ctx context.Context, arg DeleteInvitationParams) error {
	_, err := q.db.Exec(ctx, deleteInvitation, arg.Userid, arg.Projectid)
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

const getProjectInvitations = `-- name: GetProjectInvitations :many
SELECT projects.id, projects.name, projects.created_at, projects.color, users.id, users.email, users.password, users.username, users.created_at, users.refresh_token 
FROM projects 
JOIN project_invitations on projects.id = project_invitations.project_id 
JOIN users on project_invitations.sender_id = users.id
WHERE project_invitations.status = 'pending' AND receiver_id = $1
`

type GetProjectInvitationsRow struct {
	Project Project `json:"project"`
	User    User    `json:"user"`
}

func (q *Queries) GetProjectInvitations(ctx context.Context, userid uuid.UUID) ([]GetProjectInvitationsRow, error) {
	rows, err := q.db.Query(ctx, getProjectInvitations, userid)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []GetProjectInvitationsRow{}
	for rows.Next() {
		var i GetProjectInvitationsRow
		if err := rows.Scan(
			&i.Project.ID,
			&i.Project.Name,
			&i.Project.CreatedAt,
			&i.Project.Color,
			&i.User.ID,
			&i.User.Email,
			&i.User.Password,
			&i.User.Username,
			&i.User.CreatedAt,
			&i.User.RefreshToken,
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
