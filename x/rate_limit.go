package x

import (
	"time"
)

// RateLimit limits how many messages and go through in a time frame.
type RateLimit struct {
	N   int64
	Per time.Duration

	// smooth out the rate instead of bursting
	Smooth bool
}

// T is the Tfunc for the rate limiter
func (rl RateLimit) T(in <-chan interface{}, out chan<- interface{}, errs chan<- error) {
	if rl.N == 0 {
		rl.N = 1 // default to 1
	}

	burstCount := int64(0)
	rate := rl.Per.Nanoseconds()
	if rl.Smooth && rl.N > 1 {
		rate = rl.Per.Nanoseconds() / rl.N
		rl.N = 1
	}
	lastTime := time.Now()

	limiter := func(m interface{}, t time.Time) {
		diff := t.Sub(lastTime).Nanoseconds()
		if diff < rate {
			time.Sleep(time.Duration(rate - diff))
		}
		lastTime = time.Now()
		out <- m
	}

	for m := range in {
		burstCount++
		if burstCount < rl.N {
			out <- m
		} else {
			burstCount = 0
			limiter(m, time.Now())
		}
	}
}
