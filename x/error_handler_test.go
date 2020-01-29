package x_test

import (
	"errors"
	"math/rand"
	"testing"
	"time"

	l "github.com/masteryconnect/pipe/line"
	"github.com/masteryconnect/pipe/x"
)

func TestErrorHandler(t *testing.T) {
	msgCount := 10
	taskCallCount := 0
	ehCount := 0
	eh := x.ErrorHandler{
		TaskToTry: func(msg interface{}) (interface{}, error) {
			taskCallCount++
			return nil, nil
		},
		ErrorHandler: func(in <-chan interface{}, out chan<- interface{}, errs chan<- error) {
			ehCount++
		},
	}

	l.New().SetP(func(out chan<- interface{}, errs chan<- error) {
		for i := 0; i < msgCount; i++ {
			out <- i
		}
	}).Add(
		eh.T,
	).Run()

	if taskCallCount != msgCount {
		t.Errorf("want %d got %d", msgCount, taskCallCount)
	}
}

func TestErrorHandlerWithErrorNoTaskMsg(t *testing.T) {
	msgCount := 10
	taskCallCount := 0
	ehCount := 0
	err := errors.New("TaskError")
	var ehInMsg []interface{}
	eh := x.ErrorHandler{
		TaskToTry: func(msg interface{}) (interface{}, error) {
			taskCallCount++
			return nil, err
		},
		ErrorHandler: func(in <-chan interface{}, out chan<- interface{}, errs chan<- error) {
			for msg := range in {
				ehCount++
				ehInMsg = append(ehInMsg, msg)
			}
		},
	}

	l.New().SetP(func(out chan<- interface{}, errs chan<- error) {
		for i := 0; i < msgCount; i++ {
			out <- i
		}
	}).Add(
		eh.T,
	).Run()

	if taskCallCount != msgCount {
		t.Errorf("Task call count: want %d got %d", msgCount, taskCallCount)
	}

	if ehCount != msgCount {
		t.Errorf("ErrorHandler call count: want %d got %d", msgCount, ehCount)
	}

	if len(ehInMsg) != msgCount {
		t.Errorf("ErrorHandler messages: want %d got %d", msgCount, len(ehInMsg))
	}

	if ehInMsg[0] != 0 {
		t.Errorf("ErrorHandler msg: want %d got %d", 0, ehInMsg[0])
	}
}

func TestErrorHandlerWithErrorWithTaskMsg(t *testing.T) {
	msgCount := 10
	taskMsg := "TaskMsg"
	taskCallCount := 0
	ehCount := 0
	err := errors.New("TaskError")
	var ehInMsg []interface{}
	eh := x.ErrorHandler{
		TaskToTry: func(msg interface{}) (interface{}, error) {
			taskCallCount++
			return taskMsg, err
		},
		ErrorHandler: func(in <-chan interface{}, out chan<- interface{}, errs chan<- error) {
			for msg := range in {
				ehCount++
				ehInMsg = append(ehInMsg, msg)
			}
		},
	}

	l.New().SetP(func(out chan<- interface{}, errs chan<- error) {
		for i := 0; i < msgCount; i++ {
			out <- i
		}
	}).Add(
		eh.T,
	).Run()

	if taskCallCount != msgCount {
		t.Errorf("Task call count: want %d got %d", msgCount, taskCallCount)
	}

	if ehCount != msgCount {
		t.Errorf("ErrorHandler call count: want %d got %d", msgCount, ehCount)
	}

	if len(ehInMsg) != msgCount {
		t.Errorf("ErrorHandler messages: want %d got %d", msgCount, len(ehInMsg))
	}

	if ehInMsg[0] != taskMsg {
		t.Errorf("ErrorHandler msg: want %s got %s", taskMsg, ehInMsg[0])
	}
}

func TestErrorHandlerRandomError(t *testing.T) {
	msgCount := 10
	taskMsg := "TaskMsg"
	taskCallCount := 0
	errorSendCount := 0
	ehCount := 0
	rand.Seed(time.Now().UnixNano())
	err := errors.New("TaskError")
	var ehInMsg []interface{}
	eh := x.ErrorHandler{
		TaskToTry: func(msg interface{}) (interface{}, error) {
			taskCallCount++
			// Send and error message roughly half of the time
			if rand.Intn(2) == 1 {
				errorSendCount++
				return taskMsg, err
			} else {
				return taskMsg, nil
			}
		},
		ErrorHandler: func(in <-chan interface{}, out chan<- interface{}, errs chan<- error) {
			for msg := range in {
				ehCount++
				ehInMsg = append(ehInMsg, msg)
			}
		},
	}

	l.New().SetP(func(out chan<- interface{}, errs chan<- error) {
		for i := 0; i < msgCount; i++ {
			out <- i
		}
	}).Add(
		eh.T,
	).Run()

	if taskCallCount != msgCount {
		t.Errorf("Task call count: want %d got %d", msgCount, taskCallCount)
	}

	if ehCount != errorSendCount {
		t.Errorf("ErrorHandler call count: want %d got %d", errorSendCount, ehCount)
	}

	if len(ehInMsg) != errorSendCount {
		t.Errorf("ErrorHandler messages: want %d got %d", errorSendCount, len(ehInMsg))
	}

	if ehInMsg[0] != taskMsg {
		t.Errorf("ErrorHandler msg: want %s got %s", taskMsg, ehInMsg[0])
	}
}
