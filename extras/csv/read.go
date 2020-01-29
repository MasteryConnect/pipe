package csv

import (
	"encoding/csv"
	"fmt"
	"io"
	"strings"
	"sync"

	"github.com/masteryconnect/pipe/message"
)

// ReadStream will read the string messages
func ReadStream(delim rune) func(<-chan interface{}, chan<- interface{}, chan<- error) {
	if delim == rune(0) {
		delim = ','
	}
	return func(in <-chan interface{}, out chan<- interface{}, errs chan<- error) {
		pr, pw := io.Pipe()

		var header []string

		r := csv.NewReader(pr)
		r.Comma = delim
		r.ReuseRecord = true
		r.LazyQuotes = true

		var wg sync.WaitGroup
		wg.Add(1)

		go func() {
			defer wg.Done()
			defer pr.Close()

			for {
				row, err := r.Read()
				if err == io.EOF {
					break
				}
				if err != nil {
					errs <- err
					return
				}

				if header == nil {
					header = append([]string{}, row...)
				} else {
					rec := message.NewRecord()
					for i, v := range header {
						rec.Set(v, row[i])
					}
					out <- rec
				}
			}
		}()

		for msg := range in {
			fmt.Fprintf(pw, "%s\n", msg.(fmt.Stringer).String())
		}

		pw.Close()
		wg.Wait() // let things drain
	}
}

// Read will create a new csv reader per message
func Read(delim rune) func(<-chan interface{}, chan<- interface{}, chan<- error) {
	if delim == rune(0) {
		delim = ','
	}
	return func(in <-chan interface{}, out chan<- interface{}, errs chan<- error) {
		var header []string

		for m := range in {
			r := csv.NewReader(strings.NewReader(m.(fmt.Stringer).String()))
			r.Comma = delim
			//r.ReuseRecord = true
			r.LazyQuotes = true

			records, err := r.ReadAll()
			if err != nil {
				errs <- err
			}

			for _, rec := range records {
				if header == nil {
					header = rec
				} else {
					row := message.NewRecord()
					for i, v := range header {
						row.Set(v, rec[i])
					}
					out <- row
				}
			}

		}

	}
}
