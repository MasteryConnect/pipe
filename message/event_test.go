package message_test

import (
	"fmt"
	"time"

	"github.com/masteryconnect/pipe/message"
)

func ExampleEvent() {
	fmt.Println(message.Event{
		Timestamp: time.Time{},
		Source:    "/var/log/foo.log",
		Message:   "DEBUG: some formatted log even message",
	})

	// Output: 0001-01-01T00:00:00Z /var/log/foo.log DEBUG: some formatted log even message
}
