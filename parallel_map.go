package parallel

import "sync"

// MapN starts `concurrency` of goroutines, in each goroutine it creates one goroutine for each `MapFn` in `fns`,
// to compose a chain of goroutines, then it consumes from `input` channel and passes to the goroutines chain,
// you must read from the error channel returned or it will deadlock.
func MapN(concurrency int, input <-chan interface{}, fns ...MapFn) (<-chan interface{}, <-chan error) {
	output, errCh, wg := make(chan interface{}), make(chan error, concurrency), new(sync.WaitGroup)
	wg.Add(2 * concurrency)

	for i := 0; i < concurrency; i++ {
		go func() {
			defer wg.Done()
			out, errs := Map(input, fns...)
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

// Map creates one goroutine for each `MapFn` in `fns`, to compose a chain of goroutines,
// then consumes from `input` channel and passes to the goroutines chain,
// you must read from the error channel returned or it will deadlock.
func Map(input <-chan interface{}, fns ...MapFn) (<-chan interface{}, <-chan error) {
	errCh, wg := make(chan error, 1), new(sync.WaitGroup)
	wg.Add(len(fns))
	go func() {
		wg.Wait()
		close(errCh)
	}()
	return mapWithErr(input, errCh, wg, fns...), errCh
}

func mapWithErr(input <-chan interface{}, errCh chan error, wg *sync.WaitGroup, fns ...MapFn) <-chan interface{} {
	if len(fns) == 0 {
		return input
	}
	input = mapWithErr(input, errCh, wg, fns[:len(fns)-1]...)
	f := fns[len(fns)-1]
	output := make(chan interface{})

	go func() {
		defer func() {
			close(output)
			wg.Done()
		}()

		for x := range input {
			v, err := f(x)
			if err != nil {
				errCh <- err
			}
			output <- v
		}
	}()
	return output
}
