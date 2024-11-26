package sender

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/AndreyChufelin/homework/hw12_13_14_15_calendar/internal/logger"
)

type Sender struct {
	queue   Queue
	logger  Logger
	storage Storage
}

type Storage interface {
	SetNotified(context.Context, string) error
}

type Queue interface {
	Consume() (<-chan []byte, error)
}

type Logger interface {
	logger.Logger
}

func NewSender(queue Queue, logger Logger, storage Storage) Sender {
	return Sender{
		queue:   queue,
		logger:  logger,
		storage: storage,
	}
}

type Notification struct {
	ID     string
	Title  string
	Date   time.Time
	UserID string
}

func (s Sender) Start() error {
	s.logger.Info("starting sender")
	msgs, err := s.queue.Consume()
	if err != nil {
		return fmt.Errorf("faield consume notification queue: %w", err)
	}

	for msg := range msgs {
		var notification Notification
		err := json.Unmarshal(msg, &notification)
		if err != nil {
			s.logger.Error("Got wrong body", "err", err)
			continue
		}
		s.storage.SetNotified(context.TODO(), notification.ID)
		s.logger.Info("Received notification", "notification", msg)
	}

	return err
}
