-- name: CreateTask :one
INSERT INTO tasks (
		content,
		project_id,
		deadline,
		attachment_url,
		task_order
	) values (
		@Content,
		@project_id,
		@deadline,
		@attachment_url,
		@task_order
	) returning id;

-- name: AssignTask :exec
INSERT INTO task_assignment (
	user_id,
	task_id
) VALUES (
	@UserId,
	@TaskId
);

-- name: ProjectTasks :many
SELECT sqlc.embed(users), sqlc.embed(tasks)
FROM tasks JOIN task_assignment on tasks.id = task_assignment.task_id 
JOIN users on users.id = task_assignment.user_id
WHERE project_id = @ProjectId
ORDER BY tasks.id;

-- name: ProjectSubTasks :many
SELECT sub_tasks.* from sub_tasks 
JOIN tasks on tasks.id = sub_tasks.main_task_id
WHERE project_id = @ProjectId;
