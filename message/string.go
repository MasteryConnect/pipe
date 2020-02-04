package message

import "fmt"

// Stringser the interface for Strings()
// which is like String() but as a slice of strings.
// This can be useful for things like CSVs
type Stringser interface {
	Strings() []string
}

// String will extract a string from a message
// one way or another.
func String(m interface{}) string {
	if m == nil {
		return ""
	}
	switch v := m.(type) {
	case string:
		return v
	case []byte:
		return string(v)
	case fmt.Stringer:
		return v.String()
	case Inner:
		return String(v.In())
	default:
		return fmt.Sprintf("%v", m)
	}
}

// Strings will extract a []string from a message
// one way or another.
func Strings(m interface{}) []string {
	switch v := m.(type) {
	case string:
		return []string{v}
	case []byte:
		return []string{string(v)}
	case Stringser:
		return v.Strings()
	case fmt.Stringer:
		return []string{v.String()}
	case Inner:
		return Strings(v.In())
	default:
		return []string{fmt.Sprintf("%v", m)}
	}
}
