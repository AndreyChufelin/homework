package sqlstorage

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/AndreyChufelin/homework/hw12_13_14_15_calendar/internal/storage"
	"github.com/jmoiron/sqlx"
)

type Storage struct {
	db       *sqlx.DB
	user     string
	password string
	name     string
	host     string
	port     string
}

func New(user, password, name, host, port string) *Storage {
	return &Storage{
		user:     user,
		password: password,
		name:     name,
		host:     host,
		port:     port,
	}
}

func (s *Storage) Connect(ctx context.Context) error {
	db, err := sqlx.ConnectContext(ctx, "postgres",
		fmt.Sprintf(
			"user=%s dbname=%s sslmode=disable password=%s host=%s port=%s",
			s.user,
			s.name,
			s.password,
			s.host,
			s.port,
		),
	)
	if err != nil {
		return fmt.Errorf("sqlstorage.Connect: %w", err)
	}

	s.db = db

	return nil
}

func (s *Storage) Close() error {
	if s.db == nil {
		return fmt.Errorf("sqlstorage.Close: no connection to close")
	}

	err := s.db.Close()
	if err != nil {
		return fmt.Errorf("sqlstorage.Close: %w", err)
	}

	s.db = nil

	return nil
}

type eventSQL struct {
	ID                        string
	Title                     string
	Date                      time.Time
	EndDate                   time.Time `db:"end_date"`
	Description               sql.NullString
	UserID                    string         `db:"user_id"`
	AdvanceNotificationPeriod sql.NullString `db:"advance_notification_period"`
	Notified                  bool
}

func (eSQL eventSQL) sqlToEvent() storage.Event {
	var event storage.Event

	if eSQL.Description.Valid {
		event.Description = eSQL.Description.String
	}
	if eSQL.AdvanceNotificationPeriod.Valid {
		period := eSQL.AdvanceNotificationPeriod.String
		period = strings.Replace(period, ":", "h", 1)
		period = strings.Replace(period, ":", "m", 1)
		period += "s"

		event.AdvanceNotificationPeriod, _ = time.ParseDuration(period)
	}

	event.ID = eSQL.ID
	event.Title = eSQL.Title
	event.Date = eSQL.Date
	event.EndDate = eSQL.EndDate
	event.UserID = eSQL.UserID
	event.Notified = eSQL.Notified

	return event
}

func (s *Storage) CreateEvent(ctx context.Context, event storage.Event) error {
	_, err := s.db.NamedExecContext(ctx, `INSERT INTO events 
		(title, date, end_date, description, user_id, advance_notification_period) 
		VALUES (:title, :date, :enddate, :description, :userid, :advancenotificationperiod)`, &event)
	if err != nil {
		return fmt.Errorf("creating event: %w", err)
	}

	return nil
}

func (s *Storage) GetEvent(ctx context.Context, id string) (*storage.Event, error) {
	var event eventSQL
	err := s.db.GetContext(ctx, &event, "SELECT * FROM events WHERE id=$1", id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("storagesql.GetEvent: %w", storage.ErrEventDoesntExist)
		}
		return nil, fmt.Errorf("getting event with id %s: %w", id, err)
	}

	e := event.sqlToEvent()

	return &e, nil
}

func (s *Storage) EditEvent(ctx context.Context, id string, update storage.Event) error {
	params := map[string]interface{}{
		"query_id":                  id,
		"id":                        update.ID,
		"title":                     update.Title,
		"date":                      update.Date,
		"enddate":                   update.EndDate,
		"description":               update.Description,
		"userid":                    update.UserID,
		"advancenotificationperiod": update.AdvanceNotificationPeriod,
	}
	res, err := s.db.NamedExecContext(ctx, `UPDATE events SET 
		title = :title, date = :date, end_date = :enddate, description = :description,
		user_id = :userid, advance_notification_period = :advancenotificationperiod
		WHERE id = :query_id`, params)
	if err != nil {
		return fmt.Errorf("edit event with id %s: %w", id, err)
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("faild get rows in storagesql.EditEvent: %w", err)
	}
	if rows != 1 {
		return fmt.Errorf("edit event with id %s: %w", id, storage.ErrEventDoesntExist)
	}

	return nil
}

func (s *Storage) DeleteEvent(ctx context.Context, id string) error {
	_, err := s.db.ExecContext(ctx, "DELETE FROM events WHERE id = $1", id)
	if err != nil {
		return fmt.Errorf("deleting event with id %s: %w", id, err)
	}

	return nil
}

func (s *Storage) GetEventsListDay(ctx context.Context, date time.Time) ([]storage.Event, error) {
	var eventsSQL []eventSQL
	err := s.db.SelectContext(ctx, &eventsSQL, "SELECT * FROM events WHERE date::date=$1", date)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("storagesql.GetEventsListDay: %w", storage.ErrNoEventsFound)
		}
		return nil, fmt.Errorf("event list day %s: %w", date.Format("2006-01-02"), err)
	}

	events := make([]storage.Event, len(eventsSQL))
	for i, event := range eventsSQL {
		events[i] = event.sqlToEvent()
	}

	return events, nil
}

func (s *Storage) GetEventsListWeek(ctx context.Context, date time.Time) ([]storage.Event, error) {
	var eventsSQL []eventSQL
	err := s.db.SelectContext(ctx, &eventsSQL, "SELECT * FROM events WHERE date::date>$1 and date::date<$2",
		date, date.AddDate(0, 0, 7),
	)
	if err != nil {
		return nil, fmt.Errorf("event list week %s: %w", date.Format("2006-01-02"), err)
	}

	events := make([]storage.Event, len(eventsSQL))
	for i, event := range eventsSQL {
		events[i] = event.sqlToEvent()
	}

	return events, nil
}

func (s *Storage) GetEventsListMonth(ctx context.Context, date time.Time) ([]storage.Event, error) {
	var eventsSQL []eventSQL
	err := s.db.SelectContext(ctx, &eventsSQL, "SELECT * FROM events WHERE date::date>$1 and date::date<$2",
		date, date.AddDate(0, 1, 0),
	)
	if err != nil {
		return nil, fmt.Errorf("event list month %s: %w", date.Format("2006-01-02"), err)
	}
	events := make([]storage.Event, len(eventsSQL))
	for i, event := range eventsSQL {
		events[i] = event.sqlToEvent()
	}

	return events, nil
}

func (s *Storage) GetEventsToNotify(ctx context.Context) ([]storage.Event, error) {
	var eventsSQL []eventSQL
	err := s.db.SelectContext(ctx,
		&eventsSQL,
		`SELECT id, title, date, user_id FROM events 
		WHERE date - advance_notification_period <= CURRENT_DATE AND notified = false`,
	)
	if err != nil {
		return nil, fmt.Errorf("sql.GetEventsToNotify: %w", err)
	}
	events := make([]storage.Event, len(eventsSQL))
	for i, event := range eventsSQL {
		events[i] = event.sqlToEvent()
	}

	return events, nil
}

func (s *Storage) MarkNotified(ctx context.Context, ids []string) error {
	q := fmt.Sprintf("UPDATE events SET notified = true WHERE id IN (%s)", "'"+strings.Join(ids, "', '")+"'")
	_, err := s.db.ExecContext(ctx, q)
	if err != nil {
		return fmt.Errorf("sqlstorage.MarkNotified: %w", err)
	}

	return nil
}

func (s *Storage) ClearEvents(ctx context.Context, duration time.Duration) error {
	date := time.Now().Add(-duration)
	_, err := s.db.ExecContext(ctx, "DELETE FROM events WHERE date < $1", date)
	if err != nil {
		return fmt.Errorf("failed to clear events: %w", err)
	}

	return nil
}
