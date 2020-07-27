package line

import (
	"context"
	"log"
	"os"
	"sync"
)

// Run runs the whole pipeline.
func (l *Line) Run() error {
	return l.RunContext(nil)
}

// RunContext runs the whole pipeline with context.Context.
func (l *Line) RunContext(ctx context.Context) error {
	var errswg *sync.WaitGroup
	var errs chan<- error
	if l.errs != nil {
		errs = l.errs
	} else {
		errs, errswg = makeErrors()
	}

	// setup the channel for the producer
	var in chan interface{}
	var out chan interface{}

	// make the out channel for the producer
	out = make(chan interface{})

	go l.spinUpProducer(ctx, out, errs)

	for _, t := range l.t {
		in = out
		out = make(chan interface{})

		// choose the context version first if exists
		if t.TfuncContext != nil {
			go spinUpTransformersContext(ctx, t.TfuncContext, 1, in, out, errs)
		} else if t.Tfunc != nil {
			go spinUpTransformers(t.Tfunc, 1, in, out, errs)
		}
	}

	l.c(out, errs)

	if l.errs == nil {
		// if we weren't passed the channel
		// we made it and need to close it
		safeCloseErrs(errs)
	}

	if errswg != nil {
		errswg.Wait()
	}

	return nil
}

// if p is nil, then the produer is overridden and the GetIn() must be used
// to produce messages. The returned channel also must be closed to end the pipeline.
func (l *Line) spinUpProducer(ctx context.Context, out chan interface{}, errs chan<- error) {
	if ctx != nil {
		if l.pContext != nil {
			defer safeClose(out)
			l.pContext(ctx, out, errs)
		}
	} else {
		if l.p != nil {
			defer safeClose(out)
			l.p(out, errs)
		}
	}
}

func spinUpTransformers(t Tfunc, concurrency int, in chan interface{}, out chan interface{}, errs chan<- error) {
	defer safeClose(out)

	if concurrency > 1 {
		var wg sync.WaitGroup
		wg.Add(concurrency)
		for n := 0; n < concurrency; n++ {
			go func() {
				defer wg.Done()
				t(in, out, errs)
			}()
		}
		wg.Wait()
	} else {
		t(in, out, errs)
	}
}

func spinUpTransformersContext(ctx context.Context, t TfuncContext, concurrency int, in chan interface{}, out chan interface{}, errs chan<- error) {
	defer safeClose(out)

	if concurrency > 1 {
		var wg sync.WaitGroup
		wg.Add(concurrency)
		for n := 0; n < concurrency; n++ {
			go func() {
				defer wg.Done()
				t(ctx, in, out, errs)
			}()
		}
		wg.Wait()
	} else {
		t(ctx, in, out, errs)
	}
}

func makeErrors() (chan<- error, *sync.WaitGroup) {
	var wg sync.WaitGroup
	errs := make(chan error)
	wg.Add(1)

	// start the errors loop
	stderr := log.New(os.Stderr, "", 0)
	go func(errs chan error) {
		defer wg.Done()
		defer safeCloseErrs(errs)
		for err := range errs {
			stderr.Println(err)
		}
	}(errs)

	return errs, &wg
}

func safeClose(ch chan<- interface{}) {
	defer func() { recover() }()
	close(ch)
}
func safeCloseErrs(ch chan<- error) {
	defer func() { recover() }()
	close(ch)
}
