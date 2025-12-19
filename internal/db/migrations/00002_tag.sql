-- +goose Up
-- +goose StatementBegin
CREATE TYPE TAG_DATA_TYPE AS ENUM ('string', 'number', 'boolean');

CREATE TABLE IF NOT EXISTS tag_keys (
    id SERIAL PRIMARY KEY,
    team_id INTEGER NOT NULL REFERENCES teams(id) ON DELETE CASCADE,
    name TEXT NOT NULL,
    data_type TAG_DATA_TYPE NOT NULL,
    created_at TIMESTAMP DEFAULT NOW(),
    CONSTRAINT unique_tag_name UNIQUE (team_id, name)
);

-- below table is useful for predefined tag values (e.g., for dropdowns)
CREATE TABLE IF NOT EXISTS tag_values (
    id SERIAL PRIMARY KEY,
    key_id INTEGER NOT NULL REFERENCES tag_keys(id) ON DELETE CASCADE,
    value TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT NOW(),
    CONSTRAINT unique_tag_value UNIQUE (key_id, value)
);

CREATE TABLE IF NOT EXISTS member_filters (
    id SERIAL PRIMARY KEY,
    membership_id INTEGER NOT NULL REFERENCES team_memberships(id) ON DELETE CASCADE,
    key_id INTEGER NOT NULL REFERENCES tag_keys(id) ON DELETE CASCADE,
    value_id INTEGER NOT NULL REFERENCES tag_values(id) ON DELETE CASCADE,
    created_at TIMESTAMP DEFAULT NOW(),

    CONSTRAINT unique_team_member_preference UNIQUE (membership_id, key_id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS member_filters;
DROP TABLE IF EXISTS tag_values;
DROP TABLE IF EXISTS tag_keys;

DROP TYPE IF EXISTS TAG_DATA_TYPE;
-- +goose StatementEnd
