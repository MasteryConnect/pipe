package json

// StreamChunk is a struct that has the block found in the stream
// as well as the position in the stream
type StreamChunk struct {
	Src    interface{}
	Offset int
	Bytes  []byte
}

// String makes this struct a fmt.Stringer
func (s *StreamChunk) String() string {
	return string(s.Bytes)
}
