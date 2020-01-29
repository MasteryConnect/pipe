package line_test

import (
	"fmt"

	l "github.com/masteryconnect/pipe/line"
)

func ExampleMany() {
	spinupCnt := 0

	l.New().
		Add(
			l.Many(func(in <-chan interface{}, out chan<- interface{}, errs chan<- error) {
				// there are two of these running concurrently
				spinupCnt++
				for m := range in {
					out <- m // passthrough
				}
			}, 2),
		).
		Run()

	fmt.Println(spinupCnt)
	// Output: 2
}
