package message_test

import (
	"testing"

	"github.com/MasteryConnect/pipe/message"
)

func TestBasicRecord_Set(t *testing.T) {
	r := message.NewBasicRecord()

	assert := func(record message.Record, count, pos int, key string, val interface{}) {
		if len(record.GetKeys()) != count {
			t.Errorf("want %d key got %d", count, len(record.GetKeys()))
		}
		if len(record.GetVals()) != count {
			t.Errorf("want %d val got %d", count, len(record.GetVals()))
		}
		if count == 0 {
			return
		}
		if record.GetKeys()[pos] != key {
			t.Errorf("want %s key got %s", key, record.GetKeys()[pos])
		}
		if record.GetVals()[pos] != val {
			t.Errorf("want %v key got %v", val, record.GetVals()[pos])
		}
	}

	assert(r, 0, 0, "", nil)

	r.Set("foo", "bar")
	assert(r, 1, 0, "foo", "bar")

	r.Set("baz", 42)
	assert(r, 2, 0, "foo", "bar")
	assert(r, 2, 1, "baz", 42)

	strct := struct{ id int }{42}
	r.Set("quux", strct)
	assert(r, 3, 0, "foo", "bar")
	assert(r, 3, 1, "baz", 42)
	assert(r, 3, 2, "quux", strct)

	newr := r.SetKeyOrder("quux", "foo", "baz")
	assert(newr, 3, 0, "quux", strct)
	assert(newr, 3, 1, "foo", "bar")
	assert(newr, 3, 2, "baz", 42)
}
