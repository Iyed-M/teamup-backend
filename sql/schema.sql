
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email VARCHAR(255) NOT NULL UNIQUE,
    password TEXT NOT NULL,
	  username VARCHAR(255) NOT NULL,
    refresh_token TEXT
);

CREATE TABLE teams (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL,
    color VARCHAR(50),
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE team_permissions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
	team_id UUID REFERENCES teams(id) ON DELETE CASCADE,
    user_id UUID REFERENCES users(id) ON DELETE CASCADE,
    permissions TEXT[],
    UNIQUE(team_id, user_id)
);

CREATE TYPE invitation_status AS ENUM ('accepted', 'pending', 'rejected');

CREATE TABLE team_invitation (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    team_id UUID REFERENCES teams(id) ON DELETE CASCADE,
    sender_id UUID REFERENCES users(id) ON DELETE CASCADE,
    receiver_id UUID REFERENCES users(id) ON DELETE CASCADE,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    status invitation_status  NOT NULL DEFAULT 'pending' 
);

CREATE TABLE direct_messages (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    sender_id UUID REFERENCES users(id) ON DELETE CASCADE,
    receiver_id UUID REFERENCES users(id) ON DELETE CASCADE,
    message TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE projects (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL,
    team_id UUID REFERENCES teams(id) ON DELETE CASCADE,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE project_permissions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID REFERENCES users(id) ON DELETE CASCADE,
    project_id UUID REFERENCES projects(id) ON DELETE CASCADE,
    permissions TEXT[],
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

