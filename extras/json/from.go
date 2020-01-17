package json

import (
	"encoding/json"
	"fmt"
	"io"
	"reflect"
	"strings"

	"github.com/pkg/errors"
)

// From converts the json message to a map string interface
func From(msg interface{}) (interface{}, error) {
	var val interface{}
	err := json.Unmarshal([]byte(msg.(fmt.Stringer).String()), &val)
	return val, err
}

// FromAs converts the json message to an instance of the type of the passed pointer
func FromAs(ptr interface{}) func(interface{}) (interface{}, error) {
	return func(msg interface{}) (interface{}, error) {
		err := json.Unmarshal([]byte(msg.(fmt.Stringer).String()), ptr)
		if err != nil {
			return nil, err
		}
		return reflect.Indirect(reflect.ValueOf(ptr)).Interface(), err
	}
}

// ChunkStream reads out multiple json blocks from a stream of json blocks
// turning each block into it's own downstream message.
func ChunkStream(in <-chan interface{}, out chan<- interface{}, errs chan<- error) {
	for m := range in {
		send := func(offset int, block []byte) {
			out <- StreamChunk{Src: m, Offset: offset, Bytes: block}
		}
		var err error
		switch v := m.(type) {
		case string:
			err = readBlocks(strings.NewReader(v), send)
		case io.RuneReader:
			err = readBlocks(v, send)
		case fmt.Stringer:
			err = readBlocks(strings.NewReader(v.String()), send)
		}
		if err != nil {
			errs <- errors.Wrap(err, "Error while chunking stream.")
		}
	}
}

func readBlocks(reader io.RuneReader, found func(int, []byte)) error {
	var err error
	//var block []byte
	offset := 0
	cnt := 0
	for err == nil {
		blockOffset, block, _ := readBlock(reader)
		if len(block) > 0 {
			found(offset+blockOffset, block)
			offset += blockOffset + len(block) // update the offset for the next read
			cnt++
		}
		if len(block) == 0 {
			return err
		}
	}
	return err
}

func readBlock(reader io.RuneReader) (int, []byte, error) {
	var err error
	var r rune
	depth := 0
	offset := 0
	runes := []rune{}
	inSingleQuote := false
	inDoubleQuote := false
	escapeNext := false
	for err == nil {
		r, _, err = reader.ReadRune()
		if err != nil {
			break // hit an error so bail
		}

		if depth == 0 && r != '{' {
			offset++ // this is padding to track in offset
			continue
		}
		runes = append(runes, r)

		switch r {
		case '"':
			if !inSingleQuote && !escapeNext {
				inDoubleQuote = !inDoubleQuote
			}
			escapeNext = false
		case '\'':
			if !inDoubleQuote && !escapeNext {
				inSingleQuote = !inSingleQuote
			}
			escapeNext = false
		case '\\':
			if inSingleQuote || inDoubleQuote {
				escapeNext = !escapeNext
			}
		default:
			escapeNext = false // reset the escaping if nowthing was escaped
		}

		// skip evaluating the depth while in a quoted string
		if inSingleQuote || inDoubleQuote {
			continue
		}

		if r == '{' {
			depth++
		} else if r == '}' {
			depth--
			if depth == 0 {
				break // we found the entire block so drop out and return
			}
		}
	}
	return offset, []byte(string(runes)), err
}
