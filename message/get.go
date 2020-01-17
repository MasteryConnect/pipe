package message

import (
	"reflect"
)

// Inner defines the In function wrapping messages should implement.
type Inner interface {
	In() interface{}
}

// Get will get the desired type out of wrapped messages.
func Get(msg interface{}, v interface{}) bool {
	vType := reflect.TypeOf(v)
	return _get(msg, v, vType)
}

// _get is the recursive digging in the wrapped messages.
func _get(msg interface{}, v interface{}, vType reflect.Type) bool {
	mType := reflect.TypeOf(msg)
	pType := reflect.PtrTo(mType)
	if mType == vType || pType == vType {
		mVal := reflect.Indirect(reflect.ValueOf(msg))
		vVal := reflect.Indirect(reflect.ValueOf(v))
		vVal.Set(mVal)
		return true
	}

	if inner, ok := msg.(Inner); ok {
		return _get(inner.In(), v, vType) // recursive call (i.e. dig deeper)
	}

	return false
}
