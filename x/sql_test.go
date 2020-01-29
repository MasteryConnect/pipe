package x_test

import (
	"fmt"
	"testing"

	l "github.com/masteryconnect/pipe/line"
	"github.com/masteryconnect/pipe/message"
	"github.com/masteryconnect/pipe/x"
)

func ExampleSQL_T() {
	l.New().SetP(func(out chan<- interface{}, errs chan<- error) {
		out <- message.NewRecordFromMSI(map[string]interface{}{
			"foo": "bar",
		})
	}).Add(
		x.SQL{Table: "foo"}.T,
		l.Stdout,
	).Run()

	// Output: INSERT INTO foo (foo) VALUES ('bar')
}

func ExampleSQL_I() {
	l.New().SetP(func(out chan<- interface{}, errs chan<- error) {
		out <- message.NewRecordFromMSI(map[string]interface{}{
			"foo": "bar",
		})
	}).Add(
		l.I(x.SQL{Table: "foo"}.I),
		l.Stdout,
	).Run()

	// Output: INSERT INTO foo (foo) VALUES ('bar')
}

func ExampleSQL_T_batch() {
	l.New().SetP(func(out chan<- interface{}, errs chan<- error) {
		out <- message.NewRecordFromMSI(map[string]interface{}{
			"foo": "bar",
		})
		out <- message.NewRecordFromMSI(map[string]interface{}{
			"foo": "bar2",
		})
	}).Add(
		x.Batch{N: 2}.T,
		x.SQL{Table: "foo"}.T,
		l.Stdout,
	).Run()

	// Output: INSERT INTO foo (foo) VALUES ('bar'),('bar2')
}

func TestSQL_I(t *testing.T) {
	s := x.SQL{Table: "foo"}

	record := message.NewRecordFromMSI(map[string]interface{}{
		"foo": "bar",
	})

	sql, err := s.I(record)

	if err != nil {
		t.Error(err)
	}

	ssql := sql.(fmt.Stringer).String()
	if ssql != "INSERT INTO foo (foo) VALUES ('bar')" {
		t.Error(ssql)
	}
}
