package x

import (
	"sync"

	"github.com/masteryconnect/pipe/line"
)

// Tee splits a stream into multiple streams and merges the results
// Example:
//	otherpipe := line.New()
//	line.New().Add(x.Tee(otherpipe)).Run()
func Tee(targets ...line.Tfunc) line.Tfunc {
	return func(in <-chan interface{}, out chan<- interface{}, errs chan<- error) {
		var wg sync.WaitGroup
		wg.Add(len(targets))
		inChans := make([]chan interface{}, len(targets))

		for i, t := range targets {
			inChans[i] = make(chan interface{})
			go func(tin <-chan interface{}, target line.Tfunc) {
				defer wg.Done()
				target(tin, out, errs) // tee up the target
			}(inChans[i], t)
		}

		// now pass the messages along to all targets and downstream
		for m := range in {
			for _, targetIn := range inChans {
				targetIn <- m
			}
			out <- m
		}

		// our in is done so close the other ins
		for _, c := range inChans {
			close(c)
		}

		// wait for all the targets to finish now that their 'in's are closed
		wg.Wait()
	}
}
