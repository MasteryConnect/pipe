package jwx

import (
	l "github.com/masteryconnect/pipe/line"
)

func ExampleSign() {
	l.New().SetP(func(out chan<- interface{}, errs chan<- error) {
		out <- `{"sub": "1234567890", "name": "John Doe", "iat": 1516239022}` // sign a JWT like payload
		out <- `{"foo": "bar"}`
	}).Add(
		l.I(Sign([]byte("secret"))),
		l.Stdout,
	).Run()
	// output:
	// eyJhbGciOiJIUzI1NiJ9.eyJzdWIiOiAiMTIzNDU2Nzg5MCIsICJuYW1lIjogIkpvaG4gRG9lIiwgImlhdCI6IDE1MTYyMzkwMjJ9.TH3DqT4ee58ttaEScoXZQDzCntSvyMaV2L_DLyqjlos
	// eyJhbGciOiJIUzI1NiJ9.eyJmb28iOiAiYmFyIn0.UeFzIir7vMJQY60bZ-ru2COIheSBQE0ano37obDR1tQ
}
