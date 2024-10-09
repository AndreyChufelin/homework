package queue

import (
	"encoding/json"
	"fmt"

	amqp "github.com/rabbitmq/amqp091-go"
)

type Producer struct {
	name string
	ch   *amqp.Channel
	conn *amqp.Connection
}

func NewProducer(name string, conn *amqp.Connection) Producer {
	return Producer{
		name: name,
		conn: conn,
	}
}

func (p *Producer) Start() error {
	var err error
	p.ch, err = p.conn.Channel()
	if err != nil {
		return fmt.Errorf("failed to open channel: %w", err)
	}

	_, err = p.ch.QueueDeclare(
		p.name,
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return fmt.Errorf("failed to declare queue: %w", err)
	}

	return nil
}

func (p *Producer) Stop() {
	if p.ch != nil {
		p.ch.Close()
	}
}

func (p Producer) Publish(body interface{}) error {
	b, err := json.Marshal(body)
	if err != nil {
		return fmt.Errorf("failed marshal body while publishing: %w", err)
	}
	err = p.ch.Publish(
		"",
		p.name,
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        b,
		},
	)
	if err != nil {
		return fmt.Errorf("failed publish: %w", err)
	}

	return nil
}
