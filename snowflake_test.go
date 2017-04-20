package snowflake

import (
	"log"
	"sort"
	"sync"
	"testing"
)

var (
	atomicWorker, channelWorker Worker
)

func TestWorkerID(t *testing.T) {
	for i := int64(0); i < 256; i++ {
		if err := checkWorkerID(i); err != nil {
			t.Fatalf("checkWorkerID(%d) failed: %v", i, err)
		}
	}

	invalidIDs := []int64{
		-5, -4, -3, -2, -1,
		256, 257, 258, 259, 260,
	}

	for _, i := range invalidIDs {
		if _, err := NewAtomicWorker(i); err == nil {
			t.Fatalf("NewAtomicWorker(%d) passed", i)
		}

		if _, err := NewChannelWorker(i); err == nil {
			t.Fatalf("NewChannelWorker(%d) passed", i)
		}
	}
}

type LongInts []int64

func (o LongInts) Len() int {
	return len(o)
}

func (o LongInts) Swap(i, j int) {
	o[i], o[j] = o[j], o[i]
}

func (o LongInts) Less(i, j int) bool {
	return o[i] < o[j]
}

func _TestWorker(t *testing.T, w Worker) {
	var mtx sync.Mutex
	var ids []int64

	var wg sync.WaitGroup

	for i := 0; i < 1000; i++ {
		wg.Add(1)
		go func(wg *sync.WaitGroup) {
			defer wg.Done()

			var chunk []int64
			for j := 0; j < 5000; j++ {
				id, err := w.Next()
				if err != nil {
					t.Error(err)
				}
				chunk = append(chunk, id)
			}

			mtx.Lock()
			defer mtx.Unlock()
			ids = append(ids, chunk...)
		}(&wg)
	}

	wg.Wait()

	sort.Sort(LongInts(ids))
	var prev = int64(-1)
	for i, v := range ids {
		if v <= prev {
			t.Errorf("Found duplicate %d at %d\n", v, i)
			return
		}
	}
}

func TestAtomicWorker(t *testing.T) {
	_TestWorker(t, atomicWorker)
}

func TestChannelWorker(t *testing.T) {
	_TestWorker(t, channelWorker)
}

func TestMain(m *testing.M) {
	var err error

	atomicWorker, err = NewAtomicWorker(1)
	if err != nil {
		log.Fatal("NewAtomicWorker: ", err)
	}

	channelWorker, err = NewChannelWorker(2)
	if err != nil {
		log.Fatal("NewChannelWorker: ", err)
	}

	m.Run()
}
