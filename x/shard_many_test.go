package x

import (
	l "github.com/MasteryConnect/pipe/line"
	"testing"
)

func TestShardManyNil(t *testing.T) {
	t.Run("concurrency < 1", func(t *testing.T) {
		_, err := NewShardMany(0, nil, nil)
		if err == nil {
			t.Error("With a concurrency < 1 there should be an error")
		}
	})
	t.Run("tfunc is nil", func(t *testing.T) {
		_, err := NewShardMany(2, nil, nil)
		if err == nil {
			t.Error("With a nil tfunc there should be an error")
		}
	})
	t.Run("shardManyKeyFunc is nil", func(t *testing.T) {
		tfunc := func(in <-chan interface{}, out chan<- interface{}, errs chan<- error) {}
		_, err := NewShardMany(2, tfunc, nil)
		if err == nil {
			t.Error("With a nil shardManyKeyFunc there should be an error")
		}
	})
}

func TestShardManyOneConcurrency(t *testing.T) {
	var empty struct{}
	msgs := map[string]struct{}{"test": empty}
	found := 0
	tfunc := func(in <-chan interface{}, out chan<- interface{}, errs chan<- error) {
		for msg := range in {
			if m, ok := msg.(string); ok {
				if _, f := msgs[m]; f {
					found += 1
				}
			}
			out <- msg
		}
	}
	smkf := func(msg interface{}) []byte {
		return []byte("test")
	}
	sm, err := NewShardMany(1, tfunc, smkf)
	if err != nil {
		t.Error(err)
	}

	errs := make(chan error)
	go func() {
		for err := range errs {
			t.Errorf("Errors channel received the following error: %s\n", err.Error())
		}
	}()

	l.New().SetP(func(out chan<- interface{}, errs chan<- error) {
		for msg, _ := range msgs {
			out <- msg
		}
	}).
		SetErrs(errs).
		Add(
			sm.T,
			l.Stdout,
		).Run()
	if found != 1 {
		t.Errorf("Expected to see 1 message, instead there were %d", found)
	}
}

func TestShardManyThreeConcurrency(t *testing.T) {
	var empty struct{}
	msgs := map[string]struct{}{
		"test":  empty,
		"test1": empty,
		"test3": empty,
	}
	found := 0
	tfunc := func(in <-chan interface{}, out chan<- interface{}, errs chan<- error) {
		for msg := range in {
			if m, ok := msg.(string); ok {
				if _, f := msgs[m]; f {
					found += 1
				}
			}
			out <- msg
		}
	}
	smkf := func(msg interface{}) []byte {
		return []byte(msg.(string))
	}
	sm, err := NewShardMany(1, tfunc, smkf)
	if err != nil {
		t.Error(err)
	}

	errs := make(chan error)
	go func() {
		for err := range errs {
			t.Errorf("Errors channel received the following error: %s\n", err.Error())
		}
	}()

	l.New().SetP(func(out chan<- interface{}, errs chan<- error) {
		for msg, _ := range msgs {
			out <- msg
		}
	}).
		SetErrs(errs).
		Add(
			sm.T,
			l.Stdout,
		).Run()
	if found != 3 {
		t.Errorf("Expected to see 3 messages, instead there were %d", found)
	}
}
