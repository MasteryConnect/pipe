package message_test

import (
	"testing"

	"github.com/masteryconnect/pipe/message"
)

func TestRecord_Set(t *testing.T) {
	r := message.NewRecord()

	if len(r.GetKeys()) != 0 {
		t.Errorf("want no keys got %d", len(r.GetKeys()))
	}

	r.Set("foo", "bar")

	if len(r.GetKeys()) != 1 {
		t.Errorf("want 1 key got %d", len(r.GetKeys()))
	}
	if len(r.GetVals()) != 1 {
		t.Errorf("want 1 val got %d", len(r.GetVals()))
	}
	if r.GetKeys()[0] != "foo" {
		t.Errorf("want %s key got %s", "foo", r.GetKeys()[0])
	}
	if r.GetVals()[0] != "bar" {
		t.Errorf("want %s key got %s", "bar", r.GetKeys()[0])
	}

	r.Set("baz", 42)

	if len(r.GetKeys()) != 2 {
		t.Errorf("want 2 keys got %d", len(r.GetKeys()))
	}
	if len(r.GetVals()) != 2 {
		t.Errorf("want 2 vals got %d", len(r.GetVals()))
	}
	if r.GetKeys()[1] != "baz" {
		t.Errorf("want %s key got %s", "foo", r.GetKeys()[1])
	}
	if r.GetVals()[1] != 42 {
		t.Errorf("want %s key got %s", "bar", r.GetKeys()[1])
	}
}
