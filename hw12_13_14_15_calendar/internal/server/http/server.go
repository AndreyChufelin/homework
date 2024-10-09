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
	s := &Server{logger: logger, app: app, addr: fmt.Sprintf("%s:%s", host, port)}
	mux := http.NewServeMux()
	mux.Handle("/hello", loggingMiddleware(s.logger, http.HandlerFunc(s.hello)))
	mux.Handle("POST /event/create", loggingMiddleware(s.logger, http.HandlerFunc(s.createEventHandler)))
	mux.Handle("DELETE /event/delete/{id}", loggingMiddleware(s.logger, http.HandlerFunc(s.deleteEventHandler)))
	mux.Handle("PUT /event/edit/{id}", loggingMiddleware(s.logger, http.HandlerFunc(s.editEventHandler)))
	mux.Handle("GET /event/day/{date}", loggingMiddleware(s.logger, http.HandlerFunc(s.getEventsDayHandler)))
	mux.Handle("GET /event/week/{date}", loggingMiddleware(s.logger, http.HandlerFunc(s.getEventsWeekHandler)))
	mux.Handle("GET /event/month/{date}", loggingMiddleware(s.logger, http.HandlerFunc(s.getEventsMonthHandler)))
	mux.Handle("GET /event/{id}", loggingMiddleware(s.logger, http.HandlerFunc(s.getEventHandler)))

	s.server = &http.Server{
		Addr:              s.addr,
		Handler:           mux,
		ReadHeaderTimeout: 2 * time.Second,
	}

	return s
}

func (s *Server) Start() error {
	s.logger.Info("starting server")

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
