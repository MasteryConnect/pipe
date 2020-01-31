package message_test

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/masteryconnect/pipe/message"
)

func TestRecordToMSI(t *testing.T) {
	msi1 := map[string]interface{}{
		"id":     42,
		"name":   "foo",
		"sub":    map[string]string{"foo": "bar"},
		"struct": struct{ id int }{42},
	}

	r := message.NewRecordFromMSI(msi1)

	msi2 := message.RecordToMSI(r)

	if !reflect.DeepEqual(msi1, msi2) {
		t.Errorf("maps not the same: %+v != %+v", msi1, msi2)
	}
}

func TestRecordToStrings(t *testing.T) {
	msi1 := map[string]interface{}{
		"id":     42,
		"name":   "foo",
		"sub":    map[string]string{"foo": "bar"},
		"struct": struct{ id int }{42},
	}

	r := message.NewRecordFromMSI(msi1).
		SetKeyOrder("name", "sub", "struct", "id")

	want := []string{"foo", fmt.Sprint(msi1["sub"]), fmt.Sprint(msi1["struct"]), "42"}
	if !reflect.DeepEqual(want, message.RecordToStrings(r)) {
		t.Errorf("want %v got %v", want, message.RecordToStrings(r))
	}
}
