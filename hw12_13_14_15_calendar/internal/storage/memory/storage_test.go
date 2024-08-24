package memorystorage

import (
	"strconv"
	"sync"
	"testing"
	"time"

	"github.com/AndreyChufelin/homework/hw12_13_14_15_calendar/internal/storage"
	"github.com/stretchr/testify/require"
)

var eventValidation = []struct {
	Name  string
	Event storage.Event
	Err   error
}{
	{
		Name:  "returns error if EndDate too early",
		Event: storage.Event{ID: "1", Date: time.Now().Add(time.Hour * 24), EndDate: time.Now()},
		Err:   storage.ErrEndDateTooEarly,
	},
}

func TestStorage(t *testing.T) {
	s := New()
	require.Equal(t, &Storage{events: make(map[string]storage.Event)}, s)
}

func TestCreateEvent(t *testing.T) {
	t.Run("creates event in storage", func(t *testing.T) {
		s := New()
		event := storage.Event{ID: "1"}
		err := s.CreateEvent(event)

		require.NoError(t, err)
		require.Equal(t, map[string]storage.Event{"1": event}, s.events)
	})

	t.Run("returns error if event already exists", func(t *testing.T) {
		s := New()
		err := s.CreateEvent(storage.Event{ID: "1"})
		require.NoError(t, err)

		err = s.CreateEvent(storage.Event{ID: "1"})

		require.ErrorIs(t, err, storage.ErrEventAlreadyExists)
	})

	t.Run("concurency", func(t *testing.T) {
		s := New()

		const goroutines = 100
		var wg sync.WaitGroup
		wg.Add(goroutines)

		for i := range goroutines {
			go func() {
				defer wg.Done()
				err := s.CreateEvent(storage.Event{ID: strconv.Itoa(i)})
				require.NoError(t, err)
			}()
		}

		wg.Wait()
	})
}

func TestCreateEventValidation(t *testing.T) {
	for _, v := range eventValidation {
		t.Run(v.Name, func(t *testing.T) {
			s := New()
			err := s.CreateEvent(v.Event)
			require.ErrorIs(t, err, v.Err)
		})
	}
}

func TestGetEvent(t *testing.T) {
	t.Run("returns event by id", func(t *testing.T) {
		s := New()
		e := storage.Event{ID: "1"}
		err := s.CreateEvent(e)
		require.NoError(t, err)

		event, err := s.GetEvent("1")

		require.NoError(t, err)
		require.Equal(t, e, event)
	})

	t.Run("returns error if event doesn't exist", func(t *testing.T) {
		s := New()
		event, err := s.GetEvent("1")

		require.ErrorIs(t, err, storage.ErrEventDoesntExist)
		require.Equal(t, storage.Event{}, event)
	})

	t.Run("concurency", func(t *testing.T) {
		s := New()
		const goroutines = 100
		var wg sync.WaitGroup
		wg.Add(goroutines)

		for i := range goroutines / 2 {
			go func() {
				defer wg.Done()
				err := s.CreateEvent(storage.Event{ID: strconv.Itoa(i)})
				require.NoError(t, err)
			}()
		}

		for i := range goroutines / 2 {
			go func() {
				defer wg.Done()
				_, _ = s.GetEvent(strconv.Itoa(i))
			}()
		}

		wg.Wait()
	})
}

func TestDeleteEvent(t *testing.T) {
	t.Run("deletes event by id", func(t *testing.T) {
		s := New()
		err := s.CreateEvent(storage.Event{ID: "1"})
		require.NoError(t, err)

		err = s.DeleteEvent("1")
		require.NoError(t, err)
		require.Equal(t, map[string]storage.Event{}, s.events)
	})

	t.Run("returns error if event doesn't exist", func(t *testing.T) {
		s := New()
		err := s.DeleteEvent("1")
		require.ErrorIs(t, err, storage.ErrEventDoesntExist)
	})

	t.Run("concurency", func(t *testing.T) {
		s := New()

		const goroutines = 100
		var wg sync.WaitGroup
		wg.Add(goroutines)

		for i := range goroutines {
			s.CreateEvent(storage.Event{ID: strconv.Itoa(i)})
		}

		for i := range goroutines {
			go func() {
				defer wg.Done()
				err := s.DeleteEvent(strconv.Itoa(i))
				require.NoError(t, err)
			}()
		}

		wg.Wait()
	})
}

func TestEditEvent(t *testing.T) {
	t.Run("changes event data by id", func(t *testing.T) {
		s := New()
		s.CreateEvent(storage.Event{ID: "1"})

		err := s.EditEvent("1", storage.Event{Title: "Event #1"})

		require.NoError(t, err)
		require.Equal(t, map[string]storage.Event{"1": {ID: "1", Title: "Event #1"}}, s.events)
	})

	t.Run("returns error if event doesn't exist", func(t *testing.T) {
		s := New()
		err := s.EditEvent("1", storage.Event{Title: "Event #1"})

		require.ErrorIs(t, err, storage.ErrEventDoesntExist)
	})

	t.Run("returns error if EndDate is too early", func(t *testing.T) {
		s := New()
		err := s.CreateEvent(storage.Event{ID: "1"})
		require.NoError(t, err)

		err = s.EditEvent("1", storage.Event{Date: time.Now().Add(24 * time.Hour), EndDate: time.Now()})

		require.ErrorIs(t, err, storage.ErrEndDateTooEarly)
	})

	t.Run("concurency", func(t *testing.T) {
		s := New()

		const goroutines = 100
		var wg sync.WaitGroup
		wg.Add(goroutines)

		for i := range goroutines {
			s.CreateEvent(storage.Event{ID: strconv.Itoa(i)})
		}

		for i := range goroutines {
			go func() {
				defer wg.Done()
				err := s.EditEvent(strconv.Itoa(i), storage.Event{Title: "Event #1"})
				require.NoError(t, err)
			}()
		}

		wg.Wait()
	})
}

func TestEditEventValidation(t *testing.T) {
	for _, v := range eventValidation {
		t.Run(v.Name, func(t *testing.T) {
			s := New()
			err := s.CreateEvent(storage.Event{ID: "1"})
			require.NoError(t, err)

			err = s.EditEvent("1", v.Event)

			require.ErrorIs(t, err, v.Err)
		})
	}
}

func TestGetEventListDay(t *testing.T) {
	s := New()
	s.CreateEvent(storage.Event{ID: "1", Date: time.Now().Add(-time.Hour * 25), EndDate: time.Now()})
	e2 := storage.Event{ID: "2", Date: time.Now().Add(-time.Hour * 10), EndDate: time.Now()}
	s.CreateEvent(e2)
	s.CreateEvent(storage.Event{ID: "3", Date: time.Now().Add(time.Hour * 25), EndDate: time.Now().Add(time.Hour * 30)})

	list := s.GetEventsListDay(time.Now())
	require.Equal(t, []storage.Event{e2}, list)
}

func TestGetEventListWeek(t *testing.T) {
	s := New()
	s.CreateEvent(storage.Event{ID: "1", Date: time.Now().AddDate(0, 0, -7), EndDate: time.Now()})
	e2 := storage.Event{ID: "2", Date: time.Now().AddDate(0, 0, 3), EndDate: time.Now().AddDate(0, 0, 5)}
	s.CreateEvent(e2)
	s.CreateEvent(storage.Event{ID: "3", Date: time.Now().AddDate(0, 0, 10), EndDate: time.Now().AddDate(0, 0, 15)})

	list := s.GetEventsListWeek(time.Now())
	require.Equal(t, []storage.Event{e2}, list)
}

func TestGetEventListMonth(t *testing.T) {
	s := New()
	s.CreateEvent(storage.Event{ID: "1", Date: time.Now().AddDate(0, -2, 0), EndDate: time.Now()})
	e2 := storage.Event{ID: "2", Date: time.Now().AddDate(0, 0, 3), EndDate: time.Now().AddDate(0, 0, 5)}
	s.CreateEvent(e2)
	s.CreateEvent(storage.Event{ID: "3", Date: time.Now().AddDate(0, 0, 35), EndDate: time.Now().AddDate(0, 2, 0)})

	list := s.GetEventsListMonth(time.Now())
	require.Equal(t, []storage.Event{e2}, list)
}
