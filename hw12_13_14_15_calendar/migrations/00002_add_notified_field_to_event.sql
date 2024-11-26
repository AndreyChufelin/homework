-- +goose Up
-- +goose StatementBegin
ALTER TABLE events
ADD COLUMN notified BOOLEAN DEFAULT FALSE;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE your_table
DROP COLUMN notified;
-- +goose StatementEnd
