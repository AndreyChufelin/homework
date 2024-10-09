package app

import (
	"context"
	"time"

	"github.com/AndreyChufelin/homework/hw12_13_14_15_calendar/internal/storage"
)

type App struct {
	logger  Logger
	storage Storage
}

type Logger interface {
	Info(msg string, keysAndValues ...interface{})
}

type Storage interface {
	CreateEvent(context.Context, storage.Event) error
	GetEvent(context.Context, string) (*storage.Event, error)
	EditEvent(context.Context, string, storage.Event) error
	DeleteEvent(context.Context, string) error
	GetEventsListDay(ctx context.Context, date time.Time) ([]storage.Event, error)
	GetEventsListWeek(ctx context.Context, date time.Time) ([]storage.Event, error)
}

func New(logger Logger, storage Storage) *App {
	return &App{
		logger:  logger,
		storage: storage,
	}
}

func (a *App) CreateEvent(ctx context.Context, id, title string) error {
	return a.storage.CreateEvent(ctx, storage.Event{ID: id, Title: title})
}

// TODO
