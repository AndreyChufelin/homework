package internalgrpc

import (
	"context"
	"errors"
	"io"
	"testing"
	"time"

	pb "github.com/AndreyChufelin/homework/hw12_13_14_15_calendar/api"
	logger "github.com/AndreyChufelin/homework/hw12_13_14_15_calendar/internal/logger/slog"
	"github.com/AndreyChufelin/homework/hw12_13_14_15_calendar/internal/server/grpc/mocks"
	"github.com/AndreyChufelin/homework/hw12_13_14_15_calendar/internal/storage"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func newLogger(t *testing.T) *logger.Logger {
	t.Helper()
	logg, err := logger.New(io.Discard, "INFO")
	if err != nil {
		t.Fatal(err)
	}

	return logg
}

var (
	userID       = "66be96d3-3d5d-4aec-af9c-5b3769d0169a"
	eventID      = "01924888-c5a8-74c5-bf47-c87787247388"
	eventMessage = pb.Event{
		Id:                        eventID,
		Title:                     "test",
		Date:                      time.Now().Unix(),
		EndDate:                   time.Now().Add(time.Hour).Unix(),
		UserId:                    userID,
		AdvanceNotificationPeriod: 0,
	}

	eventStorage = storage.Event{
		ID:                        eventID,
		Title:                     "test",
		Date:                      time.Now().Truncate(time.Second).UTC(),
		EndDate:                   time.Now().Truncate(time.Second).UTC().Add(time.Hour),
		UserID:                    userID,
		AdvanceNotificationPeriod: 0,
	}
)

func validationErr(t *testing.T, err error) {
	t.Helper()
	st := status.Convert(err)
	t.Log(st.Details())
	expect := map[string]string{
		"id":       "not valid uuid",
		"title":    "too long",
		"end_date": "too early",
		"user_id":  "not valid uuid",
	}

	for _, detail := range st.Details() {
		if ty, ok := detail.(*pb.BadRequest); ok {
			for _, violation := range ty.GetErrors() {
				assert.Equal(t, expect[violation.Field], violation.GetDescription(), "field %s", violation.Field)
			}
		}
	}
}

func TestCreateEvent(t *testing.T) {
	tests := []struct {
		name      string
		returns   []interface{}
		event     *pb.Event
		wantEvent storage.Event
		err       error
	}{
		{
			name: "success",
			returns: []interface{}{
				nil,
			},
			event:     &eventMessage,
			wantEvent: eventStorage,
		},
		{
			name: "event already exists",
			returns: []interface{}{
				storage.ErrEventAlreadyExists,
			},
			event:     &eventMessage,
			wantEvent: eventStorage,
			err:       status.Error(codes.AlreadyExists, "event already exists"),
		},
		{
			name: "internal error",
			returns: []interface{}{
				errors.New("internal error"),
			},
			event:     &eventMessage,
			wantEvent: eventStorage,
			err:       status.Error(codes.Internal, "Internal server error"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			logg := newLogger(t)
			app := mocks.NewApplication(t)
			app.On("CreateEvent", mock.Anything, tt.wantEvent).Return(tt.returns...)
			server := NewServer(logg, app, "", "")

			_, err := server.CreateEvent(context.TODO(), &pb.CreateEventRequest{Event: tt.event})

			if tt.err == nil {
				require.NoError(t, err)
			} else {
				require.ErrorIs(t, err, tt.err)
			}
		})
	}
}

func TestCreateEventVaidation(t *testing.T) {
	logg := newLogger(t)
	app := mocks.NewApplication(t)
	server := NewServer(logg, app, "", "")

	_, err := server.CreateEvent(context.TODO(), &pb.CreateEventRequest{Event: &pb.Event{
		Id:                        "-c5a8-74c5-bf47-c87787247388",
		Title:                     "testtesttesttesttesttesttesttesttesttest",
		Date:                      time.Now().Unix(),
		EndDate:                   time.Now().Add(-time.Hour).Unix(),
		UserId:                    "-c5a8-74c5-bf47-c87787247388",
		AdvanceNotificationPeriod: 0,
	}})

	validationErr(t, err)
}

func TestGetEvent(t *testing.T) {
	tests := []struct {
		name    string
		returns []interface{}
		want    *pb.GetEventResponse
		err     error
	}{
		{
			name: "success",
			returns: []interface{}{
				&eventStorage,
				nil,
			},
			want: &pb.GetEventResponse{
				Event: &eventMessage,
			},
		},
		{
			name: "event doesn't exist",
			returns: []interface{}{
				nil,
				storage.ErrEventDoesntExist,
			},
			err: status.Error(codes.NotFound, "event doesn't exist"),
		},
		{
			name: "internal error",
			returns: []interface{}{
				nil,
				errors.New("internal error"),
			},
			err: status.Error(codes.Internal, "Internal server error"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			logg := newLogger(t)
			app := mocks.NewApplication(t)
			app.On("GetEvent", mock.Anything, eventID).Return(tt.returns...)
			server := NewServer(logg, app, "", "")

			res, err := server.GetEvent(context.TODO(), &pb.GetEventRequest{Id: eventID})

			if tt.err == nil {
				require.NoError(t, err)
				require.Equal(t, tt.want, res)
			} else {
				require.ErrorIs(t, err, tt.err)
				require.Nil(t, res)
			}
		})
	}
}

func TestEditEvent(t *testing.T) {
	tests := []struct {
		name      string
		returns   []interface{}
		event     *pb.Event
		wantEvent storage.Event
		err       error
	}{
		{
			name: "success",
			returns: []interface{}{
				nil,
			},
			event:     &eventMessage,
			wantEvent: eventStorage,
		},
		{
			name: "event doesn't exists",
			returns: []interface{}{
				storage.ErrEventDoesntExist,
			},
			event:     &eventMessage,
			wantEvent: eventStorage,
			err:       status.Error(codes.NotFound, "event doesn't exist"),
		},
		{
			name: "internal error",
			returns: []interface{}{
				errors.New("internal error"),
			},
			event:     &eventMessage,
			wantEvent: eventStorage,
			err:       status.Error(codes.Internal, "Internal server error"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			logg := newLogger(t)
			app := mocks.NewApplication(t)
			app.On("EditEvent", mock.Anything, eventID, tt.wantEvent).Return(tt.returns...)
			server := NewServer(logg, app, "", "")

			_, err := server.EditEvent(context.TODO(), &pb.EditEventRequest{Id: eventID, Event: tt.event})

			if tt.err == nil {
				require.NoError(t, err)
			} else {
				require.ErrorIs(t, err, tt.err)
			}
		})
	}
}

func TestEditEventVaidation(t *testing.T) {
	logg := newLogger(t)
	app := mocks.NewApplication(t)
	server := NewServer(logg, app, "", "")

	_, err := server.EditEvent(context.TODO(), &pb.EditEventRequest{Id: eventID, Event: &pb.Event{
		Id:                        "-c5a8-74c5-bf47-c87787247388",
		Title:                     "testtesttesttesttesttesttesttesttesttest",
		Date:                      time.Now().Unix(),
		EndDate:                   time.Now().Add(-time.Hour).Unix(),
		UserId:                    "-c5a8-74c5-bf47-c87787247388",
		AdvanceNotificationPeriod: 0,
	}})

	validationErr(t, err)
}

func TestDeleteEvent(t *testing.T) {
	tests := []struct {
		name      string
		returns   []interface{}
		event     *pb.Event
		wantEvent storage.Event
		err       error
	}{
		{
			name: "success",
			returns: []interface{}{
				nil,
			},
			event:     &eventMessage,
			wantEvent: eventStorage,
		},
		{
			name: "event doesn't exists",
			returns: []interface{}{
				storage.ErrEventDoesntExist,
			},
			event:     &eventMessage,
			wantEvent: eventStorage,
			err:       status.Error(codes.NotFound, "event doesn't exist"),
		},
		{
			name: "internal error",
			returns: []interface{}{
				errors.New("internal error"),
			},
			event:     &eventMessage,
			wantEvent: eventStorage,
			err:       status.Error(codes.Internal, "Internal server error"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			logg := newLogger(t)
			app := mocks.NewApplication(t)
			app.On("DeleteEvent", mock.Anything, eventID).Return(tt.returns...)
			server := NewServer(logg, app, "", "")

			_, err := server.DeleteEvent(context.TODO(), &pb.DeleteEventRequest{Id: eventID})

			if tt.err == nil {
				require.NoError(t, err)
			} else {
				require.ErrorIs(t, err, tt.err)
			}
		})
	}
}

var testsDate = []struct {
	name    string
	returns []interface{}
	want    *pb.GetEventsDayResponse
	err     error
}{
	{
		name: "success",
		returns: []interface{}{
			[]storage.Event{eventStorage},
			nil,
		},
		want: &pb.GetEventsDayResponse{
			Events: []*pb.Event{&eventMessage},
		},
	},
	{
		name: "no events found",
		returns: []interface{}{
			nil,
			storage.ErrNoEventsFound,
		},
		err: status.Error(codes.NotFound, "no events found"),
	},
	{
		name: "internal error",
		returns: []interface{}{
			nil,
			errors.New("internal error"),
		},
		err: status.Error(codes.Internal, "Internal server error"),
	},
}

func TestGetEventsDay(t *testing.T) {
	for _, tt := range testsDate {
		t.Run(tt.name, func(t *testing.T) {
			logg := newLogger(t)
			app := mocks.NewApplication(t)
			app.On("GetEventsListDay", mock.Anything, time.Now().Truncate(time.Second).UTC()).Return(tt.returns...)
			server := NewServer(logg, app, "", "")

			res, err := server.GetEventsDay(context.TODO(), &pb.GetEventsDayRequest{Date: time.Now().Unix()})

			if tt.err == nil {
				require.NoError(t, err)
				require.Equal(t, &pb.GetEventsDayResponse{Events: []*pb.Event{&eventMessage}}, res)
			} else {
				require.ErrorIs(t, err, tt.err)
				require.Nil(t, res)
			}
		})
	}
}

func TestGetEventsWeek(t *testing.T) {
	for _, tt := range testsDate {
		t.Run(tt.name, func(t *testing.T) {
			logg := newLogger(t)
			app := mocks.NewApplication(t)
			app.On("GetEventsListWeek", mock.Anything, time.Now().Truncate(time.Second).UTC()).Return(tt.returns...)
			server := NewServer(logg, app, "", "")

			res, err := server.GetEventsWeek(context.TODO(), &pb.GetEventsWeekRequest{Date: time.Now().Unix()})

			if tt.err == nil {
				require.NoError(t, err)
				require.Equal(t, &pb.GetEventsWeekResponse{Events: []*pb.Event{&eventMessage}}, res)
			} else {
				require.ErrorIs(t, err, tt.err)
				require.Nil(t, res)
			}
		})
	}
}

func TestGetEventsMonth(t *testing.T) {
	for _, tt := range testsDate {
		t.Run(tt.name, func(t *testing.T) {
			logg := newLogger(t)
			app := mocks.NewApplication(t)
			app.On("GetEventsListMonth", mock.Anything, time.Now().Truncate(time.Second).UTC()).Return(tt.returns...)
			server := NewServer(logg, app, "", "")

			res, err := server.GetEventsMonth(context.TODO(), &pb.GetEventsMonthRequest{Date: time.Now().Unix()})

			if tt.err == nil {
				require.NoError(t, err)
				require.Equal(t, &pb.GetEventsMonthResponse{Events: []*pb.Event{&eventMessage}}, res)
			} else {
				require.ErrorIs(t, err, tt.err)
				require.Nil(t, res)
			}
		})
	}
}
