package json

import (
	"encoding/json"

	"github.com/masteryconnect/pipe/message"
)

// To converts the message to a json message
func To(msg interface{}) (interface{}, error) {
	var err error
	var b []byte

	switch v := msg.(type) {
	case message.Record:
		b, err = json.Marshal(message.RecordToMSI(v))
	default:
		b, err = json.Marshal(msg)
	}

	return &message.Bytes{M: msg, B: b}, err
}
