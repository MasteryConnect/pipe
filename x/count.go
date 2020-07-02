package x

import (
	"fmt"
	"sync/atomic"
	"time"

	"github.com/dustin/go-humanize"
)

// Count counts up the messages coming down the pipe.
// Options:
// 		Live: show the counts real-time to the console
//		Silent: don't print to the console at all
//			This is useful if you want to use the count
//			programmatically instead of for the user on the console
//		Mod: Only print counts that are modulus this number
//		Raw: don't humanize the printed value (no commas)
//		Mul: multiply the printed number by this amount
//			instead of 1 for the batch message itself
//		AutoMod: automatically determin the modulus
//			to avoid printing to the console too much
type Count struct {
	Live    bool
	Silent  bool
	Mod     int64
	Raw     bool
	Mul     int64
	AutoMod bool

	cnt       int64
	lastPrint time.Time
}

// T will count up the messages and pass them on.
func (c Count) T(in <-chan interface{}, out chan<- interface{}, errs chan<- error) {
	(&c).Use(in, out, errs)
}

// Use will count up the messages and pass them on.
func (c *Count) Use(in <-chan interface{}, out chan<- interface{}, errs chan<- error) {

	for msg := range in {
		atomic.AddInt64(&c.cnt, 1)
		if c.Live {
			if c.Mod > 0 {
				c.print(c.Mod, "\r")
			} else {
				c.print(1, "\r") // use mod 1 instead of 0 if set to 0
			}
		}

		out <- msg
	}

	// use mod = 1 to ensure we print the last value accurately
	c.print(1, "\n")
}

// Val returns the value of the counter.
func (c *Count) Val() int64 {
	if c.Mul == 0 {
		return atomic.LoadInt64(&c.cnt)
	}
	return atomic.LoadInt64(&c.cnt) * c.Mul
}

func (c *Count) print(mod int64, tail string) {
	if c.Silent {
		return
	}
	v := c.Val()
	if v%mod == 0 {
		if c.Raw {
			fmt.Printf("%d"+tail, v)
		} else {
			fmt.Printf("%s"+tail, humanize.Comma(v))
		}

		if c.AutoMod {
			c.autoMod()
		}
	}
}

func (c *Count) autoMod() {
	elapsed := time.Since(c.lastPrint)
	if elapsed < (100 * time.Millisecond) {
		c.Mod *= 10
	} else if c.Mod > 1 && elapsed > time.Second {
		c.Mod /= 10
	}
	c.lastPrint = time.Now()
}
