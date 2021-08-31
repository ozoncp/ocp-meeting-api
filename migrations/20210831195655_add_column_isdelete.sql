-- +goose Up
-- +goose StatementBegin
ALTER TABLE meeting ADD COLUMN isDeleted BOOLEAN NOT NULL DEFAULT FALSE;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE meeting DROP COLUMN isDeleted;
-- +goose StatementEnd
