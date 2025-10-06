package hw05parallelexecution

import (
	"errors"
	"fmt"
	"sync"
	"sync/atomic"
)

var (
	ErrErrorsLimitExceeded = errors.New("errors limit exceeded")
	ErrInvalidParameters   = errors.New("parameter are invalid")
)

type Task func() error

// Run starts tasks in n goroutines and stops its work when receiving m errors from tasks.
func Run(tasks []Task, n, m int) error {
	var (
		err        error
		taskErrors int32
	)

	if n <= 0 {
		return fmt.Errorf(" %w : the number of workers must be greater than zero", ErrInvalidParameters)
	}

	if len(tasks) == 0 {
		return fmt.Errorf(" %w : tnothing to do", ErrInvalidParameters)
	}

	ch := make(chan Task, n)
	wg := sync.WaitGroup{}
	wg.Add(n)

	for i := 0; i < n; i++ {
		go func() {
			for f := range ch {
				// если m < = 0, то ошибки не считаются, иначн, если достигнуто критическое количество ошибок - выходим
				if m > 0 && atomic.LoadInt32(&taskErrors) >= int32(m) {
					break
				}

				if er := f(); er != nil {
					atomic.AddInt32(&taskErrors, 1)
				}
			}
			wg.Done()
		}()
	}

	for _, v := range tasks {
		if m > 0 && atomic.LoadInt32(&taskErrors) >= int32(m) {
			err = ErrErrorsLimitExceeded
			break
		}
		ch <- v
	}

	close(ch)
	wg.Wait()

	return err
}
