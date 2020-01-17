package message

// Metadata is a map to hold metadata about a message.
type Metadata map[string]interface{}

// Get gets a value from the metadata.
func (s Metadata) Get(key string) (interface{}, bool) {
	v, exists := s[key]
	return v, exists
}

// Set sets a value in the metadata.
func (s Metadata) Set(key string, val interface{}) {
	s[key] = val
}

// Message is a baseline implementation of a Messenger.
// This is currently used to migration from the old pipeline project.
type Message struct {
	Metadata // any structured data needed for the message
	Body     []byte
}

// String implements fmt.Stringer.
func (bm Message) String() string {
	return string(bm.Body)
}
