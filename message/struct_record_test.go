package message_test

import (
	"reflect"
	"testing"

	"github.com/MasteryConnect/pipe/message"
)

type testStruct struct {
	ID   int    `db:"id"`
	Name string `db:"name"`
}

func TestStructRecord(t *testing.T) {
	foo := testStruct{1, "foo"}

	r, err := message.NewStructRecord(foo)
	if err != nil {
		t.Error(err)
	}

	wantKeys := []string{"id", "name"}
	gotKeys := r.GetKeys()

	wantVals := []interface{}{1, "foo"}
	gotVals := r.GetVals()

	for i, key := range wantKeys {
		val, ok := r.Get(key)
		if !ok {
			t.Error("couldn't find key", key)
		}
		if !reflect.DeepEqual(val, wantVals[i]) {
			t.Errorf("want %v got %v", wantVals[i], val)
		}
	}

	if !reflect.DeepEqual(wantKeys, gotKeys) {
		t.Errorf("want %v got %v", wantKeys, gotKeys)
	}

	if !reflect.DeepEqual(wantVals, gotVals) {
		t.Errorf("want %v got %v", wantVals, gotVals)
	}
}
