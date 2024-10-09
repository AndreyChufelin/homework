package queue

import (
	"fmt"

	amqp "github.com/rabbitmq/amqp091-go"
)

type Queue struct {
	name     string
	password string
	host     string
	port     string
	Conn     *amqp.Connection
}

func NewQueue(name, password, host, port string) Queue {
	return Queue{
		name:     name,
		password: password,
		host:     host,
		port:     port,
	}
}

func (q *Queue) Start() error {
	var err error
	q.Conn, err = amqp.Dial(fmt.Sprintf("amqp://%s:%s@%s:%s/", q.name, q.password, q.host, q.port))
	if err != nil {
		return fmt.Errorf("failed connect queue: %w", err)
	}

	return nil
}

func (q *Queue) Stop() {
	if q.Conn != nil {
		q.Conn.Close()
	}
}
