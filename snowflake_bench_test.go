package snowflake

import (
	"testing"
)

func atomicWorkerGen() (int64, error) {
	return atomicWorker.Next()
}

func channelWorkerGen() (int64, error) {
	return channelWorker.Next()
}

func BenchmarkAtomicGen(b *testing.B) {
	var err error
	for i := 0; i < b.N; i++ {
		_, err = atomicWorkerGen()
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkAtomicParallel(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_, err := atomicWorkerGen()
			if err != nil {
				b.Fatal(err)
			}
		}
	})
}

func BenchmarkChannelGen(b *testing.B) {
	var err error
	for i := 0; i < b.N; i++ {
		_, err = channelWorkerGen()
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkChannelParallel(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_, err := channelWorkerGen()
			if err != nil {
				b.Fatal(err)
			}
		}
	})
}
