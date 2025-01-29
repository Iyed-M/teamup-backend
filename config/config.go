package config

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/Iyed-M/teamup-backend/internal/repository"
	"github.com/jackc/pgx/v5"
	"github.com/joho/godotenv"
)

type EnvVars struct {
	DbURL              string
	Port               string
	JWTSecret          string
	JWTAccessDuration  time.Duration
	JWTRefreshDuration time.Duration
}

func ParseEnv() (*EnvVars, error) {
	err := godotenv.Load()
	if err != nil {
		return nil, fmt.Errorf("[Config] loading .env file :%v", err)
	}
	port := os.Getenv("PORT")
	if port == "" {
		return nil, fmt.Errorf("[Config] PORT env is not set")
	}
	repourl := os.Getenv("DB")
	if repourl == "" {
		return nil, fmt.Errorf("[Config] DB env is not set")
	}

	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		return nil, fmt.Errorf("[Config] JWT_SECRET env is not set")
	}

	access, err := strconv.Atoi(os.Getenv("JWT_ACCESS_DURATION_HOURS"))
	if err != nil {
		return nil, fmt.Errorf("[Config] JWT_ACCESS_DURATION_HOURS env  err=%v", err)
	}

	refresh, err := strconv.Atoi(os.Getenv("JWT_REFRESH_DURATION_HOURS"))
	if err != nil {
		return nil, fmt.Errorf("[Config] JWT_REFRESH_DURATION_HOURS env  err=%v", err)
	}
	return &EnvVars{
		Port:               port,
		DbURL:              repourl,
		JWTSecret:          jwtSecret,
		JWTAccessDuration:  time.Hour * time.Duration(access),
		JWTRefreshDuration: time.Hour * time.Duration(refresh),
	}, nil
}

func InitDB(repourl string, ctx context.Context) (*repository.Queries, *pgx.Conn, error) {
	conn, err := pgx.Connect(ctx, repourl)
	if err != nil {
		return nil, nil, err
	}

	_, err = conn.Exec(ctx, `
CREATE OR REPLACE VIEW project_data AS 
SELECT 
    projects.id as id,
    projects.name as name,
    projects.color as color,
    projects.created_at as created_at,
    COALESCE(json_agg(
        DISTINCT jsonb_build_object(
            'id', tasks.id,
            'content', tasks.content,
            'created_at', tasks.created_at,
            'deadline', tasks.deadline,
            'attachment_url', tasks.attachment_url,
            'task_order', tasks.task_order,
            'sub_tasks', (
                SELECT COALESCE(json_agg(
                    jsonb_build_object(
                        'id', st.id,
                        'content', st.content,
                        'created_at', st.created_at,
                        'deadline', st.deadline,
                        'attachment_url', st.attachment_url,
                        'task_order', st.task_order,
                        'assignments', (
                            SELECT COALESCE(json_agg(
                                jsonb_build_object(
                                    'id', sta.id,
                                    'user_id', u.id,
                                    'email', u.email,
                                    'username', u.username
                                )
                            ), '[]')
                            FROM task_assignment sta
                            JOIN users u ON u.id = sta.user_id
                            WHERE sta.task_id = st.id
                        )
                    )
                ), '[]')
                FROM sub_tasks sub
                JOIN tasks st ON st.id = sub.sub_task_id
                WHERE sub.main_task_id = tasks.id
            ),
            'assignments', (
                SELECT COALESCE(json_agg(
                    jsonb_build_object(
                        'id', ta.id,
                        'user_id', ta.user_id
                    )
                ), '[]')
                FROM task_assignment ta
                WHERE ta.task_id = tasks.id
            )
        )
    ) FILTER (WHERE tasks.id IS NOT NULL), '[]') as tasks
FROM projects 
LEFT JOIN tasks ON projects.id = tasks.project_id 
GROUP BY projects.id;
	`)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create view: %w", err)
	}

	repo := repository.New(conn)
	return repo, conn, nil
}
