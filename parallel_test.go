package parallel_test

import (
	"errors"
	"sync"
	"testing"

	p "github.com/andy2046/parallel"
)

var errTypeConv = errors.New("fail to convert interface to type")

func generator(values ...interface{}) <-chan interface{} {
	out := make(chan interface{})
	go func() {
		defer close(out)
		for _, v := range values {
			out <- v
		}
	}()
	return out
}

func square(n interface{}) (interface{}, error) {
	i, ok := n.(int)
	if !ok {
		return nil, errTypeConv
	}
	return i * i, nil
}

func increment(n interface{}) (interface{}, error) {
	i, ok := n.(int)
	if !ok {
		return nil, errTypeConv
	}
	return i + 1, nil
}

func greaterThanOne(n interface{}) (bool, error) {
	i, ok := n.(int)
	if !ok {
		return false, errTypeConv
	}
	return i > 1, nil
}

func divisibleByTwo(n interface{}) (bool, error) {
	i, ok := n.(int)
	if !ok {
		return false, errTypeConv
	}
	return i%2 == 0, nil
}

func sum(acc, cur interface{}) (interface{}, error) {
	if acc == nil {
		return cur, nil
	}
	i, ok := acc.(int)
	if !ok {
		return nil, errTypeConv
	}
	j, ok := cur.(int)
	if !ok {
		return nil, errTypeConv
	}
	return i + j, nil
}

func TestMap(t *testing.T) {
	in := generator(1, 2, 3)
	result := []int{2, 5, 10}
	wg := new(sync.WaitGroup)
	wg.Add(1)
	out, errCh := p.Map(in, square, increment)
	go func() {
		defer wg.Done()
		for err := range errCh {
			if err != nil {
				t.Error("TestMap:", err)
			}
		}
	}()

	i := 0
	for x := range out {
		if x.(int) != result[i] {
			t.Errorf("TestMap: want %v, got %v", result[i], x.(int))
		}
		i++
	}
	wg.Wait()
}

func TestMapN(t *testing.T) {
	in := generator(1, 2, 3)
	result := []int{2, 5, 10}
	wg := new(sync.WaitGroup)
	wg.Add(1)
	out, errCh := p.MapN(3, in, square, increment)
	go func() {
		defer wg.Done()
		for err := range errCh {
			if err != nil {
				t.Error("TestMapN:", err)
			}
		}
	}()

	for x := range out {
		if !contains(result, x.(int)) {
			t.Errorf("TestMapN: want one of %v, got %v", result, x.(int))
		}
	}
	wg.Wait()
}

func TestFilter(t *testing.T) {
	in := generator(1, 3, 5, 7, 8, 9)
	result := 8
	wg := new(sync.WaitGroup)
	wg.Add(1)
	out, errCh := p.Filter(in, greaterThanOne, divisibleByTwo)
	go func() {
		defer wg.Done()
		for err := range errCh {
			if err != nil {
				t.Error("TestFilter:", err)
			}
		}
	}()

	for x := range out {
		if x.(int) != result {
			t.Errorf("TestFilter: want %v, got %v", result, x.(int))
		}
	}
	wg.Wait()
}

func TestFilterN(t *testing.T) {
	in := generator(1, 3, 5, 6, 7, 8, 9)
	result := []int{6, 8}
	wg := new(sync.WaitGroup)
	wg.Add(1)
	out, errCh := p.FilterN(3, in, greaterThanOne, divisibleByTwo)
	go func() {
		defer wg.Done()
		for err := range errCh {
			if err != nil {
				t.Error("TestFilterN:", err)
			}
		}
	}()

	for x := range out {
		if !contains(result, x.(int)) {
			t.Errorf("TestFilterN: want one of %v, got %v", result, x.(int))
		}
	}
	wg.Wait()
}

func TestReduce(t *testing.T) {
	in := generator(1, 2, 3, 4)
	result := 10
	out, err := p.Reduce(in, sum)
	if err != nil {
		t.Error("TestReduce:", err)
	}
	if out.(int) != result {
		t.Errorf("TestReduce: want %v, got %v", result, out.(int))
	}

	{
		in := generator(1, 2, 3, 4)
		result := 15
		out, err := p.Reduce(in, sum, 5)
		if err != nil {
			t.Error("TestReduce:", err)
		}
		if out.(int) != result {
			t.Errorf("TestReduce: want %v, got %v", result, out.(int))
		}
	}

}

func TestReduceN(t *testing.T) {
	in := generator(1, 2, 3, 4)
	result := 10
	out, err := p.ReduceN(3, in, sum)
	if err != nil {
		t.Error("TestReduceN:", err)
	}
	if out.(int) != result {
		t.Errorf("TestReduceN: want %v, got %v", result, out.(int))
	}

	{
		in := generator(1, 2, 3, 4, 5)
		result := 15
		out, err := p.ReduceN(3, in, sum, 0)
		if err != nil {
			t.Error("TestReduceN:", err)
		}
		if out.(int) != result {
			t.Errorf("TestReduceN: want %v, got %v", result, out.(int))
		}
	}

}

func contains(s []int, i int) bool {
	for _, n := range s {
		if n == i {
			return true
		}
	}
	return false
}
