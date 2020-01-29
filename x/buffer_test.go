package x_test

import (
	"fmt"

	l "github.com/masteryconnect/pipe/line"
	"github.com/masteryconnect/pipe/x"
)

func ExampleBuffer() {
	gateComplete := make(chan bool)
	gate := 0
	flush := make(chan bool)
	flushed := false
	msgsToSend := 100
	after := x.Count{Silent: true}

	l.New().SetP(func(out chan<- interface{}, errs chan<- error) {
		for i := 0; i < msgsToSend; i++ {
			out <- "foo"
		}

		// now wait for the buffer to load before checking the counts
		<-gateComplete
		fmt.Printf("before %d after %d\n", msgsToSend, after.Val())
		flush <- true
	}).Add(
		func(in <-chan interface{}, out chan<- interface{}, errs chan<- error) {
			for m := range in {
				gate++
				if gate == msgsToSend {
					gateComplete <- true
				}
				out <- m
			}
		},

		// Buffer will take in 100 messages regardless of downstream
		x.Buffer{N: msgsToSend}.T,

		l.Inline(func(m interface{}) (interface{}, error) {
			// don't process anything until we are flushed
			if !flushed {
				<-flush
				flushed = true
			}
			return m, nil // passthrough
		}),
		after.Use,
	).Run()

	fmt.Printf("before %d after %d\n", msgsToSend, after.Val())
	// Output:
	// before 100 after 0
	// before 100 after 100
}
