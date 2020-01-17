package message

import "strings"

// Batch is a message type that can contain a list of other messages.
type Batch []interface{}

// String implements fmt.Stringer.
func (b Batch) String() string {
	strs := []string{}
	for _, v := range b {
		strs = append(strs, String(v))
	}
	return strings.Join(strs, "\n")
}

// Size returns the length of the batch.
func (b Batch) Size() int {
	return len(b)
}
