package message

import "context"

// ContextGetter defines what it takes to get the context from a message.
type ContextGetter interface {
	GetContext() context.Context
}
