package sender

import (
	"fmt"

	"github.com/AndreyChufelin/homework/hw12_13_14_15_calendar/internal/logger"
)

type Sender struct {
	queue  Queue
	logger Logger
}

type Queue interface {
	Consume() (<-chan []byte, error)
}

type Logger interface {
	logger.Logger
}

func NewSender(queue Queue, logger Logger) Sender {
	return Sender{
		queue:  queue,
		logger: logger,
	}
}

func (s Sender) Start() error {
	msgs, err := s.queue.Consume()
	if err != nil {
		return fmt.Errorf("faield consume notification queue: %w", err)
	}

	for msg := range msgs {
		s.logger.Info("Received notification", "notification", msg)
	}

	return err
}
