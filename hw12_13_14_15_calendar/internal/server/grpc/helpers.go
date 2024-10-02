package internalgrpc

import (
	"context"
	"errors"
	"fmt"
	"time"

	pb "github.com/AndreyChufelin/homework/hw12_13_14_15_calendar/api"
	"github.com/AndreyChufelin/homework/hw12_13_14_15_calendar/internal/storage"
	"github.com/AndreyChufelin/homework/hw12_13_14_15_calendar/internal/validator"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func badRequestError(errors map[string]string) error {
	st := status.New(codes.InvalidArgument, "validation error")
	badRequest := &pb.BadRequest{}
	for field, description := range errors {
		badRequest.Errors = append(badRequest.Errors, &pb.BadRequest_FieldValiation{
			Field:       field,
			Description: description,
		})
	}
	st, err := st.WithDetails(badRequest)
	if err != nil {
		panic(fmt.Sprintf("Unexpected error attaching metadata: %v", err))
	}

	return st.Err()
}

func protoToEvent(event *pb.Event) storage.Event {
	return storage.Event{
		ID:                        event.Id,
		Title:                     event.Title,
		Date:                      time.Unix(event.Date, 0).UTC(),
		EndDate:                   time.Unix(event.EndDate, 0).UTC(),
		UserID:                    event.UserId,
		AdvanceNotificationPeriod: time.Duration(event.AdvanceNotificationPeriod) * time.Second,
	}
}

func eventToProto(event *storage.Event) *pb.Event {
	return &pb.Event{
		Id:                        event.ID,
		Title:                     event.Title,
		Date:                      event.Date.Unix(),
		EndDate:                   event.EndDate.Unix(),
		UserId:                    event.UserID,
		AdvanceNotificationPeriod: int64(event.AdvanceNotificationPeriod.Seconds()),
	}
}

func prepareEvent(event *pb.Event) (storage.Event, error) {
	e := protoToEvent(event)

	validator := validator.New()
	storage.ValidateEvent(*validator, e)
	if !validator.Valid() {
		return storage.Event{}, badRequestError(validator.Errors)
	}

	return e, nil
}

func getEventsDate(
	ctx context.Context,
	logger Logger,
	date int64,
	cb func(context.Context, time.Time) ([]storage.Event, error),
) ([]*pb.Event, error) {
	events, err := cb(ctx, time.Unix(date, 0).UTC())
	if err != nil {
		if errors.Is(err, storage.ErrNoEventsFound) {
			logger.Warn("no events found")
			return nil, status.Error(codes.NotFound, "no events found")
		}
		logger.Error("failed get day events", "error", err)
		return nil, status.Error(codes.Internal, err.Error())
	}

	resopnse := make([]*pb.Event, len(events))
	for i, event := range events {
		e := event
		resopnse[i] = eventToProto(&e)
	}

	return resopnse, nil
}
