package message

import "net/http"

// Requester is a message interface that has or is an http.Request
type Requester interface {
	Request() *http.Request
}

// Clienter is a message interface that has or is an http.Client
type Clienter interface {
	Client() *http.Client
}
