package x_test

import (
	"bytes"
	"fmt"
	"time"

	l "github.com/masteryconnect/pipe/line"
	"github.com/masteryconnect/pipe/x"
)

func ExampleProgress_T_knownTotal() {
	// known total of 3
	pb := x.NewProgress(3)

	l.New().SetP(func(out chan<- interface{}, errs chan<- error) {
		out <- bytes.NewBufferString("1")
		out <- bytes.NewBufferString("2")
		out <- bytes.NewBufferString("3")
	}).Add(
		x.RateLimit{N: 1, Per: time.Second}.T, // do some work
		pb.T,                                  // increments the progress bar
	).Run()
}

func ExampleProgress_T_unknownTotal() {
	// unknown total so set to 0
	pb := x.NewProgress(0)

	l.New().SetP(func(out chan<- interface{}, errs chan<- error) {
		out <- bytes.NewBufferString("1")
		out <- bytes.NewBufferString("2")
		out <- bytes.NewBufferString("3")
	}).Add(
		pb.AddToTotal,                         // increments the total for the progress
		x.RateLimit{N: 1, Per: time.Second}.T, // do some work
		pb.T,                                  // increments the progress bar
	).Run()
}

func ExampleProgress_AddToTotal() {
	// unknown total so set to 0
	pb := x.NewProgress(0)

	// some message to send down stream
	msg := bytes.NewBufferString("foo")

	l.New().SetP(func(out chan<- interface{}, errs chan<- error) {
		for range [30]struct{}{} { // loop 30 times
			out <- msg
		}
	}).Add(
		pb.AddToTotal, // increments the total for the progress
	).Run()

	fmt.Println(pb.Total)
	// Output: 30
}
