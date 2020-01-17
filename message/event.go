package message

import (
	"fmt"
	"time"
)

// Event is a message type that
type Event struct {
	Timestamp time.Time   // the time that the event occured
	Source    interface{} // the source or fhte event (log file or device etc...)
	Message   interface{} // the actual event message itself (usually a string)
}

// String implements the fmt.Stringer interface
func (e Event) String() string {
	return fmt.Sprintf("%v %v %v", e.Timestamp.Format(time.RFC3339), e.Source, e.Message)
}

// Strings implements the Strinsger interface for Event
func (e Event) Strings() []string {
	return []string{
		e.Timestamp.Format(time.RFC3339),
		String(e.Source),
		String(e.Message),
	}
}
