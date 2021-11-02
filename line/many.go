package line

import (
	"context"
	"sync"
)

// Many wraps a transformer to run it in multiple go routines.
func Many(t Tfunc, concurrency int) Tfunc {
	return func(in <-chan interface{}, out chan<- interface{}, errs chan<- error) {
		var wg sync.WaitGroup
		wg.Add(concurrency)
		for n := 0; n < concurrency; n++ {
			go func() {
				defer wg.Done()
				t(in, out, errs)
			}()
		}
		wg.Wait()
	}
}

// ManyContext wraps a transformer to run it in multiple go routines.
func ManyContext(t TfuncContext, concurrency int) TfuncContext {
	return func(ctx context.Context, in <-chan interface{}, out chan<- interface{}, errs chan<- error) {
		var wg sync.WaitGroup
		wg.Add(concurrency)
		for n := 0; n < concurrency; n++ {
			go func() {
				defer wg.Done()
				t(ctx, in, out, errs)
			}()
		}
		wg.Wait()
	}
}
