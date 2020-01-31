package x

/*
Batch will take N number of messages and combine them into one.
The combined records will be placed in a slice in the metadata
of the single message passed downstream. The metadata key holding the batch
is 'batch' and the type is a slice of messages ([]Message).
The body is left empty.
*/

import (
	"sync"
	"time"

	"github.com/MasteryConnect/pipe/message"
)

// Batch will take N number of messages and create a batch message
// that has the slice of the messages as the metdata key "batch".
// It will also combine the bodies into the body of the batch
// separated by newlines if the CombineBody is true.
type Batch struct {
	N    int
	Timeout time.Duration
	ByteLimit int

	closeCh chan bool
}

// CloseableBatch creates a new batch
func CloseableBatch(size int, timeout time.Duration, byteLimit int) Batch {
	return Batch{N: size, Timeout: timeout, ByteLimit: byteLimit, closeCh: make(chan bool)}
}


// T inplements the pipeline transform interface.
func (b Batch) T(in <-chan interface{}, out chan<- interface{}, errs chan<- error) {
	count := 0
	byteCount := 0
	batch := []interface{}{}
	timer := time.NewTimer(b.Timeout)
	var mx sync.Mutex

	if b.Timeout == 0 {
		timer.Stop()
	}

	sendBatch := func() {
		mx.Lock()
		timer.Stop()
		if len(batch) > 0 {
			bmsg := message.Batch(batch)
			out <- bmsg // passthrough
			count = 0
			batch = []interface{}{}
		}

		if b.Timeout > 0 {
			timer.Reset(b.Timeout)
		}
		mx.Unlock()
	}

	processMsg := func(msg interface{}) {
		if b.ByteLimit > 0 { // if we are using ByteLimit
			str := message.String(msg)
			if byteCount + len(str) > b.ByteLimit {
				sendBatch() // send without adding new msg
				byteCount = len(str)
			} else {
				byteCount += len(str)
			}
		}
		count++
		batch = append(batch, msg)
		if count == b.N {
			sendBatch()
		}
	}

	closed := false
	for {
		select {
		case <-timer.C:
			sendBatch()
		case msg, ok := <-in:
			if !ok {
				closed = true
				break
			}
			processMsg(msg)
		case <-b.closeCh:
			closed = true
			break
		}

		if closed {
			break
		}
	}

	if len(batch) > 0 {
		sendBatch()
	}
}

// Close will send a signal to the batcher to break out and shutdown.
// This will send one final batch of whatever is left.
func (b *Batch) Close() {
	b.closeCh <- true
}
