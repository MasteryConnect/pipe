package line

import (
	"fmt"
	"reflect"
)

// ErrMapArgWrongShape is the error returned when the func shape isn't correct.
var ErrMapArgWrongShape = fmt.Errorf("a func of shape func(<in>) (<out>,error) is required as the arg")

var errorType = reflect.TypeOf((*error)(nil)).Elem()

// Map is the same as ForEach except that is also
// sends the resulting value on as the new value for this message.
// If a nil value is returned, no message will be pass along.
// The passed fund needs to be of the shape
//		func(<in>) (<out>, error)
func Map(fn interface{}) Tfunc {
	msgIdx, errIdx, err := validateMapArgType(fn)
	if err != nil {
		panic(err)
	}

	fnv := reflect.ValueOf(fn)

	return func(in <-chan interface{}, out chan<- interface{}, errs chan<- error) {
		for msg := range in {
			res := fnv.Call([]reflect.Value{reflect.ValueOf(msg)})

			// examine the error response
			if errIdx >= 0 {
				err := res[errIdx].Interface()
				if err != nil {
					errs <- err.(error)
				}
			}

			// send the new message on instead of the original message
			// and filter out if nil
			if msgIdx >= 0 {
				newMsg := res[0].Interface()
				if newMsg != nil {
					out <- newMsg
				}
			} else {
				// if there is no output message in the func signature, just send on the original message
				out <- msg
			}
		}
	}
}

// ForEach is a wrapper to Map for code readability
func ForEach(fn interface{}) Tfunc {
	return Map(fn)
}

// Filter is a wrapper to Map for code readability
func Filter(fn interface{}) Tfunc {
	return Map(fn)
}

// validateMapArgType checks the shape of the func and
// returns the position of the message and error results
// if there are in the func signature.
func validateMapArgType(fn interface{}) (int, int, error) {
	t := reflect.TypeOf(fn)

	if t.Kind() != reflect.Func {
		return -1, -1, ErrMapArgWrongShape
	}

	if t.NumIn() != 1 {
		return -1, -1, ErrMapArgWrongShape
	}

	switch t.NumOut() {

	case 2:
		if t.Out(1) != errorType {
			return -1, -1, ErrMapArgWrongShape
		}
		return 0, 1, nil

	case 1:
		if t.Out(0) == errorType {
			return -1, 0, nil // only error output
		}
		return 0, -1, nil // only message output

	case 0:
		return -1, -1, nil // no output and no error, no problem
	}

	return -1, -1, ErrMapArgWrongShape
}
