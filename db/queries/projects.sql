-- name: GetProject :one
SELECT id, owner_id, name, slug, description, archived_at, last_activity_at, created_at, updated_at
FROM projects
WHERE id = $1;

-- name: ListUserProjects :many
SELECT p.id, p.owner_id, p.name, p.slug, p.description, p.archived_at, p.last_activity_at, p.created_at, p.updated_at, pm.role
FROM projects p
JOIN project_members pm ON pm.project_id = p.id
WHERE pm.user_id = $1
  AND p.archived_at IS NULL
ORDER BY p.last_activity_at DESC
LIMIT $2 OFFSET $3;

-- name: CreateProject :one
INSERT INTO projects (owner_id, name, slug, description)
VALUES ($1, $2, $3, $4)
RETURNING id, owner_id, name, slug, description, archived_at, last_activity_at, created_at, updated_at;

-- name: AddProjectMember :exec
INSERT INTO project_members (project_id, user_id, role)
VALUES ($1, $2, $3)
ON CONFLICT (project_id, user_id) DO UPDATE
SET role = EXCLUDED.role;

-- name: ListProjectMembers :many
SELECT pm.project_id, pm.user_id, pm.role, pm.joined_at, u.email, u.display_name, u.avatar_url
FROM project_members pm
JOIN users u ON u.id = pm.user_id
WHERE pm.project_id = $1
ORDER BY pm.role ASC, pm.joined_at ASC;
