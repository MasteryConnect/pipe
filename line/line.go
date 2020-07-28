// Package line is a pipeline framework inspired by unix pipes and stream processing.
//
// A pipeline is comprised of a producer, 0 or more transformers, and a consumer.
// The default producer is line-by-line reading of STDIN. The default consumer is
// a noop.
//
package line

import "sync"

// tfuncEnum holds either a Tfunc or a TfuncContext
// and can be in a slice as either one
type tfuncEnum struct {
	Tfunc
	TfuncContext
}

// Line is the order of the steps in the pipe to make a pipeline.
type Line struct {
	p        Pfunc
	pContext PfuncContext
	t        []tfuncEnum
	c        Cfunc

	errs   chan<- error
	errswg *sync.WaitGroup
}

// SetP will add the producer to the pipeline.
func (l *Line) SetP(f Pfunc) Pipeline {
	if f != nil {
		l.p = f
	}
	return l // allow chaining
}

// SetPContext will add the context aware producer to the pipeline.
// This will override the Pfunc set with SetP.
func (l *Line) SetPContext(f PfuncContext) Pipeline {
	if f != nil {
		l.pContext = f
	}
	return l // allow chaining
}

// Add will add a transformer to the pipeline.
func (l *Line) Add(f ...Tfunc) Pipeline {
	if f != nil {
		for _, fn := range f {
			l.t = append(l.t, tfuncEnum{Tfunc: fn})
		}
	}
	return l // allow chaining
}

// AddContext is like Add but with a context.Context
func (l *Line) AddContext(f ...TfuncContext) Pipeline {
	if f != nil {
		for _, fn := range f {
			l.t = append(l.t, tfuncEnum{TfuncContext: fn})
		}
	}
	return l // allow chaining
}

// SetC will add the consumer to the pipeline.
func (l *Line) SetC(f Cfunc) Pipeline {
	if f != nil {
		l.c = f
	}
	return l // allow chaining
}

// SetErrs  will set the errs channel to the pipeline.
// This can be used to hijack the errors behavior.
func (l *Line) SetErrs(errs chan<- error) Pipeline {
	if errs != nil {
		l.errs = errs
	}
	return l // allow chaining
}

// Filter is syntactic sugar around the Filter transformer
func (l *Line) Filter(fn interface{}) Pipeline {
	return l.AddContext(ForEach(fn))
}

// ForEach is syntactic sugar around the ForEach transformer
func (l *Line) ForEach(fn interface{}) Pipeline {
	return l.AddContext(ForEach(fn))
}

// Map is syntactic sugar around the ForEach transformer
func (l *Line) Map(fn interface{}) Pipeline {
	return l.AddContext(ForEach(fn))
}

// New creates a new pipeline from the built-in line package.
func New(in ...<-chan interface{}) Pipeline {
	p := Stdin

	// if we got an "in" channel, use it as the producer
	if len(in) > 0 {
		p = func(out chan<- interface{}, errs chan<- error) {
			for m := range in[0] {
				out <- m
			}
		}
	}
	return &Line{p: p, c: Consumer}
}
