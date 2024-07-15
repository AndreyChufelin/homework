package hw05parallelexecution

import (
	"errors"
	"fmt"
	"math/rand"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"go.uber.org/goleak"
)

func createTasks(tasksCount int, err error) ([]Task, *int32, time.Duration) {
	tasks := make([]Task, 0, tasksCount)
	var runTasksCount int32
	var sumTime time.Duration

	for i := 0; i < tasksCount; i++ {
		taskSleep := time.Millisecond * time.Duration(rand.Intn(100))
		sumTime += taskSleep

		tasks = append(tasks, func() error {
			time.Sleep(taskSleep)
			atomic.AddInt32(&runTasksCount, 1)

			return err
		})
	}

	return tasks, &runTasksCount, sumTime
}

func TestRun(t *testing.T) {
	defer goleak.VerifyNone(t)

	t.Run("if were errors in first M tasks, than finished not more N+M tasks", func(t *testing.T) {
		tasksCount := 50
		taskErr := fmt.Errorf("error from task")
		tasks, runTasksCount, _ := createTasks(tasksCount, taskErr)

		workersCount := 10
		maxErrorsCount := 23
		err := Run(tasks, workersCount, maxErrorsCount)

		require.Truef(t, errors.Is(err, ErrErrorsLimitExceeded), "actual err - %v", err)
		require.LessOrEqual(t, atomic.LoadInt32(runTasksCount), int32(workersCount+maxErrorsCount),
			"extra tasks were started")
	})

	t.Run("tasks without errors", func(t *testing.T) {
		tasksCount := 50

		tasks, runTasksCount, sumTime := createTasks(tasksCount, nil)

		workersCount := 5
		maxErrorsCount := 1

		start := time.Now()
		err := Run(tasks, workersCount, maxErrorsCount)
		elapsedTime := time.Since(start)
		require.NoError(t, err)

		require.Equal(t, atomic.LoadInt32(runTasksCount), int32(tasksCount), "not all tasks were completed")
		require.LessOrEqual(t, int64(elapsedTime), int64(sumTime/2), "tasks were run sequentially?")
	})

	t.Run("if m less or equal 0, than ErrErrorsLimitExceeded error", func(t *testing.T) {
		tasksCount := 10
		taskErr := fmt.Errorf("error from task")
		tasks, runTasksCount, _ := createTasks(tasksCount, taskErr)

		err := Run(tasks, 5, -1)

		require.Truef(t, errors.Is(err, ErrErrorsLimitExceeded), "actual err - %v", err)
		require.Equal(t, atomic.LoadInt32(runTasksCount), int32(0), "no tasks should be run")
	})
	t.Run("if n less or equal 0", func(t *testing.T) {
		tasksCount := 10

		tasks, runTasksCount, _ := createTasks(tasksCount, nil)

		err := Run(tasks, -1, 5)

		require.Truef(t, errors.Is(err, ErrInvalidGoroutineCount), "actual err - %v", err)
		require.Equal(t, atomic.LoadInt32(runTasksCount), int32(0), "no tasks should be run")
	})
}
