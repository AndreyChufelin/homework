package internalgrpc

import (
	"context"
	"errors"

	pb "github.com/AndreyChufelin/homework/hw12_13_14_15_calendar/api"
	"github.com/AndreyChufelin/homework/hw12_13_14_15_calendar/internal/storage"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Server) CreateEvent(ctx context.Context, request *pb.CreateEventRequest) (*pb.CreateEventResponse, error) {
	logg := s.logger.With("handler", "createEventHandler")
	event, err := prepareEvent(request.GetEvent())
	if err != nil {
		logg.Warn("failed to prepare event", "error", err)
		return nil, err
	}

	err = s.app.CreateEvent(ctx, event)
	if err != nil {
		if errors.Is(err, storage.ErrEventAlreadyExists) {
			logg.Warn("event already exists")
			return nil, status.Error(codes.AlreadyExists, "event already exists")
		}
		logg.Error("failed create event", "error", err)
		return nil, status.Error(codes.Internal, err.Error())
	}
	return nil, nil
}

func (s *Server) GetEvent(ctx context.Context, request *pb.GetEventRequest) (*pb.GetEventResponse, error) {
	logg := s.logger.With("handler", "getEventHandler")
	event, err := s.app.GetEvent(ctx, request.Id)
	if err != nil {
		if errors.Is(err, storage.ErrEventDoesntExist) {
			logg.Warn("event doesn't exist")
			return nil, status.Error(codes.NotFound, "event doesn't exist")
		}
		logg.Error("failed get event", "error", err)
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &pb.GetEventResponse{Event: eventToProto(event)}, nil
}

func (s *Server) EditEvent(ctx context.Context, request *pb.EditEventRequest) (*pb.EditEventResponse, error) {
	logg := s.logger.With("handler", "editEventHandler")
	event, err := prepareEvent(request.GetEvent())
	if err != nil {
		return nil, err
	}

	err = s.app.EditEvent(ctx, request.Id, event)
	if err != nil {
		if errors.Is(err, storage.ErrEventDoesntExist) {
			logg.Warn("event doesn't exist")
			return nil, status.Error(codes.NotFound, "event doesn't exist")
		}
		logg.Error("failed edit event", "error", err)
		return nil, status.Error(codes.Internal, err.Error())
	}
	return nil, nil
}

func (s *Server) DeleteEvent(ctx context.Context, request *pb.DeleteEventRequest) (*pb.DeleteEventResponse, error) {
	logg := s.logger.With("handler", "createEventHandler")
	err := s.app.DeleteEvent(ctx, request.Id)
	if err != nil {
		if errors.Is(err, storage.ErrEventDoesntExist) {
			logg.Warn("event doesn't exist")
			return nil, status.Error(codes.NotFound, "event doesn't exist")
		}
		logg.Error("failed delete event", "error", err)
		return nil, status.Error(codes.Internal, err.Error())
	}
	return nil, nil
}

func (s *Server) GetEventsDay(ctx context.Context, request *pb.GetEventsDayRequest) (*pb.GetEventsDayResponse, error) {
	logg := s.logger.With("handler", "getEventsDayEventHandler")
	events, err := getEventsDate(ctx, logg, request.Date, s.app.GetEventsListDay)
	if err != nil {
		return nil, err
	}
	return &pb.GetEventsDayResponse{Events: events}, nil
}

func (s *Server) GetEventsWeek(
	ctx context.Context,
	request *pb.GetEventsWeekRequest,
) (*pb.GetEventsWeekResponse, error) {
	logg := s.logger.With("handler", "getEventsWeekEventHandler")
	events, err := getEventsDate(ctx, logg, request.Date, s.app.GetEventsListWeek)
	if err != nil {
		return nil, err
	}
	return &pb.GetEventsWeekResponse{Events: events}, nil
}

func (s *Server) GetEventsMonth(
	ctx context.Context,
	request *pb.GetEventsMonthRequest,
) (*pb.GetEventsMonthResponse, error) {
	logg := s.logger.With("handler", "getEventsMonthEventHandler")
	events, err := getEventsDate(ctx, logg, request.Date, s.app.GetEventsListMonth)
	if err != nil {
		return nil, err
	}
	return &pb.GetEventsMonthResponse{Events: events}, nil
}
