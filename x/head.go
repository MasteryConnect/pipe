package x

// Head only allows the first Limit number of messages through then stops the pipeline.
type Head struct {
	N int // limit message to N
}

// T implements the pipeline interface.
func (h Head) T(inMsgs <-chan interface{}, outMsgs chan<- interface{}, errs chan<- error) {
	count := 0
	for msg := range inMsgs {
		count++

		outMsgs <- msg

		if count == h.N {
			return
		}
	}
}
