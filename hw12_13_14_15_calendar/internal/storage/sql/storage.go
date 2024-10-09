package sqlstorage

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"reflect"
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
}

func New(user, password, name string) *Storage {
	return &Storage{
		user:     user,
		password: password,
		name:     name,
	}
}

func (s *Storage) Connect(ctx context.Context) error {
	fmt.Println("Connnect...")
	db, err := sqlx.ConnectContext(ctx, "postgres",
		"user=postgres dbname=postgres sslmode=disable password=postgres host=localhost",
	)
	if err != nil {
		return fmt.Errorf("sqlstorage.Connect: %w", err)
	}

	s.db = db

	return nil
}

func (s *Storage) Close() error {
	err := s.db.Close()
	if err != nil {
		return fmt.Errorf("sqlstorage.Close: %w", err)
	}
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
}

func (eSQL eventSQL) sqlToEvent() storage.Event {
	var event storage.Event

	if eSQL.Description.Valid {
		event.Description = eSQL.Description.String
	}
	if eSQL.AdvanceNotificationPeriod.Valid {
		event.AdvanceNotificationPeriod, _ = time.ParseDuration(eSQL.AdvanceNotificationPeriod.String)
	}

	event.ID = eSQL.ID
	event.Title = eSQL.Title
	event.Date = eSQL.Date
	event.EndDate = eSQL.EndDate
	event.UserID = eSQL.UserID

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
	updateVal := reflect.ValueOf(update)
	var changed []string
	var values []interface{}

	for i := 0; i < updateVal.NumField(); i++ {
		updateField := updateVal.Field(i)

		if !updateField.IsZero() {
			f := fmt.Sprintf("%s = $%d", updateVal.Type().Field(i).Name, i)
			changed = append(changed, f)
			values = append(values, updateField.Interface())
		}
	}
	values = append(values, id)

	query := fmt.Sprintf("UPDATE events SET %s WHERE id = $%d", strings.Join(changed, ", "), len(values))

	res, err := s.db.ExecContext(ctx, query, values...)
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
