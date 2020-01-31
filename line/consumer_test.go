package line_test

import (
	"testing"

	l "github.com/MasteryConnect/pipe/line"
)

type ackable struct {
	f func()
}

func (a ackable) Ack() {
	a.f()
}

func TestConsumer(t *testing.T) {
	in := make(chan interface{})
	cnt := 0

	go func() {
		defer close(in)
		in <- "foo"
		in <- ackable{func() { cnt++ }}
		in <- "baz"
	}()

	l.Consumer(in, nil)

	if cnt != 1 {
		t.Error("Acker message didn't 'Ack'")
	}
}
