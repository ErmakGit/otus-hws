package hw05parallelexecution

import (
	"errors"
	"sync"
)

var ErrErrorsLimitExceeded = errors.New("errors limit exceeded")

type Task func() error

type ErrorCounter struct {
	mu    sync.RWMutex
	count int64
}

func worker(wg *sync.WaitGroup, tasks <-chan Task, errLimit int64, errCounter *ErrorCounter) {
	for task := range tasks {
		errCounter.mu.RLock()
		countErr := errCounter.count
		errCounter.mu.RUnlock()

		if errLimit > 0 && countErr > errLimit {
			break
		}

		err := task()
		if err != nil {
			errCounter.mu.Lock()
			errCounter.count++
			errCounter.mu.Unlock()
		}
	}
	wg.Done()
}

// Run starts tasks in n goroutines and stops its work when receiving m errors from tasks.
func Run(tasks []Task, n, m int) error {
	jobs := make(chan Task)
	var errCounter ErrorCounter

	wg := sync.WaitGroup{}
	for w := 1; w <= n; w++ {
		wg.Add(1)
		go worker(&wg, jobs, int64(m), &errCounter)
	}

	for _, task := range tasks {
		errCounter.mu.RLock()
		countErr := errCounter.count
		errCounter.mu.RUnlock()

		if m > 0 && int(countErr) > m {
			close(jobs)
			wg.Wait()

			return ErrErrorsLimitExceeded
		}

		jobs <- task
	}

	close(jobs)
	wg.Wait()

	return nil
}
