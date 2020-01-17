package yaml

import (
	"gopkg.in/yaml.v2"
)

// Message wraps a message with the YAML of the message.
type Message struct {
	M interface{}
	B []byte
}

// String implements the fmt.Stringer interface.
func (y *Message) String() string {
	return string(y.B)
}

// In returns the message wrapped in the yaml message.
func (y *Message) In() interface{} {
	return y.M
}

// To converts the message to a yaml message
func To(msg interface{}) (interface{}, error) {
	b, err := yaml.Marshal(msg)
	return &Message{M: msg, B: b}, err
}
