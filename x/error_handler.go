package x

import (
	"sync"

	l "github.com/masteryconnect/pipe/line"
)

type ErrorHandler struct {
	TaskToTry    l.InlineTfunc
	ErrorHandler l.Tfunc
}

// Process a message through the 'try' function.
// If there is an error, then pass it to the error handler. The error
// handler can determine what should happen, e.g. log in a special way
// and pass the error on, or ignore the error
func (eh ErrorHandler) T(in <-chan interface{}, out chan<- interface{}, errs chan<- error) {
	var wg sync.WaitGroup
	var errIn chan interface{}

	// Setup the error handler channel and goroutine if present
	errIn = make(chan interface{})
	wg.Add(1)
	go func() {
		defer wg.Done()
		eh.ErrorHandler(errIn, out, errs)
	}()

	// For each message processed by the 'try' function
	for msg := range in {
		outMsg, err := eh.TaskToTry(msg)

		if err == nil { // No error, then pass the message on
			out <- outMsg
		} else {
			// Error, so pass it on to the error handler if it is present
			// TODO: what should we do with err?
			if outMsg == nil {
				errIn <- msg
			} else {
				errIn <- outMsg
			}
		}
	}

	close(errIn)
	wg.Wait()
}
