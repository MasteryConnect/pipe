package x

import (
	"sync/atomic"

	"gopkg.in/cheggaaa/pb.v1"
)

// Progress shows the progress of the stream
type Progress struct {
	*pb.ProgressBar
}

// NewProgress creates a new progress transformer
func NewProgress(total int) *Progress {
	return &Progress{pb.New(total)}
}

// T will add each message to the progress bar progress
func (p *Progress) T(in <-chan interface{}, out chan<- interface{}, errs chan<- error) {
	p.Start()
	for m := range in {
		p.Increment()
		out <- m
	}
	p.Finish()
}

// AddToTotal will passthrough message in a stream and add the count to the total
func (p *Progress) AddToTotal(in <-chan interface{}, out chan<- interface{}, errs chan<- error) {
	for m := range in {
		atomic.AddInt64(&p.Total, 1)
		out <- m
	}
}
