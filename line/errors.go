package line

import "fmt"

var (
	// ErrNoErrsWaitGroup represents when the user has customized the errs channel but hasn't provided a waitgroup
	ErrNoErrsWaitGroup = fmt.Errorf("No sync.WaitGroup passed for errs channel draining")
)
