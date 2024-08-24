package memorystorage

import (
	"fmt"
	"reflect"
	"sync"
	"time"

	storage "github.com/AndreyChufelin/homework/hw12_13_14_15_calendar/internal/storage"
)

type Storage struct {
	events map[string]storage.Event
	mu     sync.RWMutex //nolint:unused
}

func New() *Storage {
	return &Storage{events: make(map[string]storage.Event)}
}

func (s *Storage) CreateEvent(event storage.Event) error {
	s.mu.RLock()
	_, ok := s.events[event.ID]
	s.mu.RUnlock()
	if ok {
		return fmt.Errorf("memorystorage.CreateEvent: %w", storage.ErrEventAlreadyExists)
	}

	err := validateEvent(event)
	if err != nil {
		return fmt.Errorf("memorystorage.CreateEvent: %w", err)
	}

	s.mu.Lock()
	s.events[event.ID] = event
	s.mu.Unlock()

	return nil
}

func (s *Storage) GetEvent(id string) (storage.Event, error) {
	s.mu.RLock()
	event, ok := s.events[id]
	s.mu.RUnlock()
	if !ok {
		return storage.Event{}, fmt.Errorf("memorystorage.GetEvent: %w", storage.ErrEventDoesntExist)
	}

	return event, nil
}

func (s *Storage) DeleteEvent(id string) error {
	s.mu.RLock()
	_, ok := s.events[id]
	s.mu.RUnlock()
	if !ok {
		return fmt.Errorf("memorystorage.DeleteEvent: %w", storage.ErrEventDoesntExist)
	}

	s.mu.Lock()
	delete(s.events, id)
	s.mu.Unlock()

	return nil
}

func (s *Storage) EditEvent(id string, update storage.Event) error {
	s.mu.RLock()
	event, ok := s.events[id]
	s.mu.RUnlock()
	if !ok {
		return fmt.Errorf("memorystorage.DeleteEvent: %w", storage.ErrEventDoesntExist)
	}

	s.mu.Lock()

	origVal := reflect.ValueOf(&event).Elem()
	updateVal := reflect.ValueOf(update)

	for i := 0; i < updateVal.NumField(); i++ {
		origField := origVal.Field(i)
		updateField := updateVal.Field(i)

		if !updateField.IsZero() {
			origField.Set(updateField)
		}
	}

	err := validateEvent(event)
	if err != nil {
		return fmt.Errorf("memorystorage.EditEvent: %w", err)
	}

	s.events[id] = event

	s.mu.Unlock()

	return nil
}

func (s *Storage) GetEventsListDay(date time.Time) []storage.Event {
	var result []storage.Event
	s.mu.RLock()
	for _, event := range s.events {
		if event.Date.Year() == date.Year() && event.Date.YearDay() == date.YearDay() {
			result = append(result, event)
		}
	}
	s.mu.RUnlock()

	return result
}

func (s *Storage) GetEventsListWeek(date time.Time) []storage.Event {
	return s.getEventsListTo(date, date.AddDate(0, 0, 7))
}

func (s *Storage) GetEventsListMonth(date time.Time) []storage.Event {
	return s.getEventsListTo(date, date.AddDate(0, 1, 0))
}

func (s *Storage) getEventsListTo(start time.Time, end time.Time) []storage.Event {
	var result []storage.Event
	s.mu.RLock()
	for _, event := range s.events {
		if event.Date.After(start) && event.Date.Before(end) {
			result = append(result, event)
		}
	}
	s.mu.RUnlock()

	return result
}

func validateEvent(event storage.Event) error {
	if event.Date.After(event.EndDate) {
		return fmt.Errorf("memorystorage.CreateEvent: %w", storage.ErrEndDateTooEarly)
	}

	return nil
}
