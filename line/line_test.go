package line_test

import (
	"bytes"
	"fmt"
	"sync"

	l "github.com/MasteryConnect/pipe/line"
)

func ExamplePipeline_new_w_chan() {
	out := make(chan interface{}, 1)
	out <- bytes.NewBufferString("foo")
	close(out)

	l.New(out).
		Add(l.Stdout).
		Run()
	// Output: foo
}

func ExamplePipeline_setP() {
	l.New().
		SetP(func(out chan<- interface{}, errs chan<- error) {
			out <- bytes.NewBufferString("foo")
		}).
		Add(l.Stdout).
		Run()
	// Output: foo
}

func ExamplePipeline_setC() {
	l.New().
		SetP(func(out chan<- interface{}, errs chan<- error) {
			out <- bytes.NewBufferString("foo")
		}).
		SetC(func(in <-chan interface{}, errs chan<- error) {
			for msg := range in {
				fmt.Println(msg)
			}
		}).
		Run()
	// Output: foo
}

func ExamplePipeline_add() {
	l.New().
		SetP(func(out chan<- interface{}, errs chan<- error) {
			out <- bytes.NewBufferString("foo")
		}).

		// one-off add calls
		Add(l.Inline(func(m interface{}) (interface{}, error) {
			return m.(fmt.Stringer).String() + " bar", nil
		})).
		Add(l.Stdout).

		// execute the pipeline starting with the producer
		Run()
	// Output: foo bar
}

func ExamplePipeline_add_combined() {
	l.New().
		SetP(func(out chan<- interface{}, errs chan<- error) {
			out <- bytes.NewBufferString("foo")
		}).

		// combined add call
		Add(
			l.Inline(func(m interface{}) (interface{}, error) {
				return m.(fmt.Stringer).String() + " bar", nil
			}),
			l.Stdout,
		).

		// execute the pipeline starting with the producer
		Run()
	// Output: foo bar
}

func ExamplePipeline_setErrs() {
	// create your own errs channel
	// Run() will close it for you
	errs := make(chan error)

	// make a sync.WaitGroup to give us time to drain errs
	// before Run() returns
	var errsDrained sync.WaitGroup
	errsDrained.Add(1)

	// start reading from the errs channel
	// with your custom error channel reader
	go func() {
		defer errsDrained.Done() // indicate we are done draining the errs

		// drain the errs
		for e := range errs {
			fmt.Println(e)
		}
	}()

	// setup a pipeline with custom error channel handling
	l.New().
		SetP(func(out chan<- interface{}, errs chan<- error) {
			errs <- fmt.Errorf("foo error")
		}).
		SetErrs(errs).
		Run()

	close(errs)
	errsDrained.Wait()
	// Output: foo error
}
