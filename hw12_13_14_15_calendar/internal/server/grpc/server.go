package internalgrpc

import (
	"context"
	"fmt"
	"net"
	"time"

	pb "github.com/AndreyChufelin/homework/hw12_13_14_15_calendar/api"
	"github.com/AndreyChufelin/homework/hw12_13_14_15_calendar/internal/logger"
	"github.com/AndreyChufelin/homework/hw12_13_14_15_calendar/internal/storage"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
	"google.golang.org/grpc/status"
)

type Server struct {
	pb.UnimplementedCalendarServer
	logger logger.Logger
	app    Application
	server *grpc.Server
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
	grpcServer := grpc.NewServer(grpc.UnaryInterceptor(LoggingInterceptor(logger)))
	return &Server{logger: logger, app: app, addr: fmt.Sprintf("%s:%s", host, port), server: grpcServer}
}

func (s *Server) Start() error {
	s.logger.Info("starting grpc server")
	l, err := net.Listen("tcp", s.addr)
	if err != nil {
		return fmt.Errorf("grpc.Start: %w", err)
	}

	pb.RegisterCalendarServer(s.server, s)

	if err := s.server.Serve(l); err != nil {
		return fmt.Errorf("grpc.Start: %w", err)
	}

	return nil
}

func (s *Server) Stop() {
	if s.server == nil {
		return
	}

	s.logger.Info("stopping grpc server")
	s.server.GracefulStop()
}

func LoggingInterceptor(logger Logger) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		date := time.Now()
		resp, err := handler(ctx, req)

		latency := time.Since(date)

		p, ok := peer.FromContext(ctx)
		ip := "unknown"
		if ok {
			ip = p.Addr.String()
		}

		md, ok := metadata.FromIncomingContext(ctx)
		userAgent := "unknown"
		if ok {
			if userAgents, exists := md["user-agent"]; exists && len(userAgents) > 0 {
				userAgent = userAgents[0]
			}
		}

		logger.Info("GRPC request handled",
			"ip", ip,
			"method", info.FullMethod,
			"date", date.Format(time.RFC822Z),
			"userAgent", userAgent,
			"latency", latency.Milliseconds(),
			"status", status.Code(err).String(),
		)
		return resp, err
	}
}
