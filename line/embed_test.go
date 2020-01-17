package line_test

import (
	"fmt"
	"sync"
	"testing"

	l "github.com/MasteryConnect/pipe/line"
)

func ExamplePipeline_embed() {
	// define a sub-pipeline to be embedded
	// this one just prints messages out to stdout
	subPipeline := l.New().Add(l.Stdout)

	// setup and run the main pipeline
	l.New().
		SetP(func(out chan<- interface{}, errs chan<- error) {
			out <- "foo from sub"
		}).
		Add(
			subPipeline.Embed, // embed it just like any other Tfunc
		).Run()
	// Output: foo from sub
}

func TestPipeline_Embed(t *testing.T) {
	in := make(chan interface{}, 2)
	out := make(chan interface{})
	errs := make(chan error)
	errCnt := 0
	var err error
	msgCnt := 0
	var msg string

	var wg sync.WaitGroup

	// start with two messages loaded in the buffer
	in <- "err"
	in <- "foo"
	close(in)

	// start the errs range
	wg.Add(1)
	go func() {
		defer wg.Done()
		for e := range errs {
			errCnt++
			err = e
		}
	}()

	// start the out range
	wg.Add(1)
	go func() {
		defer wg.Done()
		defer close(errs)
		for m := range out {
			msgCnt++
			msg = m.(string)
		}
	}()

	l.New().
		Add(
			l.Inline(func(m interface{}) (interface{}, error) {
				if m.(string) == "err" {
					return nil, fmt.Errorf("foo error")
				} else {
					return m, nil
				}
			}),
		).
		Embed(in, out, errs)
	close(out)
	wg.Wait()

	if msgCnt != 1 {
		t.Errorf("message count: want 1 got %d", msgCnt)
	}

	if msg != "foo" {
		t.Errorf("message error: want foo got %s", msg)
	}

	if errCnt != 1 {
		t.Errorf("error count: want 1 got %d", errCnt)
	}

	if err.Error() != "foo error" {
		t.Errorf("error: want 'foo error' got '%s'", err.Error())
	}
}
