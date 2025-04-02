-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
TRUNCATE TABLE feed_follows;
-- +goose StatementEnd
