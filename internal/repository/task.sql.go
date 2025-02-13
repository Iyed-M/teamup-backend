// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: task.sql

package repository

import (
	"context"
	"time"

	"github.com/google/uuid"
)

const assignTask = `-- name: AssignTask :exec
INSERT INTO task_assignment (
	user_id,
	task_id
) VALUES (
	$1,
	$2
)
`

type AssignTaskParams struct {
	Userid uuid.UUID `json:"userid"`
	Taskid uuid.UUID `json:"taskid"`
}

func (q *Queries) AssignTask(ctx context.Context, arg AssignTaskParams) error {
	_, err := q.db.Exec(ctx, assignTask, arg.Userid, arg.Taskid)
	return err
}

const createTask = `-- name: CreateTask :one
INSERT INTO tasks (
		content,
		project_id,
		deadline,
		attachment_url,
		task_order
	) values (
		$1,
		$2,
		$3,
		$4,
		$5
	) returning id
`

type CreateTaskParams struct {
	Content       string     `json:"content"`
	ProjectID     uuid.UUID  `json:"projectId"`
	Deadline      *time.Time `json:"deadline"`
	AttachmentUrl *string    `json:"attachmentUrl"`
	TaskOrder     int32      `json:"taskOrder"`
}

func (q *Queries) CreateTask(ctx context.Context, arg CreateTaskParams) (uuid.UUID, error) {
	row := q.db.QueryRow(ctx, createTask,
		arg.Content,
		arg.ProjectID,
		arg.Deadline,
		arg.AttachmentUrl,
		arg.TaskOrder,
	)
	var id uuid.UUID
	err := row.Scan(&id)
	return id, err
}

const projectSubTasks = `-- name: ProjectSubTasks :many
SELECT sub_tasks.main_task_id, sub_tasks.sub_task_id from sub_tasks 
JOIN tasks on tasks.id = sub_tasks.main_task_id
WHERE project_id = $1
`

func (q *Queries) ProjectSubTasks(ctx context.Context, projectid uuid.UUID) ([]SubTask, error) {
	rows, err := q.db.Query(ctx, projectSubTasks, projectid)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []SubTask{}
	for rows.Next() {
		var i SubTask
		if err := rows.Scan(&i.MainTaskID, &i.SubTaskID); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const projectTasks = `-- name: ProjectTasks :many
SELECT users.id, users.email, users.password, users.username, users.created_at, users.refresh_token, tasks.id, tasks.content, tasks.project_id, tasks.created_at, tasks.deadline, tasks.attachment_url, tasks.task_order
FROM tasks JOIN task_assignment on tasks.id = task_assignment.task_id 
JOIN users on users.id = task_assignment.user_id
WHERE project_id = $1
ORDER BY tasks.id
`

type ProjectTasksRow struct {
	User User `json:"user"`
	Task Task `json:"task"`
}

func (q *Queries) ProjectTasks(ctx context.Context, projectid uuid.UUID) ([]ProjectTasksRow, error) {
	rows, err := q.db.Query(ctx, projectTasks, projectid)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []ProjectTasksRow{}
	for rows.Next() {
		var i ProjectTasksRow
		if err := rows.Scan(
			&i.User.ID,
			&i.User.Email,
			&i.User.Password,
			&i.User.Username,
			&i.User.CreatedAt,
			&i.User.RefreshToken,
			&i.Task.ID,
			&i.Task.Content,
			&i.Task.ProjectID,
			&i.Task.CreatedAt,
			&i.Task.Deadline,
			&i.Task.AttachmentUrl,
			&i.Task.TaskOrder,
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
