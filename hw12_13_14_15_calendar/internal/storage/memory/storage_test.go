package memorystorage

import (
	"context"
	"strconv"
	"sync"
	"testing"
	"time"

	"github.com/AndreyChufelin/homework/hw12_13_14_15_calendar/internal/storage"
	"github.com/stretchr/testify/require"
)

func TestStorage(t *testing.T) {
	s := New()
	require.Equal(t, &Storage{events: make(map[string]storage.Event)}, s)
}

func TestCreateEvent(t *testing.T) {
	t.Run("creates event in storage", func(t *testing.T) {
		s := New()
		event := storage.Event{ID: "1"}
		err := s.CreateEvent(context.TODO(), event)

		require.NoError(t, err)
		require.Equal(t, map[string]storage.Event{"1": {ID: "1"}}, s.events)
	})

	t.Run("concurency", func(t *testing.T) {
		s := New()

		const goroutines = 100
		var wg sync.WaitGroup
		wg.Add(goroutines)

		for i := range goroutines {
			go func() {
				defer wg.Done()
				err := s.CreateEvent(context.TODO(), storage.Event{ID: strconv.Itoa(i)})
				require.NoError(t, err)
			}()
		}

		wg.Wait()
	})
}

func TestGetEvent(t *testing.T) {
	t.Run("returns event by id", func(t *testing.T) {
		s := New()
		e := storage.Event{ID: "1"}
		err := s.CreateEvent(context.TODO(), e)
		require.NoError(t, err)

		event, err := s.GetEvent(context.TODO(), "1")

		require.NoError(t, err)
		require.Equal(t, &e, event)
	})

	t.Run("returns error if event doesn't exist", func(t *testing.T) {
		s := New()
		event, err := s.GetEvent(context.TODO(), "1")

		require.ErrorIs(t, err, storage.ErrEventDoesntExist)
		require.Nil(t, event)
	})

	t.Run("concurency", func(t *testing.T) {
		s := New()
		const goroutines = 100
		var wg sync.WaitGroup
		wg.Add(goroutines)

		for i := range goroutines / 2 {
			go func() {
				defer wg.Done()
				err := s.CreateEvent(context.TODO(), storage.Event{ID: strconv.Itoa(i)})
				require.NoError(t, err)
			}()
		}

		for i := range goroutines / 2 {
			go func() {
				defer wg.Done()
				_, _ = s.GetEvent(context.TODO(), strconv.Itoa(i))
			}()
		}

		wg.Wait()
	})
}

func TestDeleteEvent(t *testing.T) {
	t.Run("deletes event by id", func(t *testing.T) {
		s := New()
		err := s.CreateEvent(context.TODO(), storage.Event{ID: "1"})
		require.NoError(t, err)

		err = s.DeleteEvent(context.TODO(), "1")
		require.NoError(t, err)
		require.Equal(t, map[string]storage.Event{}, s.events)
	})

	t.Run("returns error if event doesn't exist", func(t *testing.T) {
		s := New()
		err := s.DeleteEvent(context.TODO(), "1")
		require.ErrorIs(t, err, storage.ErrEventDoesntExist)
	})

	t.Run("concurency", func(t *testing.T) {
		s := New()

		const goroutines = 100
		var wg sync.WaitGroup
		wg.Add(goroutines)

		for i := range goroutines {
			s.CreateEvent(context.TODO(), storage.Event{ID: strconv.Itoa(i)})
		}

		for i := range goroutines {
			go func() {
				defer wg.Done()
				err := s.DeleteEvent(context.TODO(), strconv.Itoa(i))
				require.NoError(t, err)
			}()
		}

		wg.Wait()
	})
}

func TestEditEvent(t *testing.T) {
	t.Run("changes event data by id", func(t *testing.T) {
		s := New()
		s.CreateEvent(context.TODO(), storage.Event{ID: "1"})

		err := s.EditEvent(context.TODO(), "1", storage.Event{ID: "1", Title: "Event #1"})

		require.NoError(t, err)
		require.Equal(t, map[string]storage.Event{"1": {ID: "1", Title: "Event #1"}}, s.events)
	})

	t.Run("returns error if event doesn't exist", func(t *testing.T) {
		s := New()
		err := s.EditEvent(context.TODO(), "1", storage.Event{Title: "Event #1"})

		require.ErrorIs(t, err, storage.ErrEventDoesntExist)
	})

	t.Run("concurency", func(t *testing.T) {
		s := New()

		const goroutines = 100
		var wg sync.WaitGroup
		wg.Add(goroutines)

		for i := range goroutines {
			s.CreateEvent(context.TODO(), storage.Event{ID: strconv.Itoa(i)})
		}

		for i := range goroutines {
			go func() {
				defer wg.Done()
				err := s.EditEvent(context.TODO(), strconv.Itoa(i), storage.Event{Title: "Event #1"})
				require.NoError(t, err)
			}()
		}

		wg.Wait()
	})
}

func TestGetEventListDay(t *testing.T) {
	s := New()
	s.CreateEvent(context.TODO(), storage.Event{ID: "1", Date: time.Now().Add(-time.Hour * 25), EndDate: time.Now()})
	e2 := storage.Event{ID: "2", Date: time.Now().Add(-time.Hour * 10), EndDate: time.Now()}
	s.CreateEvent(context.TODO(), e2)
	s.CreateEvent(context.TODO(),
		storage.Event{
			ID:      "3",
			Date:    time.Now().Add(time.Hour * 25),
			EndDate: time.Now().Add(time.Hour * 30),
		},
	)

	list, err := s.GetEventsListDay(context.TODO(), time.Now())
	require.NoError(t, err)
	require.Equal(t, []storage.Event{e2}, list)
}

func TestGetEventListWeek(t *testing.T) {
	s := New()
	s.CreateEvent(context.TODO(), storage.Event{ID: "1", Date: time.Now().AddDate(0, 0, -7), EndDate: time.Now()})
	e2 := storage.Event{ID: "2", Date: time.Now().AddDate(0, 0, 3), EndDate: time.Now().AddDate(0, 0, 5)}
	s.CreateEvent(context.TODO(), e2)
	s.CreateEvent(context.TODO(),
		storage.Event{
			ID:      "3",
			Date:    time.Now().AddDate(0, 0, 10),
			EndDate: time.Now().AddDate(0, 0, 15),
		},
	)

	list, err := s.GetEventsListWeek(context.TODO(), time.Now())
	require.NoError(t, err)
	require.Equal(t, []storage.Event{e2}, list)
}

func TestGetEventListMonth(t *testing.T) {
	s := New()
	s.CreateEvent(context.TODO(), storage.Event{ID: "1", Date: time.Now().AddDate(0, -2, 0), EndDate: time.Now()})
	e2 := storage.Event{ID: "2", Date: time.Now().AddDate(0, 0, 3), EndDate: time.Now().AddDate(0, 0, 5)}
	s.CreateEvent(context.TODO(), e2)
	s.CreateEvent(context.TODO(),
		storage.Event{
			ID:      "3",
			Date:    time.Now().AddDate(0, 0, 35),
			EndDate: time.Now().AddDate(0, 2, 0),
		},
	)

	list, err := s.GetEventsListMonth(context.TODO(), time.Now())
	require.NoError(t, err)
	require.Equal(t, []storage.Event{e2}, list)
}

func TestGetEventsToNotify(t *testing.T) {
	s := New()
	currentTime := time.Now()
	s.events["1"] = storage.Event{
		ID:                        "1",
		Date:                      currentTime,
		EndDate:                   currentTime,
		AdvanceNotificationPeriod: 3600,
		Notified:                  false,
	}
	s.events["2"] = storage.Event{
		ID:                        "2",
		Date:                      currentTime.Add(48 * time.Hour),
		EndDate:                   currentTime.Add(48 * time.Hour),
		AdvanceNotificationPeriod: 3600,
		Notified:                  false,
	}
	s.events["3"] = storage.Event{
		ID:                        "3",
		Date:                      currentTime,
		EndDate:                   currentTime,
		AdvanceNotificationPeriod: 3600,
		Notified:                  true,
	}

	events, err := s.GetEventsToNotify(context.TODO())

	require.NoError(t, err)

	require.Len(t, events, 1)
	require.Equal(t, events[0], s.events["1"])
}

func TestMarkNotified(t *testing.T) {
	s := New()
	s.events["1"] = storage.Event{ID: "1", Notified: false}
	s.events["2"] = storage.Event{ID: "2", Notified: false}
	s.events["3"] = storage.Event{ID: "3", Notified: false}

	err := s.MarkNotified(context.TODO(), []string{"1", "3"})

	require.NoError(t, err)
	require.Equal(t, s.events["1"].Notified, true)
	require.Equal(t, s.events["2"].Notified, false)
	require.Equal(t, s.events["3"].Notified, true)
}

func TestClearEvents(t *testing.T) {
	s := New()
	currentTime := time.Now()
	s.events["1"] = storage.Event{
		ID:      "1",
		Date:    currentTime,
		EndDate: currentTime,
	}
	s.events["2"] = storage.Event{
		ID:      "2",
		Date:    currentTime.Add(-2 * time.Hour),
		EndDate: currentTime.Add(-2 * time.Hour),
	}

	err := s.ClearEvents(context.TODO(), time.Hour)

	require.NoError(t, err)
	require.Len(t, s.events, 1)

	_, ok := s.events["1"]
	require.True(t, ok)
	_, ok = s.events["2"]
	require.False(t, ok)
}
