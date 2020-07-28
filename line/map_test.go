package line_test

import (
	"context"
	"testing"

	"github.com/MasteryConnect/pipe/line"
)

type foo struct {
	ID   int
	Name string
}

func TestMap(t *testing.T) {
	ctx := context.Background()

	check := func(fn, want interface{}) {
		in := make(chan interface{}, 1)
		out := make(chan interface{}, 1)
		errs := make(chan error, 1)
		defer close(in)
		defer close(out)
		defer close(errs)

		in <- want
		go line.Map(fn)(ctx, in, out, errs)
		got := <-out

		if got != want {
			t.Errorf("want %v got %v", want, got)
		}
	}

	t.Run("InlineTfunc", func(t *testing.T) { // make sure it is backwards compatible with the InlineTfunc
		check(func(msg interface{}) (interface{}, error) {
			return msg, nil
		}, "foo")
	})

	t.Run("string arg", func(t *testing.T) {
		check(func(msg string) (string, error) {
			return msg, nil
		}, "foo")
	})

	t.Run("int arg", func(t *testing.T) {
		check(func(msg int) (int, error) {
			return msg, nil
		}, 42)
	})

	t.Run("struct arg", func(t *testing.T) {
		check(func(msg foo) (foo, error) {
			return msg, nil
		}, foo{42, "bar"})
	})

	t.Run("pointer arg", func(t *testing.T) {
		check(func(msg *foo) (*foo, error) {
			return msg, nil
		}, &foo{42, "bar"})
	})

	t.Run("with context", func(t *testing.T) {
		check(func(ctx context.Context, msg *foo) (*foo, error) {
			return msg, nil
		}, &foo{42, "bar"})
	})

	t.Run("type mismatch", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Error("want panic but didn't happen")
			}
		}()

		// mismatched type between the want message and the msg arg string -> int
		fn := func(msg int) (int, error) {
			return msg, nil
		}
		want := "foo"

		in := make(chan interface{}, 1)
		out := make(chan interface{}, 1)
		errs := make(chan error, 1)
		defer close(in)
		defer close(out)
		defer close(errs)

		in <- want
		line.Map(fn)(ctx, in, out, errs)
	})

	t.Run("wrong shape", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Error("want panic but didn't happen")
			}
		}()

		line.Map(func(msg, msg2 int) (int, error) {
			return msg, nil
		})
	})
}
