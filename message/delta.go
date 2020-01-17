package message

// DeltaType is the type of mutation for the sql to take
type DeltaType uint8

const (
	// Insert is for an INSERT SQL statement
	Insert DeltaType = iota + 1
	// Update is for an UPDATE SQL statement
	Update
	// Delete is for a DELETE SQL statement
	Delete
)

// Delta is a record that is a change record.
type Delta struct {
	*Record

	Type    DeltaType              `json:"deltatype"`
	Changes map[string]interface{} `json:"changes,omitempty"`
	Table   string                 `json:"deltatable,omitempty"`
}
