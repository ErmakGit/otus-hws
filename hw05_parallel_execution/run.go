package hw05parallelexecution

import (
	"errors"
	"sync"
	"sync/atomic"
)

var ErrErrorsLimitExceeded = errors.New("errors limit exceeded")

type Task func() error

type ErrorsCounter struct {
	mu    sync.RWMutex
	count int64
}

func worker(wg *sync.WaitGroup, tasks <-chan Task, errCounter *ErrorsCounter) {
	for task := range tasks {
		err := task()
		if err != nil {
			errCounter.mu.Lock()
			atomic.AddInt64(&errCounter.count, 1)
			errCounter.mu.Unlock()
		}
	}
	wg.Done()
}

// Run starts tasks in n goroutines and stops its work when receiving m errors from tasks.
func Run(tasks []Task, n, m int) error {
	jobs := make(chan Task)
	var errCounter ErrorsCounter

	wg := sync.WaitGroup{}
	for w := 1; w <= n; w++ {
		wg.Add(1)
		go worker(&wg, jobs, &errCounter)
	}

	for _, task := range tasks {
		errCounter.mu.RLock()
		errCount := errCounter.count
		errCounter.mu.RUnlock()

		if m > 0 && int(errCount) > m {
			close(jobs)
			return ErrErrorsLimitExceeded
		}

		jobs <- task
	}

	close(jobs)
	wg.Wait()

	return nil
}
