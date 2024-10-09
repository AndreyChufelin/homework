-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS events (
  id uuid DEFAULT gen_random_uuid(),
  title VARCHAR(50) NOT NULL,
  date TIMESTAMP NOT NULL,
  end_date TIMESTAMP NOT NULL,
  description VARCHAR(300),
  user_id uuid NOT NULL,
  advance_notification_period INTERVAL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE events;
-- +goose StatementEnd
