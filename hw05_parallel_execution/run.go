package hw05parallelexecution

import (
	"errors"
	"sync"
	"sync/atomic"
)

var (
	ErrErrorsLimitExceeded   = errors.New("errors limit exceeded")
	ErrInvalidGoroutineCount = errors.New("invalid goroutine count")
)

type Task func() error

// Run starts tasks in n goroutines and stops its work when receiving m errors from tasks.
func Run(tasks []Task, n, m int) error {
	if m <= 0 {
		return ErrErrorsLimitExceeded
	}
	if n <= 0 {
		return ErrInvalidGoroutineCount
	}

	var (
		err        error
		errorCount atomic.Int64
		wg         sync.WaitGroup
	)

	tasksCh := make(chan Task)

	for range n {
		wg.Add(1)

		go func() {
			defer wg.Done()
			for task := range tasksCh {
				err := task()
				if err != nil {
					errorCount.Add(1)
				}
			}
		}()
	}

	for _, t := range tasks {
		if errorCount.Load() >= int64(m) {
			err = ErrErrorsLimitExceeded
			break
		}

		tasksCh <- t
	}

	close(tasksCh)
	wg.Wait()

	return err
}
