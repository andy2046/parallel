package parallel

type (
	// MapFn is the interface for Map function.
	MapFn = func(interface{}) (interface{}, error)

	// FilterFn is the interface for Filter function.
	FilterFn = func(interface{}) (bool, error)

	// ReduceFn is the interface for Reduce function.
	ReduceFn = func(accumulator, current interface{}) (interface{}, error)
)
