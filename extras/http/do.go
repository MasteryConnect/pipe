package http

import (
	"net/http"

	"github.com/MasteryConnect/pipe/message"
)

// Do will execute incoming http.Request message.
type Do struct {
	Client *http.Client
}

// T implements the Tfunc interface
func (d Do) T(in <-chan interface{}, out chan<- interface{}, errs chan<- error) {
	client := d.Client
	if client == nil {
		client = http.DefaultClient
	}

	for m := range in {
		cl := client

		// get the client if set
		switch v := m.(type) {
		case message.Clienter:
			if v.Client() != nil {
				cl = v.Client()
			}
		}

		// get the request to send
		switch v := m.(type) {
		case *http.Request:
			resp, err := cl.Do(v)
			if err != nil {
				errs <- err
			} else {
				out <- resp
			}
		case message.Requester:
			resp, err := cl.Do(v.Request())
			if err != nil {
				errs <- err
			} else {
				out <- resp
			}
		}
	}
}
