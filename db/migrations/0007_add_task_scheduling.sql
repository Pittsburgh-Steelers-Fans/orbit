ALTER TABLE tasks
    ADD COLUMN priority TEXT NOT NULL DEFAULT 'normal' CHECK (priority IN ('low', 'normal', 'high', 'urgent')),
    ADD COLUMN due_date DATE;

CREATE INDEX tasks_project_priority_idx ON tasks (project_id, priority, due_date) WHERE status <> 'done';
CREATE INDEX tasks_due_date_idx ON tasks (due_date) WHERE due_date IS NOT NULL;
