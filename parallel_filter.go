package parallel

import "sync"

// FilterN starts `concurrency` of goroutines, in each goroutine it creates one goroutine for each `FilterFn` in `fns`,
// to compose a chain of goroutines, then it consumes from `input` channel and passes to the goroutines chain,
// you must read from the error channel returned or it will deadlock.
func FilterN(concurrency int, input <-chan interface{}, fns ...FilterFn) (<-chan interface{}, <-chan error) {
	output, errCh, wg := make(chan interface{}), make(chan error, concurrency), new(sync.WaitGroup)
	wg.Add(2 * concurrency)

	for i := 0; i < concurrency; i++ {
		go func() {
			defer wg.Done()
			out, errs := Filter(input, fns...)
			go func() {
				defer wg.Done()
				for err := range errs {
					errCh <- err
				}
			}()
			for x := range out {
				output <- x
			}
		}()
	}

	go func() {
		wg.Wait()
		close(output)
		close(errCh)
	}()

	return output, errCh
}

// Filter creates one goroutine for each `FilterFn` in `fns`, to compose a chain of goroutines,
// then consumes from `input` channel and passes to the goroutines chain,
// you must read from the error channel returned or it will deadlock.
func Filter(input <-chan interface{}, fns ...FilterFn) (<-chan interface{}, <-chan error) {
	errCh, wg := make(chan error, 1), new(sync.WaitGroup)
	wg.Add(len(fns))
	go func() {
		wg.Wait()
		close(errCh)
	}()
	return filterWithErr(input, errCh, wg, fns...), errCh
}

func filterWithErr(input <-chan interface{}, errCh chan error, wg *sync.WaitGroup, fns ...FilterFn) <-chan interface{} {
	if len(fns) == 0 {
		return input
	}
	input = filterWithErr(input, errCh, wg, fns[:len(fns)-1]...)
	f := fns[len(fns)-1]
	output := make(chan interface{})

	go func() {
		defer func() {
			close(output)
			wg.Done()
		}()

		for x := range input {
			b, err := f(x)
			if err != nil {
				errCh <- err
			}
			if b {
				output <- x
			}
		}
	}()
	return output
}
