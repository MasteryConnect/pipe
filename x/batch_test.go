package x_test

import (
	"fmt"
	"testing"
	"time"

	l "github.com/masteryconnect/pipe/line"
	"github.com/masteryconnect/pipe/x"
	"github.com/masteryconnect/pipe/message"
)

func ExampleBatch_T_noTimeout() {
	l.New().SetP(func(out chan<- interface{}, errs chan<- error) {
		for i := 0; i < 4; i++ {
			out <- "foo"
		}
	}).Add(
		x.Batch{N: 3}.T,
		l.Inline(func(m interface{}) (interface{}, error) {
			if b, ok := m.(message.Batch); ok {
				return fmt.Sprintf("batch of size %d", b.Size()), nil
			}
			return fmt.Sprintf("want *x.BatchMsg bot %T", m), nil
		}),
		l.Stdout,
	).Run()
	// Output:
	// batch of size 3
	// batch of size 1
}

func ExampleBatch_T_withTimeout() {
	l.New().SetP(func(out chan<- interface{}, errs chan<- error) {
		for i := 0; i < 2; i++ {
			out <- "foo"
		}

		// cause the producer to delay to allow the
		// timeout to trigger a batch that isn't full
		time.Sleep(10 * time.Millisecond)

		for i := 0; i < 2; i++ {
			out <- "foo"
		}
	}).Add(
		x.Batch{N: 3, Timeout: 9*time.Millisecond}.T,
		l.Inline(func(m interface{}) (interface{}, error) {
			if b, ok := m.(message.Batch); ok {
				return fmt.Sprintf("batch of size %d", b.Size()), nil
			}
			return fmt.Sprintf("want *x.BatchMsg bot %T", m), nil
		}),
		l.Stdout,
	).Run()
	// Output:
	// batch of size 2
	// batch of size 2
}

func ExampleBatch_T_withByteLimit() {
	l.New().SetP(func(out chan<- interface{}, errs chan<- error) {
		for i := 0; i < 10; i++ {
			out <- "foo"
		}
	}).Add(
		x.Batch{N: 4, ByteLimit: 10}.T,
		l.Inline(func(m interface{}) (interface{}, error) {
			if b, ok := m.(message.Batch); ok {
				return fmt.Sprintf("batch of size %d", b.Size()), nil
			}
			return fmt.Sprintf("want *x.BatchMsg bot %T", m), nil
		}),
		l.Stdout,
	).Run()
	// Output:
	// batch of size 3
	// batch of size 3
	// batch of size 3
	// batch of size 1
}

func TestBatch_T_tailBatch(t *testing.T) {
	l.New().SetP(func(out chan<- interface{}, errs chan<- error) {
		for i := 0; i < 2; i++ {
			out <- "foo"
		}
	}).Add(
		x.Batch{N: 3, Timeout: 9*time.Millisecond}.T,
		l.Inline(func(m interface{}) (interface{}, error) {
			b := m.(message.Batch)
			if b.Size() != 2 {
				t.Errorf("want batch size of 2 got %d", b.Size())
			}
			return nil, nil
		}),
	).Run()
}

func TestBatch_T_closeBatch(t *testing.T) {
	batch := x.CloseableBatch(3, 9*time.Millisecond, 0)
	l.New().SetP(func(out chan<- interface{}, errs chan<- error) {
		for i := 0; i < 2; i++ {
			out <- "foo"
		}
		// force the current batch to close and send downstream
		// due to external reasons
		batch.Close()
		for i := 0; i < 2; i++ {
			out <- "foo"
		}
	}).Add(
		batch.T,
		l.Inline(func(m interface{}) (interface{}, error) {
			b := m.(message.Batch)
			if b.Size() != 2 {
				t.Errorf("want batch size of 2 got %d", b.Size())
			}
			return nil, nil
		}),
	).Run()
}
