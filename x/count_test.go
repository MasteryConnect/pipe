package x_test

import (
	"fmt"
	"testing"

	l "github.com/masteryconnect/pipe/line"
	"github.com/masteryconnect/pipe/x"
)

func ExampleCount() {
	l.New().SetP(func(out chan<- interface{}, errs chan<- error) {
		for i := 0; i < 4; i++ {
			out <- "foo"
		}
	}).Add(
		x.Count{}.T,
	).Run()
	// Output: 4
}

func ExampleCount_silent() {
	c := x.Count{Silent: true}

	l.New().SetP(func(out chan<- interface{}, errs chan<- error) {
		for i := 0; i < 4; i++ {
			out <- "foo"
		}
	}).Add(
		c.Use, // Use is on the pointer so the counted value is kept
	).Run()

	fmt.Println(c.Val())
	// Output: 4
}

func TestCount_concurrency(t *testing.T) {
	c := x.Count{Silent: true}

	l.New().SetP(func(out chan<- interface{}, errs chan<- error) {
		for i := 0; i < 4000; i++ {
			out <- i
		}
	}).Add(
		l.Many(c.Use, 100),
	).Run()

	final := c.Val()
	if final != 4000 {
		t.Errorf("want %d got %d", 4000, final)
	}
}
