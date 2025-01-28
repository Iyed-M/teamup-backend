-- name: GetUserByEmail :one
SELECT * FROM users where users.email = @email limit 1;

-- name: GetUserByID :one
SELECT * FROM users where users.id = @id limit 1;

-- name: CreateUser :one
INSERT INTO users (
		id,
    email,
    password,
		username,
		refresh_token
) VALUES (
		@id,
    @email,
    @password,
		@username,
		@refresh_token
) RETURNING *;


-- name: UpdateRefreshToken :exec
UPDATE users 
SET refresh_token = @refresh_token
WHERE id = @user_id;

-- name: RemoveRefreshToken :exec
UPDATE users 
SET refresh_token = NULL
WHERE id = @user_id;

-- name: GetRefreshToken :one
SELECT refresh_token 
FROM users 
WHERE id = @user_id;
