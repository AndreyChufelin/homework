package hw05parallelexecution

import (
	"errors"
	"sync"
	"sync/atomic"
)

var ErrErrorsLimitExceeded = errors.New("errors limit exceeded")

type Task func() error

// Run starts tasks in n goroutines and stops its work when receiving m errors from tasks.
func Run(tasks []Task, n, m int) error {
	tasksCh := make(chan Task)
	var errorCount atomic.Int64

	var wg sync.WaitGroup
	defer func() {
		close(tasksCh)
		wg.Wait()
	}()

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
			return ErrErrorsLimitExceeded
		}

		tasksCh <- t
	}

	return nil
}
