-- name: GetOrCreateUpload :one
INSERT INTO uploads (storage_key, file_sha256, file_size, file_mime_type)
VALUES ($1, $2, $3, $4)
ON CONFLICT (file_sha256, file_size) DO NOTHING
RETURNING *, (xmax = 0) as created;

-- name: CreateUploadRef :one
INSERT INTO upload_refs (upload_id, team_id, uploader_id, file_name)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: CreateUploadRefTags :exec
INSERT INTO upload_ref_tags (upload_ref_id, key_id, value_id)
VALUES ($1, $2, $3);
