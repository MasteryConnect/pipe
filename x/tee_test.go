package x_test

import (
	"sync"
	"testing"

	"github.com/MasteryConnect/pipe/x"
)

func TestTee(t *testing.T) {
	// implement a line.Tfunc
	target1 := func(in <-chan interface{}, out chan<- interface{}, errs chan<- error) {
		for m := range in {
			if m.(string) != "foo" {
				t.Errorf("got %v want foo", m)
			}
			out <- "bar" // pass "bar" on
		}
	}

	in := make(chan interface{})
	out := make(chan interface{})
	errs := make(chan error)
	var wg sync.WaitGroup
	wg.Add(2) // out and errs

	// errs
	go func() {
		defer wg.Done()
		for e := range errs {
			t.Error(e)
		}
	}()

	// out
	output := []interface{}{}
	go func() {
		defer wg.Done()
		for m := range out {
			output = append(output, m)
		}
	}()

	go func() {
		defer close(errs)
		defer close(out)
		x.Tee(target1)(in, out, errs)
	}()

	in <- "foo"
	close(in)

	wg.Wait()

	// make assertions
	if len(output) != 2 {
		t.Errorf("got %d want 2", len(output))
		return
	}

	if output[0] == "foo" && output[1] != "bar" {
		t.Errorf("got %v want bar", output[1])
	} else if output[0] == "bar" && output[1] != "foo" {
		t.Errorf("got %v want foo", output[1])
	} else if output[0] != "foo" && output[0] != "bar" {
		t.Errorf("got %v want [foo bar] or [bar foo]", output)
	}
}
