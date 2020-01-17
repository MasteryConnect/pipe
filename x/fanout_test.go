package x

import (
	"testing"
)

func TestFanoutNoTfuncs(t *testing.T) {
	fo := NewFanout([]FanoutTfunc{}, nil)

	in := make(chan interface{})
	out := make(chan interface{})
	errs := make(chan error)
	go func() {
		for err := range errs {
			t.Errorf("Errors channel received the following error: %s\n", err.Error())
		}
	}()

	go fo.T(in, out, errs)

	in <- "test"

	close(in)
	close(out)
	close(errs)
}

func TestFanoutWithTfuncs(t *testing.T) {
	fo := NewFanout([]FanoutTfunc{&testTfunc1{}, &testTfunc2{includeTypes: []string{"type1", "type2"}}}, nil)

	in := make(chan interface{})
	out := make(chan interface{})
	errs := make(chan error)
	go func() {
		for err := range errs {
			t.Errorf("Errors channel received the following error: %s\n", err.Error())
		}
	}()

	go fo.T(in, out, errs)

	in <- "test"

	for i := 0; i < 2; i++ {
		outMsg := <-out
		if outMsg != "test" {
			t.Errorf("The out message is not the same as the in message sent: %s\n", outMsg)
		}
	}

	close(in)
	close(out)
	close(errs)
}

func TestFanoutWithFiltering(t *testing.T) {
	t1 := &testTfunc1{} // Get all messages
	t2 := &testTfunc2{
		includeTypes: []string{"type1", "type2"},
	}
	t3 := &testTfunc2{
		includeTypes: []string{"type2", "type3"},
	}
	fo := NewFanout([]FanoutTfunc{t1, t2, t3}, msgTypes)

	in := make(chan interface{})
	out := make(chan interface{})
	errs := make(chan error)
	go func() {
		for err := range errs {
			t.Errorf("Errors channel received the following error: %s\n", err.Error())
		}
	}()

	var outMsgs []interface{}
	go func() {
		for msg := range out {
			outMsgs = append(outMsgs, msg)
		}
	}()

	go func() {
		in <- testMsg{[]string{"type1"}}
		in <- testMsg{[]string{"type1", "type3"}}
		in <- testMsg{[]string{"type2", "type3"}}
		in <- testMsg{[]string{"type3", "type4"}}
		close(in)
	}()

	fo.T(in, out, errs)

	if len(t1.inMsgs) != 4 {
		t.Errorf("Tfunc 1 should have received all messages: %+v\n", t1.inMsgs)
	}
	if len(t2.inMsgs) != 3 {
		t.Errorf("Tfunc 2 should have received 3 messages: %+v\n", t2.inMsgs)
	}
	if len(t3.inMsgs) != 3 {
		t.Errorf("Tfunc 3 should have received 3 messages: %+v\n", t3.inMsgs)
	}

	close(out)
	close(errs)
}

type testMsg struct {
	types []string
}

func msgTypes(msg interface{}) (types []string) {
	return msg.(testMsg).types
}

type testTfunc1 struct {
	inMsgs []interface{}
}

func (t *testTfunc1) T(in <-chan interface{}, out chan<- interface{}, errs chan<- error) {
	for msg := range in {
		t.inMsgs = append(t.inMsgs, msg)
		out <- msg
	}
}

type testTfunc2 struct {
	inMsgs       []interface{}
	includeTypes []string
}

func (t *testTfunc2) T(in <-chan interface{}, out chan<- interface{}, errs chan<- error) {
	for msg := range in {
		t.inMsgs = append(t.inMsgs, msg)
		out <- msg
	}
}

func (t *testTfunc2) I() []string {
	return t.includeTypes
}
