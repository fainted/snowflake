Snowflake
=====

Unique ID generator based on Twitter's snowflake.
IDs are 64bits positive integers(highest bit zero), format in detail:
```
+------------------------+--------------------------+-----------------+
| 43bits timestamp in ms | 8bits worker(machine) ID | 12bits sequence |
+------------------------+--------------------------+-----------------+
```

Installation
------

```
go get github.com/fainted/snowflake
```

Usage
------

```go
import "github.com/fainted/snowflake"

// worker ID: 1
worker, err := snowflake.NewChannelWorker(1)
if err != nil {
    return
}

id, err := worker.Next()
```

You can use snowflake.AtomicWorker as well, which is very fast and costs more CPU.

Benchmark
------

```
$ go test -v -bench . -run "Gen" -benchmem
BenchmarkAtomicGen-2             5000000           244 ns/op           0 B/op          0 allocs/op
BenchmarkAtomicParallel-2        5000000           244 ns/op           0 B/op          0 allocs/op
BenchmarkChannelGen-2            3000000           475 ns/op           0 B/op          0 allocs/op
BenchmarkChannelParallel-2       3000000           496 ns/op           0 B/op          0 allocs/op
```
