package fs

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
)

// Read will read the messages from a file much like from stdin.
type Read struct {
	Path string
}

// T is the Tfunc for a pipe/line.
func (r Read) T(in <-chan interface{}, out chan<- interface{}, errs chan<- error) {
	for m := range in {
		r.run(m.(fmt.Stringer).String(), out, errs)
	}
}

// P is the producer
func (r Read) P(out chan<- interface{}, errs chan<- error) {
	r.run(r.Path, out, errs)
}

func (r Read) run(path string, out chan<- interface{}, errs chan<- error) {
	file, err := os.Open(path)
	if err != nil {
		errs <- err
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		msgSrc := scanner.Bytes()
		msg := make([]byte, len(msgSrc))
		copy(msg, msgSrc)
		out <- bytes.NewBuffer(msg)
	}

	if err := scanner.Err(); err != nil {
		errs <- err
	}
}
