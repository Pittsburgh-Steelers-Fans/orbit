-- name: GetTask :one
SELECT id, project_id, creator_id, assignee_id, title, description, status, position, priority, due_date, created_at, updated_at, completed_at
FROM tasks
WHERE id = $1;

-- name: ListProjectTasks :many
SELECT id, project_id, creator_id, assignee_id, title, description, status, position, priority, due_date, created_at, updated_at, completed_at
FROM tasks
WHERE project_id = $1
  AND ($2::text IS NULL OR status = $2)
ORDER BY position ASC, created_at DESC
LIMIT $3 OFFSET $4;

-- name: CreateTask :one
INSERT INTO tasks (project_id, creator_id, assignee_id, title, description, priority, due_date, position)
VALUES ($1, $2, $3, $4, $5, COALESCE(sqlc.narg('priority'), 'normal'), sqlc.narg('due_date'), $6)
RETURNING id, project_id, creator_id, assignee_id, title, description, status, position, priority, due_date, created_at, updated_at, completed_at;

-- name: UpdateTaskStatus :one
UPDATE tasks
SET status = $2,
    completed_at = CASE WHEN $2 = 'done' THEN NOW() ELSE NULL END,
    updated_at = NOW()
WHERE id = $1
RETURNING id, project_id, creator_id, assignee_id, title, description, status, position, priority, due_date, created_at, updated_at, completed_at;

-- name: AssignTask :one
UPDATE tasks
SET assignee_id = $2,
    updated_at = NOW()
WHERE id = $1
RETURNING id, project_id, creator_id, assignee_id, title, description, status, position, priority, due_date, created_at, updated_at, completed_at;

-- name: SetTaskPriority :one
UPDATE tasks
SET priority = $2,
    due_date = $3,
    updated_at = NOW()
WHERE id = $1
RETURNING id, project_id, creator_id, assignee_id, title, description, status, position, priority, due_date, created_at, updated_at, completed_at;
