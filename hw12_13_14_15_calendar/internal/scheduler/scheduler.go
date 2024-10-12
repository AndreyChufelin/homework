package scheduler

import (
	"context"
	"time"

	"github.com/AndreyChufelin/homework/hw12_13_14_15_calendar/internal/logger"
	"github.com/AndreyChufelin/homework/hw12_13_14_15_calendar/internal/storage"
)

type Scheduler struct {
	clearInterval int
	interval      int
	storage       Storage
	logger        Logger
	queue         Queue
}

type Logger interface {
	logger.Logger
}

type Queue interface {
	Publish(interface{}) error
}

type Storage interface {
	GetEventsToNotify(context.Context) ([]storage.Event, error)
	MarkNotified(context.Context, []string) error
	ClearEvents(context.Context, time.Duration) error
}

type Notification struct {
	ID     string
	Title  string
	Date   time.Time
	UserID string
}

func NewScheduler(queue Queue, clearInteval int, interval int, logger Logger, storage Storage) Scheduler {
	return Scheduler{
		queue:         queue,
		clearInterval: clearInteval,
		interval:      interval,
		logger:        logger,
		storage:       storage,
	}
}

func (s *Scheduler) Start(ctx context.Context) {
	s.logger.Info("starting scheduler")
	ticker := time.NewTicker(time.Duration(s.interval) * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			s.notifyEvents(ctx)
			s.storage.ClearEvents(ctx, time.Duration(s.clearInterval)*24*time.Hour)
		}
	}
}

func (s *Scheduler) notifyEvents(ctx context.Context) {
	logg := s.logger.With("at", "notifyEvents")
	events, err := s.storage.GetEventsToNotify(ctx)
	if err != nil {
		logg.Error("failed get events to notify", "err", err)
	}

	sended := make([]string, len(events))
	for i, event := range events {
		notification := Notification{
			ID:     event.ID,
			Title:  event.Title,
			Date:   event.Date,
			UserID: event.UserID,
		}

		err = s.queue.Publish(notification)
		if err != nil {
			logg.Warn("failed to publish notification", "id", notification.ID, "err", err)
			continue
		}
		sended[i] = notification.ID
		logg.Info("notification published", "id", notification.ID)
	}

	if len(sended) > 0 {
		err = s.storage.MarkNotified(ctx, sended)
		if err != nil {
			logg.Warn("failed to mark notification", "ids", sended, "err", err)
		}
	}
}
