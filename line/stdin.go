package line

import (
	"bufio"
	"bytes"
	"io"
	"os"
)

// Stdin reads stdin and sends each line into the pipeline as a message.
func Stdin(out chan<- interface{}, errs chan<- error) {
	reader := bufio.NewReader(os.Stdin)
	msg := []byte("")
	line, prefix, err := reader.ReadLine()
	for err == nil && line != nil {
		msg = append(msg, line...)
		if prefix {
			line, prefix, err = reader.ReadLine()
			continue
		}
		out <- bytes.NewBuffer(msg)
		msg = []byte("")
		line, prefix, err = reader.ReadLine()
	}
	if err != nil {
		if err != io.EOF {
			errs <- err
		}
	}
}
