package hw05parallelexecution

import (
	"errors"
	"sync/atomic"
)

var ErrErrorsLimitExceeded = errors.New("errors limit exceeded")

type Task func() error

func producer(tasks []Task, jobs chan<- Task, stop <-chan struct{}) {
	for _, task := range tasks {
		select {
		case <-stop:
			return
		default:
		}

		jobs <- task
	}
}

func consumer(jobs <-chan Task, results chan<- bool, stop <-chan struct{}) {
	for {
		select {
		case <-stop:
			return
		case job := <-jobs:
			err := job()
			executed := true
			if err != nil {
				executed = false
			}

			results <- executed
		}
	}
}

// Run starts tasks in n goroutines and stops its work when receiving m errors from tasks.
func Run(tasks []Task, n, m int) error {
	jobs := make(chan Task)
	results := make(chan bool)

	stopsCount := n + 1
	stop := make(chan struct{}, stopsCount)

	go producer(tasks, jobs, stop)

	for c := 1; c <= n; c++ {
		go consumer(jobs, results, stop)
	}

	var errorsCounter int64
	var executeCount int64
	for result := range results {
		atomic.AddInt64(&executeCount, 1)

		if !result {
			atomic.AddInt64(&errorsCounter, 1)
		}

		if m > 0 && int(errorsCounter) == m {
			for w := 1; w <= stopsCount; w++ {
				stop <- struct{}{}
			}

			return ErrErrorsLimitExceeded
		}

		if int(executeCount) == len(tasks) {
			return nil
		}
	}

	return nil
}
