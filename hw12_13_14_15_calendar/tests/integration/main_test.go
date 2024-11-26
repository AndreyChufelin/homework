//go:build integration

package integration

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"testing"
	"time"

	pb "github.com/AndreyChufelin/homework/hw12_13_14_15_calendar/api"
	"github.com/AndreyChufelin/homework/hw12_13_14_15_calendar/internal/app"
	loggerslog "github.com/AndreyChufelin/homework/hw12_13_14_15_calendar/internal/logger/slog"
	internalgrpc "github.com/AndreyChufelin/homework/hw12_13_14_15_calendar/internal/server/grpc"
	"github.com/AndreyChufelin/homework/hw12_13_14_15_calendar/internal/storage"
	sqlstorage "github.com/AndreyChufelin/homework/hw12_13_14_15_calendar/internal/storage/sql"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/suite"
	"google.golang.org/grpc/status"
)

type IntegrationSuite struct {
	suite.Suite
	handlers grpcHandlers
}

type grpcHandlers interface {
	CreateEvent(context.Context, *pb.CreateEventRequest) (*pb.CreateEventResponse, error)
	GetEventsDay(context.Context, *pb.GetEventsDayRequest) (*pb.GetEventsDayResponse, error)
	GetEventsWeek(context.Context, *pb.GetEventsWeekRequest) (*pb.GetEventsWeekResponse, error)
	GetEventsMonth(context.Context, *pb.GetEventsMonthRequest) (*pb.GetEventsMonthResponse, error)
}

var (
	configFile = "./config.toml"
	db         *sqlx.DB
	config     Config
	logg       *loggerslog.Logger
	store      *sqlstorage.Storage
)

func TestMain(m *testing.M) {
	var err error
	config, err = LoadConfig(configFile)
	if err != nil {
		log.Fatalf("failed to read config from %q: %v", configFile, err)
	}
	logg, err = loggerslog.New(os.Stdout, config.Logger.Level)
	if err != nil {
		log.Fatalf("failed to create logger: %v", err)
	}

	db, err = sqlx.Connect("postgres",
		fmt.Sprintf(
			"user=%s dbname=%s sslmode=disable password=%s host=%s port=%s",
			config.DB.User,
			config.DB.Name,
			config.DB.Password,
			config.DB.Host,
			config.DB.Port,
		),
	)
	if err != nil {
		log.Fatalf("failed to connect to db: %v", err)
	}

	store = sqlstorage.New(config.DB.User, config.DB.Password, config.DB.Name, config.DB.Host, config.DB.Port)
	err = store.Connect(context.TODO())
	if err != nil {
		log.Fatalf("failed to connect storage: %v", err)
	}

	code := m.Run()

	db.Close()
	os.Exit(code)
}

func (s *IntegrationSuite) SetupSuite() {
	app := app.New(logg, store)
	s.handlers = internalgrpc.NewServer(logg, app, config.GRPC.Host, config.GRPC.Port)
}

func (s *IntegrationSuite) TearDownTest() {
	clearEvents()
}

func clearEvents() {
	_, err := db.Exec("TRUNCATE events")
	if err != nil {
		log.Fatalf("failed delete all events: %v", err)
	}
}

type eventSQL struct {
	ID                        string
	Title                     string
	Date                      time.Time
	EndDate                   time.Time `db:"end_date"`
	Description               sql.NullString
	UserID                    string                     `db:"user_id"`
	AdvanceNotificationPeriod sql.NullString             `db:"advance_notification_period"`
	NotificationStatus        storage.NotificationStatus `db:"notification_status"`
}

func (s *IntegrationSuite) TestCreateEvent() {
	now := time.Now()
	response, err := s.handlers.CreateEvent(context.TODO(), &pb.CreateEventRequest{
		Event: &pb.Event{
			Id:                        "",
			Title:                     "Test",
			Date:                      now.Unix(),
			EndDate:                   now.Unix(),
			UserId:                    "cf7ef14b-a43e-4449-a462-3b45620dca93",
			AdvanceNotificationPeriod: 86400,
		},
	})

	s.NoError(err)
	s.Nil(response)

	var events []eventSQL
	err = db.Select(&events, "SELECT * FROM events")
	s.NoError(err)

	s.Len(events, 1)
	s.Equal("Test", events[0].Title)
	s.Equal(now.Truncate(time.Second).UTC(), events[0].Date.UTC())
	s.Equal(now.Truncate(time.Second).UTC(), events[0].EndDate.UTC())
	s.Equal("cf7ef14b-a43e-4449-a462-3b45620dca93", events[0].UserID)
	s.Equal("24:00:00", events[0].AdvanceNotificationPeriod.String)
}

func (s *IntegrationSuite) validationErr(err error) {
	st := status.Convert(err)
	expect := map[string]string{
		"id":       "not valid uuid",
		"title":    "too long",
		"end_date": "too early",
		"user_id":  "not valid uuid",
	}

	var errs []*pb.BadRequest_FieldValiation
	for _, detail := range st.Details() {
		if ty, ok := detail.(*pb.BadRequest); ok {
			errs = ty.GetErrors()
			for _, violation := range errs {
				s.Equal(expect[violation.Field], violation.GetDescription(), "field %s", violation.Field)
			}
		}
	}
	s.Len(errs, len(expect))
}

func (s *IntegrationSuite) TestCreateEventValidation() {
	now := time.Now()
	response, err := s.handlers.CreateEvent(context.TODO(), &pb.CreateEventRequest{
		Event: &pb.Event{
			Id:                        "test",
			Title:                     "Toolong dfas fdsa gsgs gasfg sag asg asg asfgsafg sdgasd gsagd",
			Date:                      now.Add(time.Hour).Unix(),
			EndDate:                   now.Unix(),
			UserId:                    "cf7ef14b-a43e--a462-3b45620dca93",
			AdvanceNotificationPeriod: 86400,
		},
	})

	s.validationErr(err)
	s.Nil(response)

	var events []eventSQL
	err = db.Select(&events, "SELECT * FROM events")

	s.Len(events, 0)
}

func (s *IntegrationSuite) fillTable(events []storage.Event) {
	for _, e := range events {
		_, err := db.NamedExec(`INSERT INTO events 
		(title, date, end_date, description, user_id, advance_notification_period) 
		VALUES (:title, :date, :enddate, :description, :userid, :advancenotificationperiod)`, &e)
		s.NoError(err)
	}
}

func (s *IntegrationSuite) compareEvents(events []*pb.Event, createEvents []storage.Event) {
	s.Len(events, 1)
	s.Equal(createEvents[0].Title, events[0].Title)
	s.Equal(createEvents[0].Date.Unix(), events[0].Date)
	s.Equal(createEvents[0].EndDate.Unix(), events[0].EndDate)
	s.Equal(createEvents[0].UserID, events[0].UserId)
	s.Equal(int64(createEvents[0].AdvanceNotificationPeriod), events[0].AdvanceNotificationPeriod)
}

func (s *IntegrationSuite) TestGetEventsDay() {
	now := time.Now()
	createEvents := []storage.Event{
		{
			Title:                     "Today",
			Date:                      now.UTC(),
			EndDate:                   now.UTC(),
			UserID:                    "cf7ef14b-a43e-4449-a462-3b45620dca93",
			AdvanceNotificationPeriod: time.Duration(time.Hour.Seconds()),
		},
		{
			Title:                     "Tomorrow",
			Date:                      now.Add(48 * time.Hour).UTC(),
			EndDate:                   now.Add(48 * time.Hour).UTC(),
			UserID:                    "cf7ef14b-a43e-4449-a462-3b45620dca93",
			AdvanceNotificationPeriod: time.Duration(time.Hour.Seconds()),
		},
	}

	s.fillTable(createEvents)

	response, err := s.handlers.GetEventsDay(context.TODO(), &pb.GetEventsDayRequest{
		Date: now.Unix(),
	})

	s.NoError(err)
	s.compareEvents(response.Events, createEvents)
}

func (s *IntegrationSuite) TestGetEventsWeek() {
	now := time.Now()
	createEvents := []storage.Event{
		{
			Title:                     "This week",
			Date:                      now.AddDate(0, 0, 3).UTC(),
			EndDate:                   now.AddDate(0, 0, 3).UTC(),
			UserID:                    "cf7ef14b-a43e-4449-a462-3b45620dca93",
			AdvanceNotificationPeriod: time.Duration(time.Hour.Seconds()),
		},
		{
			Title:                     "Next week",
			Date:                      now.AddDate(0, 0, 7).UTC(),
			EndDate:                   now.AddDate(0, 0, 7).UTC(),
			UserID:                    "cf7ef14b-a43e-4449-a462-3b45620dca93",
			AdvanceNotificationPeriod: time.Duration(time.Hour.Seconds()),
		},
	}

	s.fillTable(createEvents)

	response, err := s.handlers.GetEventsWeek(context.TODO(), &pb.GetEventsWeekRequest{
		Date: now.Unix(),
	})

	s.NoError(err)
	s.compareEvents(response.Events, createEvents)
}

func (s *IntegrationSuite) TestGetEventsMonth() {
	now := time.Now()
	createEvents := []storage.Event{
		{
			Title:                     "This month",
			Date:                      now.AddDate(0, 0, 3).UTC(),
			EndDate:                   now.AddDate(0, 0, 3).UTC(),
			UserID:                    "cf7ef14b-a43e-4449-a462-3b45620dca93",
			AdvanceNotificationPeriod: time.Duration(time.Hour.Seconds()),
		},
		{
			Title:                     "Next month",
			Date:                      now.AddDate(0, 2, 0).UTC(),
			EndDate:                   now.AddDate(0, 2, 0).UTC(),
			UserID:                    "cf7ef14b-a43e-4449-a462-3b45620dca93",
			AdvanceNotificationPeriod: time.Duration(time.Hour.Seconds()),
		},
	}

	s.fillTable(createEvents)

	response, err := s.handlers.GetEventsMonth(context.TODO(), &pb.GetEventsMonthRequest{
		Date: now.Unix(),
	})

	s.NoError(err)
	s.compareEvents(response.Events, createEvents)
}

func TestIntegrationSuite(t *testing.T) {
	suite.Run(t, new(IntegrationSuite))
	suite.Run(t, new(NotificationSuite))
}
