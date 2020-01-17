package x

// Tail only allows the last Limit number of messages through the pipeline.
type Tail struct {
	N int // limit message to N
}

// T implements the pipeline interface.
func (t Tail) T(inMsgs <-chan interface{}, outMsgs chan<- interface{}, errs chan<- error) {
	buf := []interface{}{}
	for msg := range inMsgs {
		if len(buf) == t.N {
			buf = append(buf[1:t.N], msg)
		} else {
			buf = append(buf, msg)
		}
	}

	for _, msg := range buf {
		outMsgs <- msg
	}
}
