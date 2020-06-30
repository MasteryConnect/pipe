package line

import (
	"context"
	"crypto/rand"
	"testing"
	"time"
)

// Test to make sure that the context from RunContext gets used in the
// producer.
func TestRunContext(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
	defer cancel() // should already be cancelled by the timeout

	p := New().SetPContext(func(ctx context.Context, out chan<- interface{}, errs chan<- error) {
		for {
			select {
			case <-ctx.Done():
				return // context was cancelled so return
			case <-time.Tick(1 * time.Millisecond):
				out <- nil
			}
		}
	})
	p.RunContext(ctx)
}

func lotsOfWork(in <-chan interface{}, out chan<- interface{}, errs chan<- error) {
	for msg := range in {
		time.Sleep(10 * time.Microsecond)
		out <- msg // passthrough
	}
}

func BenchmarkRun(b *testing.B) {
	p := New()
	body := make([]byte, 400)
	rand.Read(body)

	p.SetP(func(out chan<- interface{}, errs chan<- error) {
		b.ResetTimer()
		for n := 0; n < b.N; n++ {
			out <- body
		}
	})
	p.Add(func(in <-chan interface{}, out chan<- interface{}, errs chan<- error) {
		for msg := range in {
			out <- msg // passthrough
		}
	})
	p.SetC(NoopC)

	p.Run()
}

func BenchmarkOneRoutine(b *testing.B) {
	p := New()
	body := make([]byte, 400)
	rand.Read(body)

	p.SetP(func(out chan<- interface{}, errs chan<- error) {
		b.ResetTimer()
		for n := 0; n < b.N; n++ {
			out <- body
		}
	})
	p.Add(lotsOfWork)
	p.SetC(NoopC)

	p.Run()
}

func BenchmarkManyRoutines(b *testing.B) {
	p := New()
	body := make([]byte, 400)
	rand.Read(body)

	p.SetP(func(out chan<- interface{}, errs chan<- error) {
		b.ResetTimer()
		for n := 0; n < b.N; n++ {
			out <- body
		}
	})
	p.Add(Many(lotsOfWork, 2))
	p.SetC(NoopC)

	p.Run()
}
