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

func(er *ErrorCounter) Inc() {
	er.mu.Lock()
	er.count++
	er.mu.Unlock()
}

func(er *ErrorCounter) Get() int64 {
	er.mu.RLock()
	countErr := er.count
	er.mu.RUnlock()

	return countErr
}

func worker(wg *sync.WaitGroup, tasks <-chan Task, errLimit int64, errCounter *ErrorCounter) {
	for task := range tasks {
		if errLimit > 0 && errCounter.Get() > errLimit {
			break
		}

		err := task()
		if err != nil {
			errCounter.Inc()
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
		if m > 0 && int(errCounter.Get()) > m {
			break
		}

		jobs <- task
	}

	close(jobs)
	wg.Wait()

	if m > 0 && int(errCounter.Get()) > m {
		return ErrErrorsLimitExceeded
	}

	return nil
}
