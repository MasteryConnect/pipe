package message_test

import (
	"testing"

	"github.com/masteryconnect/pipe/message"
)

func TestRecord_Set(t *testing.T) {
	r := message.NewRecord()

	if len(r.Keys) != 0 {
		t.Errorf("want no keys got %d", len(r.Keys))
	}


	r.Set("foo", "bar")

	if len(r.Keys) != 1 {
		t.Errorf("want 1 key got %d", len(r.Keys))
	}
	if len(r.Vals) != 1 {
		t.Errorf("want 1 val got %d", len(r.Vals))
	}
	if r.Keys[0] != "foo" {
		t.Errorf("want %s key got %s", "foo", r.Keys[0])
	}
	if r.Vals[0] != "bar" {
		t.Errorf("want %s key got %s", "bar", r.Keys[0])
	}

	r.Set("baz", 42)

	if len(r.Keys) != 2 {
		t.Errorf("want 2 keys got %d", len(r.Keys))
	}
	if len(r.Vals) != 2 {
		t.Errorf("want 2 vals got %d", len(r.Vals))
	}
	if r.Keys[1] != "baz" {
		t.Errorf("want %s key got %s", "foo", r.Keys[1])
	}
	if r.Vals[1] != 42 {
		t.Errorf("want %s key got %s", "bar", r.Keys[1])
	}
}
