package hw05parallelexecution

import (
	"errors"
	"sync"
	"sync/atomic"
)

var ErrErrorsLimitExceeded = errors.New("errors limit exceeded")

type Task func() error

type Counter struct {
	mu       sync.RWMutex
	errors   int64
	executed int64
}

func worker(tasks <-chan Task, counter *Counter) {
	for task := range tasks {
		err := task()
		if err != nil {
			counter.mu.Lock()
			atomic.AddInt64(&counter.errors, 1)
			counter.mu.Unlock()
		}

		counter.mu.Lock()
		atomic.AddInt64(&counter.executed, 1)
		counter.mu.Unlock()
	}
}

// Run starts tasks in n goroutines and stops its work when receiving m errors from tasks.
func Run(tasks []Task, n, m int) error {
	jobs := make(chan Task)
	var counter Counter

	for w := 1; w <= n; w++ {
		go worker(jobs, &counter)
	}

	for _, task := range tasks {
		counter.mu.RLock()
		errCount := counter.errors
		counter.mu.RUnlock()

		if m > 0 && int(errCount) > m {
			close(jobs)
			return ErrErrorsLimitExceeded
		}

		jobs <- task
	}

	close(jobs)

	for {
		counter.mu.RLock()
		executed := counter.executed
		counter.mu.RUnlock()

		if int(executed) == len(tasks) {
			break
		}
	}

	return nil
}
