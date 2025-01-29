
CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email VARCHAR(255) NOT NULL UNIQUE,
    password TEXT NOT NULL,
	  username VARCHAR(255) NOT NULL,
		created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    refresh_token TEXT
);


CREATE TYPE invitation_status AS ENUM ('accepted', 'pending', 'rejected');
CREATE TABLE projects (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
		color TEXT NOT NULL
);

CREATE TABLE project_invitations (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    project_id UUID REFERENCES projects(id) ON DELETE CASCADE,
    sender_id UUID REFERENCES users(id) ON DELETE CASCADE,
    receiver_id UUID REFERENCES users(id) ON DELETE CASCADE,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    status invitation_status  NOT NULL DEFAULT 'pending' 
);



CREATE TABLE user_projects (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID REFERENCES users(id) ON DELETE CASCADE,
    project_id UUID REFERENCES projects(id) ON DELETE CASCADE,
		is_owner BOOLEAN NOT NULL DEFAULT FALSE,
    UNIQUE(user_id, project_id)
);

CREATE TABLE tasks (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    content TEXT NOT NULL,
    project_id UUID REFERENCES projects(id) ON DELETE CASCADE,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deadline TIMESTAMP,
    attachment_url TEXT,
    task_order INTEGER NOT NULL,
    UNIQUE(project_id, task_order)
);

CREATE TABLE sub_tasks (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    main_task_id UUID REFERENCES tasks(id) ON DELETE CASCADE,
    sub_task_id UUID REFERENCES tasks(id) ON DELETE CASCADE,
    UNIQUE(main_task_id, sub_task_id)
);

CREATE TABLE task_assignment (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID REFERENCES users(id) ON DELETE CASCADE,
    task_id UUID REFERENCES tasks(id) ON DELETE CASCADE,
    UNIQUE(user_id, task_id)
);

CREATE OR REPLACE VIEW project_data AS 
SELECT 
    projects.id,
    projects.name,
    projects.color,
    projects.created_at,
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
