package message_test

import (
	"testing"
	"bytes"

	"github.com/masteryconnect/pipe/message"
)

type wrap struct {
	M interface{}
}

func (w wrap) In() interface{} {
	return w.M
}

func TestString(t *testing.T) {
	testVals := map[interface{}]interface{}{
		"foo": "foo",
		"bar": wrap{"bar"},
		"baz": bytes.NewBufferString("baz"),
		"quux": message.Batch{"quux"},
		"fnord": message.Batch{wrap{"fnord"}},
	}

	for want, v := range testVals {
		if message.String(v) != want {
			t.Errorf("want %s got %s", want, message.String(v))
		}
	}
}
