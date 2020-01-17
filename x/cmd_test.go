package x_test

import (
	l "github.com/MasteryConnect/pipe/line"
	"github.com/MasteryConnect/pipe/x"
)

func ExampleCmd_T() {
	l.New().SetP(func(out chan<- interface{}, errs chan<- error) {
		out <- nil
		out <- nil
	}).Add(
		x.Command("echo", "foo").T,
		l.Stdout,
	).Run()
	// output:
	// foo
	// foo
}

func ExampleCmd_T_template_name() {
	l.New().SetP(func(out chan<- interface{}, errs chan<- error) {
		out <- "ls"
		out <- "echo"
	}).Add(
		x.Command("{{.}}", "foo").T,
		l.Stdout,
	).Run()
	// output: foo
}

func ExampleCmd_T_template_args() {
	l.New().SetP(func(out chan<- interface{}, errs chan<- error) {
		out <- "ls"
		out <- "echo"
	}).Add(
		x.Command("echo", "{{.}}").T,
		l.Stdout,
	).Run()
	// output:
	// ls
	// echo
}
