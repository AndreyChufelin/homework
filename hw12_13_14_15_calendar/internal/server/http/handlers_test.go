package internalhttp

import (
	"bytes"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	logger "github.com/AndreyChufelin/homework/hw12_13_14_15_calendar/internal/logger/slog"
	"github.com/AndreyChufelin/homework/hw12_13_14_15_calendar/internal/server/http/mocks"
	"github.com/AndreyChufelin/homework/hw12_13_14_15_calendar/internal/storage"
	mock "github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func newLogger(t *testing.T) *logger.Logger {
	t.Helper()
	logg, err := logger.New(io.Discard, "INFO")
	if err != nil {
		t.Fatal(err)
	}

	return logg
}

func TestGetEventHandler(t *testing.T) {
	tests := []struct {
		name    string
		returns []interface{}
		want    string
		status  int
	}{
		{
			name:    "success",
			returns: []interface{}{&storage.Event{ID: "1", Title: "test"}, nil},
			status:  http.StatusOK,
			want: `{
	"event": {
		"id": "1",
		"title": "test",
		"date": "0001-01-01T00:00:00Z",
		"end_date": "0001-01-01T00:00:00Z",
		"description": "",
		"user_id": "",
		"advance_notification_period": 0
	}
}`,
		},
		{
			name:    "not found",
			returns: []interface{}{nil, storage.ErrEventDoesntExist},
			status:  http.StatusNotFound,
			want: `{
	"error": "Event not found"
}`,
		},
		{
			name:    "internal error",
			returns: []interface{}{nil, errors.New("internal error")},
			status:  http.StatusInternalServerError,
			want: `{
	"error": "Unknown error"
}`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/events/1", nil)
			w := httptest.NewRecorder()

			logg := newLogger(t)
			app := mocks.NewApplication(t)
			app.On("GetEvent", mock.Anything, "1").Return(tt.returns...)
			server := NewServer(logg, app, "", "")

			server.getEventHandler(w, req)

			require.Equal(t, tt.status, w.Code)
			require.Equal(t, tt.want, w.Body.String())
		})
	}
}

func TestCreateEventHandler(t *testing.T) {
	eventArg := storage.Event{
		ID:                        "66be96d3-3d5d-4aec-af9c-5b3769d0169a",
		Title:                     "test",
		UserID:                    "66be96d3-3d5d-4aec-af9c-5b3769d0169a",
		Date:                      time.Date(2024, time.September, 23, 0, 0, 0, 0, time.UTC),
		EndDate:                   time.Date(2024, time.September, 25, 0, 0, 0, 0, time.UTC),
		Description:               "",
		AdvanceNotificationPeriod: 0,
	}
	event := `{
	"id": "66be96d3-3d5d-4aec-af9c-5b3769d0169a",
	"title": "test",
	"user_id": "66be96d3-3d5d-4aec-af9c-5b3769d0169a",
	"date": "2024-09-23T00:00:00Z",
	"end_date": "2024-09-25T00:00:00Z"
}`
	tests := []struct {
		name    string
		body    string
		returns []interface{}
		want    string
		status  int
	}{
		{
			name:    "success",
			body:    event,
			returns: []interface{}{nil},
			status:  http.StatusOK,
			want: `{
	"message": "Success"
}`,
		},
		{
			name:    "event already exist",
			body:    event,
			returns: []interface{}{storage.ErrEventAlreadyExists},
			status:  http.StatusConflict,
			want: `{
	"error": "Event already exist"
}`,
		},
		{
			name:    "event already exist",
			body:    event,
			returns: []interface{}{errors.New("internal error")},
			status:  http.StatusInternalServerError,
			want: `{
	"error": "Unknown error"
}`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/events/create", bytes.NewBufferString(tt.body))
			w := httptest.NewRecorder()

			logg := newLogger(t)
			app := mocks.NewApplication(t)
			app.On("CreateEvent", mock.Anything, eventArg).Return(tt.returns...)
			server := NewServer(logg, app, "", "")

			server.createEventHandler(w, req)

			require.Equal(t, tt.status, w.Code)
			require.Equal(t, tt.want, w.Body.String())
		})
	}
}

func TestCreateEventHandlerValidation(t *testing.T) {
	body := `{
	"id": "66be96d3-3d5d-4aec",
	"title": "testtitletoolongtesttitletoolong",
	"user_id": "66be96d3-3d5d-4aec",
	"date": "2024-09-23T00:00:00Z",
	"end_date": "2024-09-20T00:00:00Z"
}`
	req := httptest.NewRequest(http.MethodGet, "/events/create", bytes.NewBufferString(body))
	w := httptest.NewRecorder()

	logg := newLogger(t)
	app := mocks.NewApplication(t)
	server := NewServer(logg, app, "", "")

	server.createEventHandler(w, req)

	require.Equal(t, http.StatusPartialContent, w.Code)
	require.Equal(t, `{
	"error": {
		"end_date": "too early",
		"id": "not valid uuid",
		"title": "too long",
		"user_id": "not valid uuid"
	}
}`, w.Body.String())
}

func TestDeleteEventHandler(t *testing.T) {
	tests := []struct {
		name    string
		returns []interface{}
		want    string
		status  int
	}{
		{
			name:    "success",
			returns: []interface{}{nil},
			status:  http.StatusOK,
			want: `{
	"message": "Success"
}`,
		},
		{
			name:    "not found",
			returns: []interface{}{storage.ErrEventDoesntExist},
			status:  http.StatusNotFound,
			want: `{
	"error": "Event not found"
}`,
		},
		{
			name:    "internal error",
			returns: []interface{}{errors.New("internal error")},
			status:  http.StatusInternalServerError,
			want: `{
	"error": "Unknown error"
}`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/events/delete/1", nil)
			w := httptest.NewRecorder()

			logg := newLogger(t)
			app := mocks.NewApplication(t)
			app.On("DeleteEvent", mock.Anything, "1").Return(tt.returns...)
			server := NewServer(logg, app, "", "")

			server.deleteEventHandler(w, req)

			require.Equal(t, tt.status, w.Code)
			require.Equal(t, tt.want, w.Body.String())
		})
	}
}

func TestEditEventHandler(t *testing.T) {
	eventArg := storage.Event{
		ID:                        "66be96d3-3d5d-4aec-af9c-5b3769d0169a",
		Title:                     "test",
		UserID:                    "66be96d3-3d5d-4aec-af9c-5b3769d0169a",
		Date:                      time.Date(2024, time.September, 23, 0, 0, 0, 0, time.UTC),
		EndDate:                   time.Date(2024, time.September, 25, 0, 0, 0, 0, time.UTC),
		Description:               "",
		AdvanceNotificationPeriod: 0,
	}
	event := `{
	"id": "66be96d3-3d5d-4aec-af9c-5b3769d0169a",
	"title": "test",
	"user_id": "66be96d3-3d5d-4aec-af9c-5b3769d0169a",
	"date": "2024-09-23T00:00:00Z",
	"end_date": "2024-09-25T00:00:00Z"
}`
	tests := []struct {
		name    string
		body    string
		returns []interface{}
		want    string
		status  int
	}{
		{
			name:    "success",
			body:    event,
			returns: []interface{}{nil},
			status:  http.StatusOK,
			want: `{
	"message": "Success"
}`,
		},
		{
			name:    "not found",
			body:    event,
			returns: []interface{}{storage.ErrEventDoesntExist},
			status:  http.StatusNotFound,
			want: `{
	"error": "Event not found"
}`,
		},
		{
			name:    "event already exist",
			body:    event,
			returns: []interface{}{errors.New("internal error")},
			status:  http.StatusInternalServerError,
			want: `{
	"error": "Unknown error"
}`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/events/edit/1", bytes.NewBufferString(tt.body))
			w := httptest.NewRecorder()

			logg := newLogger(t)
			app := mocks.NewApplication(t)
			app.On("EditEvent", mock.Anything, "1", eventArg).Return(tt.returns...)
			server := NewServer(logg, app, "", "")

			server.editEventHandler(w, req)

			require.Equal(t, tt.status, w.Code)
			require.Equal(t, tt.want, w.Body.String())
		})
	}
}

func TestEditEventHandlerValidation(t *testing.T) {
	body := `{
	"id": "66be96d3-3d5d-4aec",
	"title": "testtitletoolongtesttitletoolong",
	"user_id": "66be96d3-3d5d-4aec",
	"date": "2024-09-23T00:00:00Z",
	"end_date": "2024-09-20T00:00:00Z"
}`
	req := httptest.NewRequest(http.MethodGet, "/events/edit/1", bytes.NewBufferString(body))
	w := httptest.NewRecorder()

	logg := newLogger(t)
	app := mocks.NewApplication(t)
	server := NewServer(logg, app, "", "")

	server.editEventHandler(w, req)

	require.Equal(t, http.StatusPartialContent, w.Code)
	require.Equal(t, `{
	"error": {
		"end_date": "too early",
		"id": "not valid uuid",
		"title": "too long",
		"user_id": "not valid uuid"
	}
}`, w.Body.String())
}

var testsEventList = []struct {
	name    string
	returns []interface{}
	want    string
	status  int
}{
	{
		name:    "success",
		returns: []interface{}{[]storage.Event{{ID: "1", Title: "test"}}, nil},
		status:  http.StatusOK,
		want: `{
	"events": [
		{
			"id": "1",
			"title": "test",
			"date": "0001-01-01T00:00:00Z",
			"end_date": "0001-01-01T00:00:00Z",
			"description": "",
			"user_id": "",
			"advance_notification_period": 0
		}
	]
}`,
	},
	{
		name:    "not found",
		returns: []interface{}{nil, storage.ErrNoEventsFound},
		status:  http.StatusNotFound,
		want: `{
	"error": "No events found"
}`,
	},
	{
		name:    "internal error",
		returns: []interface{}{nil, errors.New("internal error")},
		status:  http.StatusInternalServerError,
		want: `{
	"error": "Unknown error"
}`,
	},
}

func getDate(t *testing.T) time.Time {
	t.Helper()
	date, err := time.Parse("2006-01-02", "2024-09-23")
	if err != nil {
		t.Fatal(err)
	}

	return date
}

func TestGetEventsDayHandler(t *testing.T) {
	for _, tt := range testsEventList {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/events/day/2024-09-23", nil)
			w := httptest.NewRecorder()

			logg := newLogger(t)
			app := mocks.NewApplication(t)
			app.On("GetEventsListDay", mock.Anything, getDate(t)).Return(tt.returns...)
			server := NewServer(logg, app, "", "")

			server.getEventsDayHandler(w, req)
			require.Equal(t, tt.status, w.Code)
			require.Equal(t, tt.want, w.Body.String())
		})
	}
}

func TestGetEventsWeekHandler(t *testing.T) {
	for _, tt := range testsEventList {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/events/week/2024-09-23", nil)
			w := httptest.NewRecorder()

			logg := newLogger(t)
			app := mocks.NewApplication(t)
			app.On("GetEventsListWeek", mock.Anything, getDate(t)).Return(tt.returns...)
			server := NewServer(logg, app, "", "")

			server.getEventsWeekHandler(w, req)
			require.Equal(t, tt.status, w.Code)
			require.Equal(t, tt.want, w.Body.String())
		})
	}
}

func TestGetEventsMonthHandler(t *testing.T) {
	for _, tt := range testsEventList {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/events/week/2024-09-23", nil)
			w := httptest.NewRecorder()

			logg := newLogger(t)
			app := mocks.NewApplication(t)
			app.On("GetEventsListMonth", mock.Anything, getDate(t)).Return(tt.returns...)
			server := NewServer(logg, app, "", "")

			server.getEventsMonthHandler(w, req)
			require.Equal(t, tt.status, w.Code)
			require.Equal(t, tt.want, w.Body.String())
		})
	}
}
