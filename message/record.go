package message

import (
	"encoding/json"
)

// Record defines what it takes to be a record message
type Record interface {
	Get(string) (interface{}, bool)
	Set(string, interface{})
	GetKeys() []string
	GetVals() []interface{}
}

// OrderedRecord defines a record that specifies the column order.
// This is useful for things like CSVs or bulk SQL inserts where the order matters.
type OrderedRecord interface {
	Record
	SetOrder([]string) OrderedRecord
}

// BasicRecord is a collection of key/value pairs where the keys are strings
// It is similar to a map[string]interface{} accept it keeps column order.
type BasicRecord struct {
	Keys  []string
	Vals  []interface{}
	index map[string]uint
}

// NewRecord creates a *Record empty
func NewRecord() OrderedRecord {
	return &BasicRecord{
		Keys:  []string{},
		Vals:  []interface{}{},
		index: map[string]uint{},
	}
}

// NewRecordFromMSI creates a *Record from a map[string]interface{}
// This is a convenience function to combine new and FromMSI
func NewRecordFromMSI(msi map[string]interface{}) OrderedRecord {
	r := &BasicRecord{
		Keys:  []string{},
		Vals:  []interface{}{},
		index: map[string]uint{},
	}
	r.FromMSI(msi)
	return r
}

// RecordToMSI converts any record to a map[string]interface{}
// (note: column order is not preserved in a map[string]interface{})
func RecordToMSI(r Record) map[string]interface{} {
	m := map[string]interface{}{}
	for _, key := range r.GetKeys() {
		if val, ok := r.Get(key); ok {
			m[key] = val
		}
	}
	return m
}

// RecordToStrings converts a record to a slice of the values as strings.
// This is useful for things like CSVs.
func RecordToStrings(r Record) []string {
	vals := r.GetVals()
	row := make([]string, len(vals))
	for i, val := range vals {
		row[i] = String(val) // convert each array element to a string
	}
	return row
}

// SetOrder sets the order of the keys. This is useful if you need to maintain column order
// for things like CSV output etc...
func (r BasicRecord) SetOrder(keys []string) OrderedRecord {
	newR := BasicRecord{}
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

	return &newR
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
// IDRecord
//

// IDRecord is an identifyable record.
// That means the ID for the record (composite or not) is defined on the record.
type IDRecord struct {
	Record
	IDKeys []string
}

// GetIDKeys gets the identifying keys from the IDRecord
func (idr IDRecord) GetIDKeys() []string {
	return idr.IDKeys
}

// GetIDVals gets the identifying values from the IDRecord
func (idr IDRecord) GetIDVals() []interface{} {
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
