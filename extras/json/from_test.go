package json

import (
	"encoding/json"
	"strings"
	"testing"
)

func TestChunkStream(t *testing.T) {
	// offsets:                    |2       |11                    |34              x4                     |76
	stream := strings.NewReader(`  {"id":1} {"id":2,"name":"foo\""}{"id":3, "name":"💩 { test\"this out\\"}{"id":4}`)
	in := make(chan interface{})     // make a buffered channel
	out := make(chan interface{}, 4) // make a buffered channel
	errs := make(chan error, 3)      // make a buffered channel

	go func() {
		defer close(in)
		in <- stream // send the stream in to be chunked
	}()

	ChunkStream(in, out, errs) // run it

	if len(errs) > 0 {
		err := <-errs
		t.Errorf("Received errors chunking stream with first error of: %s\n", err.Error())
		return
	}

	if len(out) != 4 {
		t.Errorf("Expected 4 blocks, only got %d", len(out))
		return
	}

	type Foo struct {
		ID   int `json:"id"`
		Name string
	}

	// Block 1
	block := (<-out).(StreamChunk)
	if block.Offset != 2 {
		t.Errorf("Expected block1 offset to be 2 but got (%d)", block.Offset)
	}
	var f Foo
	err := json.Unmarshal(block.Bytes, &f)
	if err != nil {
		t.Errorf("Received an error parsing the chunked json block1: %s\n", err.Error())
	}
	if f.ID != 1 {
		t.Errorf("Expected ID of block1 to be 1 but got (%d)", f.ID)
	}

	// Block 2
	block = (<-out).(StreamChunk)
	if block.Offset != 11 {
		t.Errorf("Expected block2 offset to be 11 but got (%d)", block.Offset)
	}
	err = json.Unmarshal(block.Bytes, &f)
	if err != nil {
		t.Errorf("Received an error parsing the chunked json block2: %s\n", err.Error())
	}
	if f.ID != 2 {
		t.Errorf("Expected ID of block2 to be 2 but got (%d)", f.ID)
	}

	// Block 3
	block = (<-out).(StreamChunk)
	if block.Offset != 34 {
		t.Errorf("Expected block3 offset to be 19 but got (%d)", block.Offset)
	}
	err = json.Unmarshal(block.Bytes, &f)
	if err != nil {
		t.Errorf("Received an error parsing the chunked json block3: %s\n", err.Error())
	}
	if f.ID != 3 {
		t.Errorf("Expected ID of block3 to be 3 but got (%d)", f.ID)
	}
	if f.Name != "💩 { test\"this out\\" {
		t.Errorf("Expected Name of block3 to be (💩 { test\"this out\\) but got (%s)", f.Name)
	}

	// Block 4
	block = (<-out).(StreamChunk)
	if block.Offset != 76 {
		t.Errorf("Expected block4 offset to be 61 but got (%d)", block.Offset)
	}
}
