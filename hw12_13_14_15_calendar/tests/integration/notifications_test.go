//go:build integration

package integration

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/AndreyChufelin/homework/hw12_13_14_15_calendar/internal/queue"
	"github.com/AndreyChufelin/homework/hw12_13_14_15_calendar/internal/scheduler"
	"github.com/AndreyChufelin/homework/hw12_13_14_15_calendar/internal/sender"
	"github.com/AndreyChufelin/homework/hw12_13_14_15_calendar/internal/storage"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/stretchr/testify/suite"
)

type NotificationSuite struct {
	suite.Suite
	queue    queue.Queue
	producer queue.Producer
	consumer queue.Consumer
	channel  *amqp.Channel
}

const queueName = "test_queue"

func (s *NotificationSuite) SetupSuite() {
	conn, err := amqp.Dial(fmt.Sprintf("amqp://%s:%s@%s:%s/", config.Queue.User, config.Queue.Password, config.Queue.Host, config.Queue.Port))
	if err != nil {
		log.Fatalf("failed to connect queue: %v", err)
	}
	s.channel, err = conn.Channel()
	if err != nil {
		log.Fatalf("failed to connect channel queue: %v", err)
	}

	s.queue = queue.NewQueue(config.Queue.User, config.Queue.Password, config.Queue.Host, config.Queue.Port)
	err = s.queue.Start()
	if err != nil {
		log.Fatal("failed to start queue", err)
	}

	s.producer = queue.NewProducer(queueName, s.queue.Conn)
	err = s.producer.Start()
	if err != nil {
		logg.Error("failed to start producer", "err", err)
	}

	go func() {
		sch := scheduler.NewScheduler(s.producer, config.ClearInterval, config.Interval, logg, store)
		sch.Start(context.TODO())
	}()

	s.consumer = queue.NewConsumer(queueName, s.queue.Conn)
	err = s.consumer.Start()
	if err != nil {
		log.Fatal("failed start consumer", err)
	}

	go func() {
		sen := sender.NewSender(s.consumer, logg, store)
		err = sen.Start()
		if err != nil {
			log.Fatal("failed start sender", err)
		}
	}()
}

func (s *NotificationSuite) TearDownSuite() {
	s.queue.Stop()
	s.producer.Stop()
	s.consumer.Stop()
}

func (s *NotificationSuite) TearDownTest() {
	clearEvents()
	_, err := s.channel.QueuePurge(queueName, false)
	if err != nil {
		log.Fatalf("failed purge queue: %v", err)
	}
}

func (s *NotificationSuite) TestNotificationSends() {
	e := storage.Event{
		Title:                     "test",
		Date:                      time.Now(),
		EndDate:                   time.Now(),
		UserID:                    "fd3195e5-17a9-4b61-8d9d-0d1bbb4edf93",
		AdvanceNotificationPeriod: time.Duration((24 * time.Hour).Seconds()),
	}
	_, err := db.NamedExec(`INSERT INTO events 
		(title, date, end_date, description, user_id, advance_notification_period) 
		VALUES (:title, :date, :enddate, :description, :userid, :advancenotificationperiod)`, &e)
	s.NoError(err)

	wait := 3 * time.Duration(config.Interval) * time.Second
	s.Eventually(func() bool {
		var events []eventSQL
		err = db.Select(&events, "SELECT * FROM events")
		s.NoError(err)

		s.Require().Len(events, 1)
		return events[0].NotificationStatus == storage.StatusSent
	}, wait, time.Second)
}

func (s *NotificationSuite) TestEventClears() {
	e := storage.Event{
		Title:                     "test",
		Date:                      time.Now().AddDate(-1, 0, 0),
		EndDate:                   time.Now(),
		UserID:                    "fd3195e5-17a9-4b61-8d9d-0d1bbb4edf93",
		AdvanceNotificationPeriod: time.Duration((24 * time.Hour).Seconds()),
	}
	_, err := db.NamedExec(`INSERT INTO events 
		(title, date, end_date, description, user_id, advance_notification_period) 
		VALUES (:title, :date, :enddate, :description, :userid, :advancenotificationperiod)`, &e)
	s.NoError(err)

	wait := 3 * time.Duration(config.Interval) * time.Second
	s.Eventually(func() bool {
		var events []eventSQL
		err = db.Select(&events, "SELECT * FROM events")
		s.NoError(err)
		s.T().Log("events", events)

		return len(events) == 0
	}, wait, time.Second)
}
