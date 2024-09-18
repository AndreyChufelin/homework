package memorystorage

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/AndreyChufelin/homework/hw12_13_14_15_calendar/internal/storage"
	"github.com/google/uuid"
)

type Storage struct {
	events map[string]storage.Event
	mu     sync.RWMutex
}

func New() *Storage {
	return &Storage{events: make(map[string]storage.Event)}
}

func (s *Storage) CreateEvent(_ context.Context, event storage.Event) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	id := event.ID
	if id == "" {
		id = uuid.New().String()
	} else {
		_, ok := s.events[id]
		if ok {
			return fmt.Errorf("creating event with id %s: %w", id, storage.ErrEventAlreadyExists)
		}
	}

	s.events[id] = event

	return nil
}

func (s *Storage) GetEvent(_ context.Context, id string) (*storage.Event, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	event, ok := s.events[id]
	if !ok {
		return nil, fmt.Errorf("getting event with id %s: %w", id, storage.ErrEventDoesntExist)
	}

	return &event, nil
}

func (s *Storage) DeleteEvent(_ context.Context, id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	_, ok := s.events[id]
	if !ok {
		return fmt.Errorf("deleting event with id: %w", storage.ErrEventDoesntExist)
	}

	delete(s.events, id)

	return nil
}

func (s *Storage) EditEvent(_ context.Context, id string, update storage.Event) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	_, ok := s.events[id]
	if !ok {
		return fmt.Errorf("edit event with id %s: %w", id, storage.ErrEventDoesntExist)
	}

	s.events[id] = update

	return nil
}

func (s *Storage) GetEventsListDay(_ context.Context, date time.Time) ([]storage.Event, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var result []storage.Event
	for _, event := range s.events {
		if event.Date.Year() == date.Year() && event.Date.YearDay() == date.YearDay() {
			result = append(result, event)
		}
	}

	return result, nil
}

func (s *Storage) GetEventsListWeek(_ context.Context, date time.Time) ([]storage.Event, error) {
	return s.getEventsListTo(date, date.AddDate(0, 0, 7)), nil
}

func (s *Storage) GetEventsListMonth(_ context.Context, date time.Time) ([]storage.Event, error) {
	return s.getEventsListTo(date, date.AddDate(0, 1, 0)), nil
}

func (s *Storage) getEventsListTo(start time.Time, end time.Time) []storage.Event {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var result []storage.Event
	for _, event := range s.events {
		if event.Date.After(start) && event.Date.Before(end) {
			result = append(result, event)
		}
	}

	return result
}
