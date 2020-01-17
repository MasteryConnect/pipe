package message

// Bytes wraps a message with the bytes version of the message
// and keeps the original message around.
type Bytes struct {
	M interface{}
	B []byte
}

// String implements the fmt.Stringer interface.
func (j *Bytes) String() string {
	return string(j.B)
}

// In returns the message wrapped in the json message.
func (j *Bytes) In() interface{} {
	return j.M
}
