-- +goose Up
CREATE UNIQUE INDEX payees_name_lower_unique_idx ON payees (lower(name));

-- +goose Down
DROP INDEX IF EXISTS payees_name_lower_unique_idx;
