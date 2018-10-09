

# parallel
`import "github.com/andy2046/parallel"`

* [Overview](#pkg-overview)
* [Index](#pkg-index)

## <a name="pkg-overview">Overview</a>



## <a name="pkg-index">Index</a>
* [func Filter(input &lt;-chan interface{}, fns ...FilterFn) (&lt;-chan interface{}, &lt;-chan error)](#Filter)
* [func FilterN(concurrency int, input &lt;-chan interface{}, fns ...FilterFn) (&lt;-chan interface{}, &lt;-chan error)](#FilterN)
* [func Map(input &lt;-chan interface{}, fns ...MapFn) (&lt;-chan interface{}, &lt;-chan error)](#Map)
* [func MapN(concurrency int, input &lt;-chan interface{}, fns ...MapFn) (&lt;-chan interface{}, &lt;-chan error)](#MapN)
* [func Reduce(input &lt;-chan interface{}, f ReduceFn, initial ...interface{}) (interface{}, error)](#Reduce)
* [func ReduceN(concurrency int, input &lt;-chan interface{}, f ReduceFn, initial ...interface{}) (interface{}, error)](#ReduceN)
* [type FilterFn](#FilterFn)
* [type MapFn](#MapFn)
* [type ReduceFn](#ReduceFn)


#### <a name="pkg-files">Package files</a>
[parallel.go](./parallel.go) [parallel_filter.go](./parallel_filter.go) [parallel_map.go](./parallel_map.go) [parallel_reduce.go](./parallel_reduce.go) 





## <a name="Filter">func</a> [Filter](./parallel_filter.go?s=1150:1239#L40)
``` go
func Filter(input <-chan interface{}, fns ...FilterFn) (<-chan interface{}, <-chan error)
```
Filter creates one goroutine for each `FilterFn` in `fns`, to compose a chain of goroutines,
then consumes from `input` channel and passes to the goroutines chain,
you must read from the error channel returned or it will deadlock.



## <a name="FilterN">func</a> [FilterN](./parallel_filter.go?s=334:441#L8)
``` go
func FilterN(concurrency int, input <-chan interface{}, fns ...FilterFn) (<-chan interface{}, <-chan error)
```
FilterN starts `concurrency` of goroutines, in each goroutine it creates one goroutine for each `FilterFn` in `fns`,
to compose a chain of goroutines, then it consumes from `input` channel and passes to the goroutines chain,
you must read from the error channel returned or it will deadlock.



## <a name="Map">func</a> [Map](./parallel_map.go?s=1129:1212#L40)
``` go
func Map(input <-chan interface{}, fns ...MapFn) (<-chan interface{}, <-chan error)
```
Map creates one goroutine for each `MapFn` in `fns`, to compose a chain of goroutines,
then consumes from `input` channel and passes to the goroutines chain,
you must read from the error channel returned or it will deadlock.



## <a name="MapN">func</a> [MapN](./parallel_map.go?s=328:429#L8)
``` go
func MapN(concurrency int, input <-chan interface{}, fns ...MapFn) (<-chan interface{}, <-chan error)
```
MapN starts `concurrency` of goroutines, in each goroutine it creates one goroutine for each `MapFn` in `fns`,
to compose a chain of goroutines, then it consumes from `input` channel and passes to the goroutines chain,
you must read from the error channel returned or it will deadlock.



## <a name="Reduce">func</a> [Reduce](./parallel_reduce.go?s=262:356#L12)
``` go
func Reduce(input <-chan interface{}, f ReduceFn, initial ...interface{}) (interface{}, error)
```
Reduce consumes from `input` channel and passes to the ReduceFn `f`,
if `initial` is not provided, it will block and wait for the first value from `input` channel,
and use it as `initial`.



## <a name="ReduceN">func</a> [ReduceN](./parallel_reduce.go?s=873:985#L39)
``` go
func ReduceN(concurrency int, input <-chan interface{}, f ReduceFn, initial ...interface{}) (interface{}, error)
```
ReduceN starts `concurrency` of goroutines, each goroutine consumes from `input` channel
and passes to the ReduceFn `f`, if `initial` is not provided, `initial` will be `nil`,
if `initial` is provided, all goroutines use the same `initial`.




## <a name="FilterFn">type</a> [FilterFn](./parallel.go?s=171:213#L8)
``` go
type FilterFn = func(interface{}) (bool, error)
```
FilterFn is the interface for Filter function.










## <a name="MapFn">type</a> [MapFn](./parallel.go?s=71:117#L5)
``` go
type MapFn = func(interface{}) (interface{}, error)
```
MapFn is the interface for Map function.










## <a name="ReduceFn">type</a> [ReduceFn](./parallel.go?s=267:337#L11)
``` go
type ReduceFn = func(accumulator, current interface{}) (interface{}, error)
```
ReduceFn is the interface for Reduce function.




