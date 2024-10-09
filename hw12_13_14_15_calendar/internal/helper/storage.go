package helper

import (
	"context"
	"fmt"
	"time"

	"github.com/AndreyChufelin/homework/hw12_13_14_15_calendar/internal/storage"
	memorystorage "github.com/AndreyChufelin/homework/hw12_13_14_15_calendar/internal/storage/memory"
	sqlstorage "github.com/AndreyChufelin/homework/hw12_13_14_15_calendar/internal/storage/sql"
)

type Storage interface {
	CreateEvent(context.Context, storage.Event) error
	GetEvent(context.Context, string) (*storage.Event, error)
	EditEvent(context.Context, string, storage.Event) error
	DeleteEvent(context.Context, string) error
	GetEventsListDay(ctx context.Context, date time.Time) ([]storage.Event, error)
	GetEventsListWeek(ctx context.Context, date time.Time) ([]storage.Event, error)
	GetEventsListMonth(ctx context.Context, date time.Time) ([]storage.Event, error)
}

type DBConfig struct {
	User     string
	Password string
	Name     string
}

type closeStorage = func() error

func cl() error {
	return nil
}

func InitStorage(ctx context.Context, dbConfig DBConfig, storageType string) (Storage, closeStorage, error) {
	var storage Storage
	c := cl
	if storageType == "sql" {
		sql := sqlstorage.New(dbConfig.User, dbConfig.Password, dbConfig.Name)
		c = sql.Close
		err := sql.Connect(ctx)
		if err != nil {
			return nil, c, fmt.Errorf("InitStorage: %w", err)
		}

		storage = sql
	} else {
		storage = memorystorage.New()
	}

	return storage, c, nil
}
