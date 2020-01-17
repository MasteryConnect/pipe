package line

import (
	"fmt"
)

// Stdout prints out the message to standard out.
func Stdout(in <-chan interface{}, out chan<- interface{}, errs chan<- error) {
	for msg := range in {
		if v, ok := msg.(fmt.Stringer); ok {
			fmt.Println(v.String())
		} else {
			fmt.Printf("%s\n", msg)
		}
		if out != nil {
			out <- msg
		}
	}
}

// StdoutC prints out the message to standard out.
func StdoutC(in <-chan interface{}, errs chan<- error) {
	Stdout(in, nil, errs)
}
