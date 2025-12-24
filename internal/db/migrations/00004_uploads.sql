-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS uploads (
  id SERIAL PRIMARY KEY,
  storage_key TEXT UNIQUE NOT NULL,
  file_sha256 TEXT NOT NULL,
  file_size BIGINT NOT NULL,
  file_mime_type TEXT NOT NULL,
  created_at TIMESTAMP NOT NULL DEFAULT NOW(),
  UNIQUE (file_sha256, file_size)
);

CREATE TABLE IF NOT EXISTS upload_refs (
  id SERIAL PRIMARY KEY,
  upload_id INTEGER NOT NULL REFERENCES uploads(id),
  uploader_id INTEGER NOT NULL REFERENCES users(id),
  team_id INTEGER NOT NULL REFERENCES teams(id),
  file_name TEXT NOT NULL,
  created_at TIMESTAMP NOT NULL DEFAULT NOW(),
  UNIQUE (upload_id, team_id)
);

CREATE TABLE IF NOT EXISTS upload_ref_tags (
  id SERIAL PRIMARY KEY,
  upload_ref_id INTEGER NOT NULL REFERENCES upload_refs(id),
  key_id INTEGER NOT NULL REFERENCES tag_keys(id),
  value_id INTEGER NOT NULL REFERENCES tag_values(id),
  created_at TIMESTAMP NOT NULL DEFAULT NOW(),
  UNIQUE (upload_ref_id, key_id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS upload_ref_tags;
DROP TABLE IF EXISTS upload_refs;
DROP TABLE IF EXISTS uploads;
-- +goose StatementEnd
