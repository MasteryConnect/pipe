package x_test

import (
	"time"

	l "github.com/MasteryConnect/pipe/line"
	"github.com/MasteryConnect/pipe/x"
)

func ExampleRateLimit_T() {
	l.New().SetP(func(out chan<- interface{}, errs chan<- error) {
		out <- "1"
		out <- "2"
		out <- "3"
	}).Add(
		x.RateLimit{N: 1, Per: time.Millisecond}.T,
		l.Stdout,
	).Run()
	// Output:
	// 1
	// 2
	// 3
}
