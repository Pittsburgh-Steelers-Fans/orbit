ALTER TABLE projects
    ADD COLUMN last_activity_at TIMESTAMPTZ NOT NULL DEFAULT NOW();

ALTER TABLE users
    ADD COLUMN disabled_at TIMESTAMPTZ;

CREATE INDEX projects_last_activity_at_idx ON projects (last_activity_at DESC);
CREATE INDEX tasks_project_updated_at_idx ON tasks (project_id, updated_at DESC);
CREATE INDEX users_active_email_idx ON users (email) WHERE disabled_at IS NULL;
