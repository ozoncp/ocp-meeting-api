-- +goose Up
-- +goose StatementBegin
CREATE TABLE meeting
(
    id      SERIAL PRIMARY KEY,
    user_id INT NOT NULL,
    link    TEXT,
    start TIMESTAMP(0) WITH TIME ZONE,
    "end" TIMESTAMP(0) WITH TIME ZONE
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS meeting;
-- +goose StatementEnd
