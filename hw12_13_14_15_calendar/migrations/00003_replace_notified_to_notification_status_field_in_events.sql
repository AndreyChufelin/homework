-- +goose Up
-- +goose StatementBegin
CREATE TYPE notification_status AS ENUM ('idle', 'sending', 'sent');

ALTER TABLE events
DROP COLUMN notified;

ALTER TABLE events
ADD COLUMN notification_status notification_status DEFAULT 'idle';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE events
DROP COLUMN notification_status;
ALTER TABLE events
ADD COLUMN notified BOOLEAN DEFAULT FALSE;
-- +goose StatementEnd
