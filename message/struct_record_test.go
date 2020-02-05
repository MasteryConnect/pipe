package message_test

import (
	"reflect"
	"testing"

	"github.com/MasteryConnect/pipe/message"
)

type testStructWithTags struct {
	ID   int    `db:"id"`
	Name string `db:"name"`
}

type testStructWithoutTags struct {
	ID   int
	Name string
}

func TestStructRecord(t *testing.T) {
	t.Run("with tags", func(t *testing.T) {
		foo := testStructWithTags{1, "foo"}

		r, err := message.NewStructRecord(foo, "db")
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
	})

	t.Run("without tags", func(t *testing.T) {
		foo := testStructWithoutTags{1, "foo"}

		r, err := message.NewStructRecord(foo)
		if err != nil {
			t.Error(err)
		}

		wantKeys := []string{"ID", "Name"}
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
	})
}
