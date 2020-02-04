package csv

import (
	"bytes"
	"encoding/csv"

	"github.com/MasteryConnect/pipe/message"
	"github.com/pkg/errors"
)

// ErrUnknownType is the error for when we don't know how to conver the message to a []string
var ErrUnknownType = errors.New("unknown type")

// To will convert incoming messages to a csv message output and send downstream
type To struct {
	ShowHeader bool     // indicate if we should show the header on the first record
	Header     []string // let the user specify the header to be shown

	headerShown bool
	csv         *csv.Writer
}

// T implements the Tfunc interface
func (t To) T(in <-chan interface{}, out chan<- interface{}, errs chan<- error) {
	var buf bytes.Buffer
	t.csv = csv.NewWriter(&buf)

	sendRow := func(m interface{}) {
		var err error
		buf.Reset()

		switch v := m.(type) {
		case message.Record:
			t.csv.Write(message.RecordToStrings(v))
		case []string:
			err = t.csv.Write(v)
		case []interface{}:
			row := make([]string, len(v))
			for i, val := range v {
				row[i] = message.String(val) // convert each array element to a string
			}
			err = t.csv.Write(row)
		case message.Stringser: // should cover things like message.Event and message.Record
			t.csv.Write(v.Strings())
		case map[string]interface{}:
			if len(t.Header) > 0 {
				rec := message.NewRecordFromMSI(v).SetKeyOrder(t.Header...)
				err = t.csv.Write(message.RecordToStrings(rec))
			} else {
				row := []string{}
				for _, val := range v {
					row = append(row, message.String(val)) // convert each array element to a string
				}
				err = t.csv.Write(row)
			}
		default:
			err = errors.Wrapf(ErrUnknownType, "got type %T", v)
		}

		if err != nil {
			errs <- err
		} else {
			t.csv.Flush()
			b := buf.Bytes()
			newB := make([]byte, len(b)-1) // -1 to trim off the "\n" at the end
			copy(newB, b)
			out <- newB
		}
	}

	for m := range in {
		if t.ShowHeader && !t.headerShown {
			t.headerShown = true
			if len(t.Header) > 0 {
				sendRow(t.Header)
			} else if r, ok := m.(message.Record); ok {
				sendRow(r.GetKeys())
			}
		}
		sendRow(m)
	}
}
