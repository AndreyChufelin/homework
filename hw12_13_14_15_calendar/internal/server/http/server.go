package internalhttp

import (
	"context"
	"fmt"
	"net/http"
	"time"
)

type Server struct {
	logger Logger
	app    Application
	server *http.Server
	addr   string
}

type Logger interface {
	Info(msg string, keysAndValues ...interface{})
}

type Application interface { // TODO
}

func NewServer(logger Logger, app Application, host, port string) *Server {
	return &Server{logger: logger, app: app, addr: fmt.Sprintf("%s:%s", host, port)}
}

func (s *Server) Start() error {
	s.logger.Info("starting server")
	mux := http.NewServeMux()
	mux.Handle("/hello", loggingMiddleware(s.logger, http.HandlerFunc(s.hello)))

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
