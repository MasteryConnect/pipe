package message

import (
	"encoding/json"
)

// Record is a collection of key/value pairs where the keys are strings
// It is similar to a map[string]interface{} accept it keeps column order.
type Record struct {
	Keys  []string
	Vals  []interface{}
	index map[string]uint
}

// NewRecord creates a *Record empty
func NewRecord() *Record {
	return &Record{
		Keys:  []string{},
		Vals:  []interface{}{},
		index: map[string]uint{},
	}
}

// NewRecordFromMSI creates a *Record from a map[string]interface{}
// This is a convenience function to combine new and FromMSI
func NewRecordFromMSI(msi map[string]interface{}) *Record {
	r := Record{
		Keys:  []string{},
		Vals:  []interface{}{},
		index: map[string]uint{},
	}
	return r.FromMSI(msi)
}

// SetOrder sets the order of the keys. This is useful if you need to maintain column order
// for things like CSV output etc...
func (r Record) SetOrder(keys []string) *Record {
	newR := NewRecord()
	newR.Keys = keys
	newR.Vals = make([]interface{}, len(keys))
	for i, k := range keys {
		newR.index[k] = uint(i)
	}

	if len(r.Vals) > 0 {
		for k, v := range r.index {
			newR.Set(k, r.Vals[v])
		}
	}

	return newR
}

// FromMSI sets the keys and vals of the record to the data in the map[string]interface{}
func (r *Record) FromMSI(msi map[string]interface{}) *Record {
	for k, v := range msi {
		r.Set(k, v)
	}
	return r
}

// Set will add or update the key value pair
func (r *Record) Set(key string, val interface{}) *Record {
	if _, exists := r.index[key]; exists {
		r.Vals[r.index[key]] = val
		return r
	}

	r.Keys = append(r.Keys, key)
	r.Vals = append(r.Vals, val)
	r.index[key] = uint(len(r.Vals) - 1)
	return r
}

// Get will return the value for the specified key or nil.
// The returned bool indicates if the key existed or not.
func (r Record) Get(key string) (interface{}, bool) {
	if _, exists := r.index[key]; exists {
		return r.Vals[r.index[key]], true
	}
	return nil, false
}

// MSI makes a map[string]interface{} out of the record
func (r Record) MSI() map[string]interface{} {
	msi := make(map[string]interface{})
	for i, k := range r.Keys {
		if i < len(r.Vals) {
			msi[k] = r.Vals[i]
		}
	}
	return msi
}

// String implements the fmt.Stringer interface
func (r Record) String() string {
	jsn, _ := json.Marshal(r.MSI())
	return string(jsn)
}

// Strings implements the Strinsger interface
func (r Record) Strings() []string {
	row := make([]string, len(r.Vals))
	for i, val := range r.Vals {
		row[i] = String(val) // convert each array element to a string
	}
	return row
}
