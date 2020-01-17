package line_test

import (
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
