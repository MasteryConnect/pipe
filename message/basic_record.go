package message

import (
	"encoding/json"
)

//
// OrderedRecord implementation
//

// BasicRecord is a collection of key/value pairs where the keys are strings.
// It is similar to a map[string]interface{} accept it keeps column order
// and implements the OrderedRecord interface. Use the NewBasicRecord func
// instead of BasicRecord{} as the properties won't be initialized.
type BasicRecord struct {
	Keys  []string
	Vals  []interface{}
	index map[string]uint
}

// NewBasicRecord creates a *BasicRecord empty
// with the properties initialized.
func NewBasicRecord() *BasicRecord {
	return &BasicRecord{
		Keys:  []string{},
		Vals:  []interface{}{},
		index: map[string]uint{},
	}
}

// SetKeyOrder sets the order of the keys and returns a new struct.
// This is useful if you need to maintain column order
// for things like CSV output etc... This can also be useful
// to extract only the needed columns into a new record.
func (r BasicRecord) SetKeyOrder(keys ...string) OrderedRecord {
	newR := NewBasicRecord()
	newR.Keys = keys
	newR.Vals = make([]interface{}, len(keys))
	for i, k := range keys {
		newR.index[k] = uint(i)
		if v, ok := r.Get(k); ok {
			newR.Vals[i] = v
		}
	}

	return newR
}

// FromMSI sets the keys and vals of the record to the data in the map[string]interface{}
func (r *BasicRecord) FromMSI(msi map[string]interface{}) *BasicRecord {
	for k, v := range msi {
		r.Set(k, v)
	}
	return r
}

// Set will add or update the key value pair
func (r *BasicRecord) Set(key string, val interface{}) {
	if _, exists := r.index[key]; exists {
		r.Vals[r.index[key]] = val
		return
	}

	r.Keys = append(r.Keys, key)
	r.Vals = append(r.Vals, val)
	r.index[key] = uint(len(r.Vals) - 1)
	return
}

// Get will return the value for the specified key or nil.
// The returned bool indicates if the key existed or not.
func (r BasicRecord) Get(key string) (interface{}, bool) {
	if _, exists := r.index[key]; exists {
		return r.Vals[r.index[key]], true
	}
	return nil, false
}

// GetKeys implements the Record interface
func (r BasicRecord) GetKeys() []string {
	return r.Keys
}

// GetVals implements the Record interface
func (r BasicRecord) GetVals() []interface{} {
	return r.Vals
}

// String implements the fmt.Stringer interface
func (r BasicRecord) String() string {
	jsn, _ := json.Marshal(RecordToMSI(&r))
	return string(jsn)
}

//
// IDRecord implementation
//

// BasicIDRecord is a basic implementation of the IDRecord interface
type BasicIDRecord struct {
	Record
	IDKeys []string
}

// GetIDKeys gets the identifying keys from the IDRecord
func (idr BasicIDRecord) GetIDKeys() []string {
	return idr.IDKeys
}

// GetIDVals gets the identifying values from the IDRecord
// and returns nil if any of the values weren't found.
func (idr BasicIDRecord) GetIDVals() []interface{} {
	id := []interface{}{}
	for _, k := range idr.GetIDKeys() {
		if v, ok := idr.Get(k); ok {
			id = append(id, v)
		} else {
			return nil // couldn't find identifying key so bail
		}
	}
	return id
}
