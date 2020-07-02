package line_test

import (
	"fmt"
	"sync/atomic"

	l "github.com/MasteryConnect/pipe/line"
)

func ExampleMany() {
	spinupCnt := uint32(0)

	l.New().
		Add(
			l.Many(func(in <-chan interface{}, out chan<- interface{}, errs chan<- error) {
				// there are two of these running concurrently
				atomic.AddUint32(&spinupCnt, 1)
				for m := range in {
					out <- m // passthrough
				}
			}, 2),
		).
		Run()

	fmt.Println(spinupCnt)
	// Output: 2
}
