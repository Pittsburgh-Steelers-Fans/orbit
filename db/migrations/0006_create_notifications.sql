CREATE TABLE notifications (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    actor_id BIGINT REFERENCES users (id) ON DELETE SET NULL,
    project_id BIGINT REFERENCES projects (id) ON DELETE CASCADE,
    task_id BIGINT REFERENCES tasks (id) ON DELETE CASCADE,
    kind TEXT NOT NULL CHECK (kind IN ('task_assigned', 'task_commented', 'task_completed', 'project_invite')),
    payload JSONB NOT NULL DEFAULT '{}'::jsonb,
    read_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX notifications_user_unread_idx ON notifications (user_id, created_at DESC) WHERE read_at IS NULL;
CREATE INDEX notifications_task_id_idx ON notifications (task_id) WHERE task_id IS NOT NULL;
