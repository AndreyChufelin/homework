package app

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/AndreyChufelin/homework/hw12_13_14_15_calendar/internal/storage"
)

type App struct {
	logger  Logger
	storage Storage
}

type Logger interface {
	Info(msg string, keysAndValues ...interface{})
	Error(msg string, keysAndValues ...interface{})
}

type Storage interface {
	CreateEvent(context.Context, storage.Event) error
	GetEvent(context.Context, string) (*storage.Event, error)
	EditEvent(context.Context, string, storage.Event) error
	DeleteEvent(context.Context, string) error
	GetEventsListDay(ctx context.Context, date time.Time) ([]storage.Event, error)
	GetEventsListWeek(ctx context.Context, date time.Time) ([]storage.Event, error)
	GetEventsListMonth(ctx context.Context, date time.Time) ([]storage.Event, error)
}

func New(logger Logger, storage Storage) *App {
	return &App{
		logger:  logger,
		storage: storage,
	}
}

func (a *App) CreateEvent(ctx context.Context, event storage.Event) error {
	err := a.storage.CreateEvent(ctx, event)
	if err != nil {
		a.logger.Error("failed to create event", slog.String("error", err.Error()))
		return fmt.Errorf("failed to create event: %w", err)
	}

	return nil
}

func (a *App) GetEvent(ctx context.Context, id string) (*storage.Event, error) {
	event, err := a.storage.GetEvent(ctx, id)
	if err != nil {
		a.logger.Error("failed to get event", slog.String("error", err.Error()))
		return nil, fmt.Errorf("failed to get event: %w", err)
	}

	return event, nil
}

func (a *App) DeleteEvent(ctx context.Context, id string) error {
	err := a.storage.DeleteEvent(ctx, id)
	if err != nil {
		a.logger.Error("failed to get event", slog.String("error", err.Error()))
		return fmt.Errorf("failed to get event: %w", err)
	}

	return nil
}

func (a *App) EditEvent(ctx context.Context, id string, event storage.Event) error {
	err := a.storage.EditEvent(ctx, id, event)
	if err != nil {
		a.logger.Error("failed to edit event", slog.String("error", err.Error()))
		return fmt.Errorf("failed to edit event: %w", err)
	}

	return nil
}

func (a *App) GetEventsListDay(ctx context.Context, date time.Time) ([]storage.Event, error) {
	events, err := a.storage.GetEventsListDay(ctx, date)
	if err != nil {
		a.logger.Error("failed to get events", slog.String("error", err.Error()))
		return nil, fmt.Errorf("failed to get events: %w", err)
	}

	return events, nil
}

func (a *App) GetEventsListWeek(ctx context.Context, date time.Time) ([]storage.Event, error) {
	events, err := a.storage.GetEventsListWeek(ctx, date)
	if err != nil {
		a.logger.Error("failed to get events", slog.String("error", err.Error()))
		return nil, fmt.Errorf("failed to get events: %w", err)
	}

	return events, nil
}

func (a *App) GetEventsListMonth(ctx context.Context, date time.Time) ([]storage.Event, error) {
	events, err := a.storage.GetEventsListMonth(ctx, date)
	if err != nil {
		a.logger.Error("failed to get events", slog.String("error", err.Error()))
		return nil, fmt.Errorf("failed to get events: %w", err)
	}

	return events, nil
}
