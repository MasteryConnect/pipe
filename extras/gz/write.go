package gz

import (
	"compress/gzip"
	"fmt"
	"os"
)

// Write will write the messages to a file much like to stdout.
type Write struct {
	Path    string
	Prefix  string // add to the beginning of each message
	Postfix string // add to the end of each message (useful for adding newlines at the end)
}

// T is the Tfunc for a pipe/line.
func (w Write) T(in <-chan interface{}, out chan<- interface{}, errs chan<- error) {
	file, err := os.Create(w.Path)
	if err != nil {
		errs <- err
		return
	}
	defer file.Close()

	zw := gzip.NewWriter(file)
	defer zw.Close()
	defer zw.Flush()

	for msg := range in {
		if v, ok := msg.(fmt.Stringer); ok {
			fmt.Fprintf(zw, "%s%s%s", w.Prefix, v.String(), w.Postfix)
		} else {
			fmt.Fprintf(zw, "%s%+v%s", w.Prefix, msg, w.Postfix)
		}

		if out != nil {
			out <- msg
		}
	}
}
