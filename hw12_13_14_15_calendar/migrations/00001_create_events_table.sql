-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS events (
  id uuid DEFAULT gen_random_uuid(),
  title VARCHAR(30) NOT NULL,
  date TIMESTAMP NOT NULL,
  end_date TIMESTAMP NOT NULL,
  description TEXT,
  user_id uuid NOT NULL,
  advance_notification_period INTERVAL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE events;
-- +goose StatementEnd
