package internalhttp

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/AndreyChufelin/homework/hw12_13_14_15_calendar/internal/logger"
	"github.com/AndreyChufelin/homework/hw12_13_14_15_calendar/internal/storage"
)

type Server struct {
	logger logger.Logger
	app    Application
	server *http.Server
	addr   string
}

//go:generate mockery --name=Application
type Application interface {
	CreateEvent(ctx context.Context, event storage.Event) error
	GetEvent(ctx context.Context, id string) (*storage.Event, error)
	DeleteEvent(ctx context.Context, id string) error
	EditEvent(ctx context.Context, id string, event storage.Event) error
	GetEventsListDay(ctx context.Context, date time.Time) ([]storage.Event, error)
	GetEventsListWeek(ctx context.Context, date time.Time) ([]storage.Event, error)
	GetEventsListMonth(ctx context.Context, date time.Time) ([]storage.Event, error)
}

type Logger interface {
	logger.Logger
}

func NewServer(logger Logger, app Application, host, port string) *Server {
	return &Server{logger: logger, app: app, addr: fmt.Sprintf("%s:%s", host, port)}
}

func (s *Server) Start() error {
	s.logger.Info("starting server")
	mux := http.NewServeMux()
	mux.Handle("/hello", loggingMiddleware(s.logger, http.HandlerFunc(s.hello)))
	mux.Handle("/event/create", loggingMiddleware(s.logger, methodHandler(http.MethodPost, s.createEventHandler)))
	mux.Handle("/event/delete/", loggingMiddleware(s.logger, methodHandler(http.MethodDelete, s.deleteEventHandler)))
	mux.Handle("/event/edit/", loggingMiddleware(s.logger, methodHandler(http.MethodPut, s.editEventHandler)))
	mux.Handle("/event/day/", loggingMiddleware(s.logger, methodHandler(http.MethodGet, s.getEventsDayHandler)))
	mux.Handle("/event/week/", loggingMiddleware(s.logger, methodHandler(http.MethodGet, s.getEventsWeekHandler)))
	mux.Handle("/event/month/", loggingMiddleware(s.logger, methodHandler(http.MethodGet, s.getEventsMonthHandler)))
	mux.Handle("/event/", loggingMiddleware(s.logger, methodHandler(http.MethodGet, s.getEventHandler)))

	s.server = &http.Server{
		Addr:              s.addr,
		Handler:           mux,
		ReadHeaderTimeout: 2 * time.Second,
	}

	if err := s.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return fmt.Errorf("server.Start: %w", err)
	}

	return nil
}

func (s *Server) Stop(ctx context.Context) error {
	s.logger.Info("stopping server")
	err := s.server.Shutdown(ctx)
	if err != nil {
		return fmt.Errorf("server.Stop: %w", err)
	}
	s.logger.Info("server stopped")

	return nil
}

func (s *Server) hello(res http.ResponseWriter, _ *http.Request) {
	time.Sleep(2 * time.Second)
	fmt.Fprintf(res, "Hello, world!")
}
