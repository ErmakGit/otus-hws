package hw05parallelexecution

import (
	"errors"
)

var ErrErrorsLimitExceeded = errors.New("errors limit exceeded")

type Task func() error

func worker(jobs <-chan Task, results chan<- bool, stop <-chan struct{}) {
	select {
	case <-stop:
		return
	default:
	}

	for job := range jobs {
		// fmt.Printf("Goroutine %v getting task \n", id)
		execute(job, results)
	}
}

func execute(task Task, results chan<- bool) {
	err := task()
	executed := true
	if err != nil {
		executed = false
	}

	results <- executed
}

// Run starts tasks in n goroutines and stops its work when receiving m errors from tasks.
func Run(tasks []Task, n, m int) error {
	// Place your code here.

	jobs := make(chan Task, len(tasks))
	results := make(chan bool, len(tasks))
	stop := make(chan struct{}, 1)

	// wg := sync.WaitGroup{}

	for w := 1; w <= n; w++ {
		// wg.Add(1)
		// go func(w int) {
		go worker(jobs, results, stop)
		// wg.Done()
		// }(w)
	}
	// wg.Wait()

	for _, job := range tasks {
		jobs <- job
	}
	close(jobs)

	errorsCounter := 0
	for res := 1; res <= len(tasks); res++ {
		if m > 0 && errorsCounter == m {
			stop <- struct{}{}

			return ErrErrorsLimitExceeded
		}

		result := <-results
		if !result {
			errorsCounter++
		}
	}

	return nil
}
