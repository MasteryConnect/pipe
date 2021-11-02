package line_test

import (
	"context"
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

func ExampleManyContext() {
	type ctxKey string
	var multiplyBy ctxKey = "multi"
	outsideCtx := context.WithValue(context.Background(), multiplyBy, uint32(10))

	spinupCnt := uint32(0)

	l.New().
		AddContext(
			l.ManyContext(func(ctx context.Context, in <-chan interface{}, out chan<- interface{}, errs chan<- error) {
				multiBy := ctx.Value(multiplyBy).(uint32)
				// there are two of these running concurrently
				atomic.AddUint32(&spinupCnt, 1*multiBy)
				for m := range in {
					out <- m // passthrough
				}
			}, 2),
		).
		RunContext(outsideCtx)

	fmt.Println(spinupCnt)
	// Output: 20
}
