/*
Package snowflake provides unique ID generator based on twitter's snowflake.
	IDs' format
	+----------------+-----------------+------------------------+--------------+
	|   unused(1)    |  timestamp(43)  |      worker id(8)      | sequence(12) |
	+----------------+-----------------+------------------------+--------------+
	| ensure all IDs | in milliseconds | need assign in advance | max 4096/ms  |
	| are positive   | max: 2248-09-26 | max worker id 256      |              |
	+----------------+-----------------+------------------------+--------------+
*/
package snowflake

import (
	"errors"
	"time"
)

// errors define
var (
	ErrWrongWorkerID = errors.New("Invalid snowflake worker ID")
	ErrTimeGoesBack  = errors.New("System time goes back")
)

// format definition
const (
	TimestampBits = 43
	WorkerIDBits  = 8
	SequenceBits  = 12
)

const (
	timestampShift = WorkerIDBits + SequenceBits
	sequenceMask   = (1 << SequenceBits) - 1
)

// Worker snowflake ID generator interface
type Worker interface {
	// Next get next ID
	Next() (int64, error)
}

func nowMillis() int64 {
	return int64(time.Now().UnixNano() / 1000000)
}

func waitUntilNextMillis(old int64) int64 {
	now := nowMillis()
	for old == now {
		now = nowMillis()
	}

	return now
}

func checkWorkerID(workerID int64) error {
	if workerID < 0 || workerID >= (1<<WorkerIDBits) {
		return ErrWrongWorkerID
	}

	return nil
}

// combine get ID by bit operation
func combine(timestamp, workerID, sequence int64) int64 {
	return (timestamp << timestampShift) | (workerID << SequenceBits) | sequence
}

func init() {
	if TimestampBits+WorkerIDBits+SequenceBits != 64-1 {
		panic("TimestampBits+WorkerIDBits+SequenceBits != 64-1")
	}
}
