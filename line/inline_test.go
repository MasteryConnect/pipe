package line_test

import (
	"context"
	"fmt"

	l "github.com/MasteryConnect/pipe/line"
)

func ExampleInline() {
	l.New().
		SetP(func(out chan<- interface{}, errs chan<- error) {
			out <- "foo"
		}).
		Add(
			l.Inline(func(m interface{}) (interface{}, error) {
				return fmt.Sprintf("inline func says: %s", m), nil
			}),
			l.Stdout,
		).Run()
	// Output: inline func says: foo
}

func ExampleInlineContext() {
	type ctxKey string
	var key ctxKey = "somekey"
	outsideCtx, cancel := context.WithCancel(context.Background())
	outsideCtx = context.WithValue(outsideCtx, key, "Go")

	l.New().
		SetP(func(out chan<- interface{}, errs chan<- error) {
			out <- "foo"
			out <- "last"
			out <- "never"
		}).
		AddContext(
			l.InlineContext(func(ctx context.Context, m interface{}) (interface{}, error) {
				if m.(string) == "last" {
					// the InlineContext will stop if the context is Done()
					cancel() // cancel the context
				}
				return fmt.Sprintf(`inline func with context "%s" says: %s`, ctx.Value(key), m), nil
			}),
		).
		Add(l.Stdout).
		RunContext(outsideCtx) // be sure to run with context

	// Output:
	// inline func with context "Go" says: foo
	// inline func with context "Go" says: last
}
