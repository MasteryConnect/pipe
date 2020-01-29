package http

import (
	"net/http"

	l "github.com/masteryconnect/pipe/line"
)

func ExampleDo_T() {
	l.New().SetP(func(out chan<- interface{}, errs chan<- error) {
		req, err := http.NewRequest("GET", "https://google.com", nil)
		if err != nil {
			errs <- err
		} else {
			out <- req
			out <- req
		}
	}).Add(
		Do{}.T,
		l.I(func(m interface{}) (interface{}, error) {
			return m.(*http.Response).Status, nil
		}),
		l.Stdout,
	).Run()

	// Output:
	// 200 OK
	// 200 OK
}
