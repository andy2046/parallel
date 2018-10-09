package parallel

import (
	"sync"

	"github.com/pkg/errors"
)

// Reduce consumes from `input` channel and passes to the ReduceFn `f`,
// if `initial` is not provided, it will block and wait for the first value from `input` channel,
// and use it as `initial`.
func Reduce(input <-chan interface{}, f ReduceFn, initial ...interface{}) (interface{}, error) {
	var init interface{}
	if len(initial) > 0 {
		init = initial[0]
	} else {
		for x := range input {
			init = x
			break
		}
	}

	acc := init

	for x := range input {
		v, err := f(acc, x)
		if err != nil {
			return acc, err
		}
		acc = v
	}

	return acc, nil
}

// ReduceN starts `concurrency` of goroutines, each goroutine consumes from `input` channel
// and passes to the ReduceFn `f`, if `initial` is not provided, `initial` will be `nil`,
// if `initial` is provided, all goroutines use the same `initial`.
func ReduceN(concurrency int, input <-chan interface{}, f ReduceFn, initial ...interface{}) (interface{}, error) {
	output, errCh := make(chan interface{}, concurrency), make(chan error, concurrency)
	wg, done := new(sync.WaitGroup), make(chan struct{})
	closeSignal := make(chan struct{}, concurrency)

	go func() {
		select {
		case <-closeSignal:
			close(done)
		}
	}()
	wg.Add(concurrency)

	var init interface{}
	if len(initial) > 0 {
		init = initial[0]
	}

	for i := 0; i < concurrency; i++ {
		go func() {
			defer wg.Done()
			acc := init
			for {
				select {
				case _, open := <-done:
					if !open {
						return
					}
				default:
				}
				x, open := <-input
				if !open {
					if acc != nil {
						output <- acc
					}
					return
				}
				r, err := f(acc, x)
				if err != nil {
					closeSignal <- struct{}{}
					errCh <- err
					return
				}
				acc = r
			}
		}()
	}

	wg.Wait()
	close(closeSignal)
	close(errCh)
	close(output)
	var err error
	for e := range errCh {
		err = errors.Wrap(err, e.Error())
	}
	if err != nil {
		return nil, err
	}

	return Reduce(output, f)
}
