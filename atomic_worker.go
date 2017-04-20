package snowflake

import (
	"runtime"
	"sync/atomic"
)

// AtomicWorker snowflake worker implemented by package sync/atomic
type AtomicWorker struct {
	workerID int64
	lastTime int64
	lastSeq  int64
}

// NewAtomicWorker create a new AtomicWorker object
func NewAtomicWorker(workerID int64) (Worker, error) {
	if err := checkWorkerID(workerID); err != nil {
		return nil, err
	}

	return &AtomicWorker{
		workerID: workerID,
		lastTime: -1,
		lastSeq:  -1,
	}, nil
}

// Next get next ID
func (aw *AtomicWorker) Next() (int64, error) {
	var lastTime, lastSeq int64
	var seq, now int64

	for {
		seq, now = 0, nowMillis()

		lastTime = atomic.LoadInt64(&aw.lastTime)
		if lastTime > now {
			continue
		}

		lastSeq = atomic.LoadInt64(&aw.lastSeq)
		if now == lastTime {
			seq = sequenceMask & (lastSeq + 1)
			if seq == 0 {
				// reach to max sequence, wait
				now = waitUntilNextMillis(now)
			}
		}

		if !atomic.CompareAndSwapInt64(&aw.lastTime, lastTime, now) ||
			!atomic.CompareAndSwapInt64(&aw.lastSeq, lastSeq, seq) {
			runtime.Gosched()
			continue
		}

		return combine(now, aw.workerID, seq), nil
	}
}
