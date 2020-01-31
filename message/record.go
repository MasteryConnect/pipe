package message

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
	SetKeyOrder(...string) OrderedRecord
}

// NewRecord creates a *Record empty
func NewRecord() OrderedRecord {
	return NewBasicRecord()
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

//
// IDRecord
//

// IDRecord is an identifyable record.
// That means the ID for the record (composite or not) is defined on the record.
type IDRecord interface {
	Record
	GetIDKeys() []string
	GetIDVals() []interface{}
}
