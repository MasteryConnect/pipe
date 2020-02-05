package message

// Record defines what it takes to be a record message
type Record interface {
	Get(string) (interface{}, bool)
	GetKeys() []string
	GetVals() []interface{}
}

// MutableRecord defines what it takes to mutate a record
type MutableRecord interface {
	Record
	Set(string, interface{})
}

// OrderedRecord defines a record that specifies the column order.
// This is useful for things like CSVs or bulk SQL inserts where the order matters.
type OrderedRecord interface {
	Record
	SetKeyOrder(...string) OrderedRecord
}

// MutableOrderedRecord is the same as OrderedRecord but wih a MutableRecord
type MutableOrderedRecord interface {
	MutableRecord
	SetKeyOrder(...string) OrderedRecord
}

// NewRecord creates a *Record empty
func NewRecord() MutableRecord {
	return NewBasicRecord()
}

// NewIDRecord created a record that implements IDRecord and as such
// is identifiable // with some of value or combination of values from the record.
func NewIDRecord(IDKeys ...string) MutableIDRecord {
	return &BasicIDRecord{
		MutableRecord: NewRecord(),
		IDKeys:        IDKeys,
	}
}

// NewRecordFromMSI creates a *Record from a map[string]interface{}
// This is a convenience function to combine new and FromMSI
func NewRecordFromMSI(msi map[string]interface{}) MutableOrderedRecord {
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

// MutableIDRecord is the same as IDRecord but with a MutableRecord
// underlying it. It can be casted to an IDRecord.
type MutableIDRecord interface {
	MutableRecord
	GetIDKeys() []string
	GetIDVals() []interface{}
}

// GetNonIDKeys gets all the keys for the record that aren't a part of the ID.
// This can be useful when building SQL queries for a record.
func GetNonIDKeys(r IDRecord) []string {
	// get the non-ID keys
	idkeys := r.GetIDKeys()
	var nonIDKeys []string
	for _, key := range r.GetKeys() {
		exclude := false
		for _, idkey := range idkeys {
			if key == idkey {
				exclude = true
				break // skip this col as it exists in the id
			}
		}
		if !exclude {
			nonIDKeys = append(nonIDKeys, key)
		}
	}
	return nonIDKeys
}
