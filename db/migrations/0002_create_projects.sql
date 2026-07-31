CREATE TABLE projects (
    id BIGSERIAL PRIMARY KEY,
    owner_id BIGINT NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    name TEXT NOT NULL,
    slug TEXT NOT NULL,
    description TEXT,
    archived_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT projects_owner_slug_key UNIQUE (owner_id, slug)
);

CREATE INDEX projects_owner_id_idx ON projects (owner_id);
CREATE INDEX projects_archived_at_idx ON projects (archived_at) WHERE archived_at IS NOT NULL;
