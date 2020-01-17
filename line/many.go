package line

import (
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
