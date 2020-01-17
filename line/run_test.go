package line

import (
	"crypto/rand"
	"testing"
	"time"
)

func lotsOfWork(in <-chan interface{}, out chan<- interface{}, errs chan<- error) {
	for msg := range in {
		time.Sleep(1 * time.Microsecond)
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
