package queue

import (
	"fmt"

	amqp "github.com/rabbitmq/amqp091-go"
)

type Consumer struct {
	queue string
	conn  *amqp.Connection
	ch    *amqp.Channel
}

func NewConsumer(queue string, conn *amqp.Connection) Consumer {
	return Consumer{
		queue: queue,
		conn:  conn,
	}
}

func (c *Consumer) Start() error {
	var err error
	c.ch, err = c.conn.Channel()
	if err != nil {
		return err
	}

	return nil
}

func (c *Consumer) Stop() {
	if c.ch != nil {
		c.ch.Close()
	}
}

func (c Consumer) Consume() (<-chan []byte, error) {
	msgs, err := c.ch.Consume(
		c.queue,
		"sender",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to consume: %w", err)
	}

	result := make(chan []byte)
	go func() {
		defer close(result)
		for m := range msgs {
			result <- m.Body
		}
	}()

	return result, nil
}
