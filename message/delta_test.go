package message_test

import (
	"reflect"
	"testing"
	"time"

	"github.com/masteryconnect/pipe/message"
)

func TestInsertDelta(t *testing.T) {
	vals := []interface{}{1, "quux", time.Now()}
	r := message.NewRecord()
	r.Set("id", vals[0])
	r.Set("name", vals[1])
	r.Set("created_at", vals[2])

	d := message.NewInsertDelta(r, "foo")

	want := "INSERT INTO foo (id,name,created_at) VALUES (?,?,?)"
	if d.GetSQL() != want {
		t.Errorf("want '%s' got '%s'", want, d.GetSQL())
	}

	if !reflect.DeepEqual(vals, d.GetVals()) {
		t.Errorf("want '%s' got '%s'", vals, d.GetVals())
	}
}

func TestUpdateDelta(t *testing.T) {
	vals := []interface{}{1, "quux", 2, time.Now()}
	r := message.NewIDRecord("id")
	r.Set("id", vals[0])
	r.Set("name", vals[1])
	r.Set("version", vals[2])
	r.Set("created_at", vals[3])

	d := message.NewUpdateDelta(r, "foo")

	want := "UPDATE foo SET name=?, version=?, created_at=? WHERE id=?"
	if d.GetSQL() != want {
		t.Errorf("want '%s' got '%s'", want, d.GetSQL())
	}

	wantVals := []interface{}{vals[1], vals[2], vals[3]}
	if !reflect.DeepEqual(wantVals, d.GetArgs()) {
		t.Errorf("want '%s' got '%s'", vals, d.GetVals())
	}

	// test compsite keys
	r.(*message.BasicIDRecord).IDKeys = []string{"id", "version"}
	want = "UPDATE foo SET name=?, created_at=? WHERE id=? AND version=?"
	if d.GetSQL() != want {
		t.Errorf("want '%s' got '%s'", want, d.GetSQL())
	}
}
