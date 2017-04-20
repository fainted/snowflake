package snowflake

import ()

// ChannelWorker snowflake worker implemented using golang channel
type ChannelWorker struct {
	workerID int64

	chanIn  chan struct{}
	chanOut chan struct {
		id  int64
		err error
	}
}

// NewChannelWorker create a new ChannelWorker object
func NewChannelWorker(workerID int64) (Worker, error) {
	if err := checkWorkerID(workerID); err != nil {
		return nil, err
	}

	var cw = &ChannelWorker{
		workerID: workerID,
	}

	cw.chanIn = make(chan struct{}, 10)
	cw.chanOut = make(chan struct {
		id  int64
		err error
	}, 1)

	go cw.startBackground()

	return cw, nil
}

// Next get next ID
func (cw *ChannelWorker) Next() (int64, error) {
	cw.chanIn <- struct{}{}

	var out = struct {
		id  int64
		err error
	}{}

	out = <-cw.chanOut
	return out.id, out.err
}

func (cw *ChannelWorker) startBackground() {
	var lastSeq, lastTime = int64(-1), int64(-1)
	var seq, now int64

	for {
		_, ok := <-cw.chanIn
		if !ok {
			return
		}

		seq, now = int64(0), nowMillis()
		if lastTime > now {
			cw.chanOut <- struct {
				id  int64
				err error
			}{-1, ErrTimeGoesBack}
		} else if now == lastTime {
			seq = sequenceMask & (lastSeq + 1)
			if seq == 0 {
				// reach to max sequence, wait
				now = waitUntilNextMillis(now)
			}
		}

		cw.chanOut <- struct {
			id  int64
			err error
		}{
			id:  combine(now, cw.workerID, seq),
			err: nil,
		}

		lastSeq, lastTime = seq, now
	}
}
