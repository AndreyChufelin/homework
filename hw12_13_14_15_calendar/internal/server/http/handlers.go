package internalhttp

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/AndreyChufelin/homework/hw12_13_14_15_calendar/internal/storage"
	"github.com/AndreyChufelin/homework/hw12_13_14_15_calendar/internal/validator"
)

func (s *Server) createEventHandler(w http.ResponseWriter, r *http.Request) {
	logg := s.logger.With("handler", "createEventHandler")
	var event storage.Event
	r.Body = http.MaxBytesReader(w, r.Body, 1048576)
	err := json.NewDecoder(r.Body).Decode(&event)
	if err != nil {
		logg.Error("event failed to decode json", "error", err)
		s.errorResponse(w, http.StatusBadRequest, "Bad request")
		return
	}

	validator := validator.New()
	storage.ValidateEvent(*validator, event)
	if !validator.Valid() {
		logg.Warn("event validation failed", "error", validator.Errors)
		s.errorResponse(w, http.StatusPartialContent, validator.Errors)
		return
	}

	err = s.app.CreateEvent(r.Context(), event)
	if err != nil {
		if errors.Is(err, storage.ErrEventAlreadyExists) {
			logg.Warn("event already exist")
			s.errorResponse(w, http.StatusConflict, "Event already exist")
			return
		}
		logg.Error("failed create event", "error", err)
		s.errorResponse(w, http.StatusInternalServerError, "Unknown error")
		return
	}

	s.writeJSON(w, http.StatusOK, wrapper{"message": "Success"})
}

func (s *Server) getEventHandler(w http.ResponseWriter, r *http.Request) {
	logg := s.logger.With("handler", "getEventHandler")
	id := r.PathValue("id")

	event, err := s.app.GetEvent(r.Context(), id)
	if err != nil {
		if errors.Is(err, storage.ErrEventDoesntExist) {
			logg.Warn("event not found")
			s.errorResponse(w, http.StatusNotFound, "Event not found")
			return
		}
		logg.Error("failed get event", "error", err)
		s.errorResponse(w, http.StatusInternalServerError, "Unknown error")
		return
	}

	s.writeJSON(w, http.StatusOK, wrapper{"event": event})
}

func (s *Server) deleteEventHandler(w http.ResponseWriter, r *http.Request) {
	logg := s.logger.With("handler", "deleteEventHandler")
	id := r.PathValue("id")

	err := s.app.DeleteEvent(r.Context(), id)
	if err != nil {
		if errors.Is(err, storage.ErrEventDoesntExist) {
			logg.Warn("Delete event not found")
			s.errorResponse(w, http.StatusNotFound, "Event not found")
			return
		}
		logg.Error("Failed delete event", "error", err)
		s.errorResponse(w, http.StatusInternalServerError, "Unknown error")
		return
	}

	s.writeJSON(w, http.StatusOK, wrapper{"message": "Success"})
}

func (s *Server) editEventHandler(w http.ResponseWriter, r *http.Request) {
	logg := s.logger.With("handler", "editEventHandler")
	var event storage.Event
	r.Body = http.MaxBytesReader(w, r.Body, 1048576)
	err := json.NewDecoder(r.Body).Decode(&event)
	if err != nil {
		logg.Error("failed to decode json", "error", err)
		s.errorResponse(w, http.StatusBadRequest, "Bad request")
		return
	}
	id := r.PathValue("id")

	validator := validator.New()
	storage.ValidateEvent(*validator, event)
	if !validator.Valid() {
		logg.Warn("event validation failed", "error", validator.Errors)
		s.errorResponse(w, http.StatusPartialContent, validator.Errors)
		return
	}

	err = s.app.EditEvent(r.Context(), id, event)
	if err != nil {
		if errors.Is(err, storage.ErrEventDoesntExist) {
			logg.Warn("event not found")
			s.errorResponse(w, http.StatusNotFound, "Event not found")
			return
		}
		logg.Error("failed edit event", "error", err)
		s.errorResponse(w, http.StatusInternalServerError, "Unknown error")
		return
	}

	s.writeJSON(w, http.StatusOK, wrapper{"message": "Success"})
}

func getDateParam(r *http.Request) (time.Time, error) {
	param := r.PathValue("date")
	fmt.Println("param", param, r.URL.Path)

	date, err := time.Parse("2006-01-02", param)
	if err != nil {
		return time.Time{}, fmt.Errorf("wrong date parameter")
	}
	fmt.Println("date", date, r.URL.Path)

	return date, nil
}

func (s *Server) getEventsDayHandler(w http.ResponseWriter, r *http.Request) {
	logg := s.logger.With("handler", "getEventsDayEventHandler")
	date, err := getDateParam(r)
	if err != nil {
		logg.Warn("wrong date parameter")
		s.errorResponse(w, http.StatusBadRequest, "Wrong date parameter")
		return
	}

	events, err := s.app.GetEventsListDay(r.Context(), date)
	if err != nil {
		if errors.Is(err, storage.ErrNoEventsFound) {
			logg.Warn("no events found for day %s", date.Format("2006-01-02"))
			s.errorResponse(w, http.StatusNotFound, "No events found")
			return
		}
		logg.Error("failed get day events")
		s.errorResponse(w, http.StatusInternalServerError, "Unknown error")
		return
	}

	s.writeJSON(w, http.StatusOK, wrapper{"events": events})
}

func (s *Server) getEventsWeekHandler(w http.ResponseWriter, r *http.Request) {
	logg := s.logger.With("handler", "getEventsWeekEventHandler")
	date, err := getDateParam(r)
	if err != nil {
		logg.Warn("wrong date parameter")
		s.errorResponse(w, http.StatusBadRequest, "Wrong date parameter")
		return
	}

	events, err := s.app.GetEventsListWeek(r.Context(), date)
	if err != nil {
		if errors.Is(err, storage.ErrNoEventsFound) {
			logg.Warn("no events found for day %s", date.Format("2006-01-02"))
			s.errorResponse(w, http.StatusNotFound, "No events found")
			return
		}
		logg.Error("failed get day events")
		s.errorResponse(w, http.StatusInternalServerError, "Unknown error")
		return
	}

	s.writeJSON(w, http.StatusOK, wrapper{"events": events})
}

func (s *Server) getEventsMonthHandler(w http.ResponseWriter, r *http.Request) {
	logg := s.logger.With("handler", "getEventsMonthEventHandler")
	date, err := getDateParam(r)
	if err != nil {
		logg.Warn("wrong date parameter")
		s.errorResponse(w, http.StatusBadRequest, "Wrong date parameter")
		return
	}

	events, err := s.app.GetEventsListMonth(r.Context(), date)
	if err != nil {
		if errors.Is(err, storage.ErrNoEventsFound) {
			logg.Warn("no events found for day %s", date.Format("2006-01-02"))
			s.errorResponse(w, http.StatusNotFound, "No events found")
			return
		}
		logg.Error("failed get day events")
		s.errorResponse(w, http.StatusInternalServerError, "Unknown error")
		return
	}

	s.writeJSON(w, http.StatusOK, wrapper{"events": events})
}
