package parallel_test

import (
	"fmt"
	"strings"

	"github.com/andy2046/parallel"
)

var (
	specialChar = "!@#$%^&*()"
)

func ExampleFilter_map_reduce() {
	words := []string{
		"The ",
		"quick ",
		"brown",
		" fox",
		" jumps",
		"over",
		"the ",
		"lazy ",
		"dog",
		specialChar,
	}
	size := 200
	concurrency := 10
	datum := make([]string, 0, size*len(words))
	for i := 0; i < size; i++ {
		for _, w := range words {
			datum = append(datum, w)
		}
	}

	filtered, filterErrs := parallel.FilterN(concurrency, genDatum(datum...), filterSpecialCharacters)
	go func() {
		for err := range filterErrs {
			if err != nil {
				fmt.Printf("Filter error: %v", err)
			}
		}
	}()

	mapped, mapErrs := parallel.MapN(concurrency, filtered, trimTitleWords)
	go func() {
		for err := range mapErrs {
			if err != nil {
				fmt.Printf("Map error: %v", err)
			}
		}
	}()

	result, err := parallel.Reduce(mapped, count, make(map[string]int))
	if err != nil {
		fmt.Printf("Reduce error: %v", err)
	}

	m, ok := result.(map[string]int)
	if !ok {
		fmt.Printf("Reduce error: %v", errTypeConv)
	}
	for k, v := range m {
		fmt.Printf("%v:%v\n", k, v)
	}
	// Unordered output:
	// The:400
	// Quick:200
	// Brown:200
	// Fox:200
	// Jumps:200
	// Over:200
	// Lazy:200
	// Dog:200
}

func genDatum(values ...string) <-chan interface{} {
	out := make(chan interface{})
	go func() {
		defer close(out)
		for _, v := range values {
			out <- v
		}
	}()
	return out
}

func trimTitleWords(s interface{}) (interface{}, error) {
	str, ok := s.(string)
	if !ok {
		return nil, errTypeConv
	}
	return strings.Title(strings.TrimSpace(str)), nil
}

func filterSpecialCharacters(s interface{}) (bool, error) {
	str, ok := s.(string)
	if !ok {
		return false, errTypeConv
	}
	return !strings.ContainsAny(str, specialChar), nil
}

func count(acc, cur interface{}) (interface{}, error) {
	m, ok := acc.(map[string]int)
	if !ok {
		return nil, errTypeConv
	}
	str, ok := cur.(string)
	if !ok {
		return nil, errTypeConv
	}
	if c, ok := m[str]; ok {
		m[str] = c + 1
	} else {
		m[str] = 1
	}
	return m, nil
}
