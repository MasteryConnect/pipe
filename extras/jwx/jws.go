package jwx

import (
	l "github.com/masteryconnect/pipe/line"
	"github.com/masteryconnect/pipe/message"

	"github.com/lestrrat-go/jwx/jwa"
	"github.com/lestrrat-go/jwx/jws"
)

// Sign creates an InlineTfunc to sign the string of incoming messages
func Sign(key interface{}) l.InlineTfunc {
	return func(m interface{}) (interface{}, error) {
		payload := message.String(m)
		return jws.Sign([]byte(payload), jwa.HS256, key)
	}
}
