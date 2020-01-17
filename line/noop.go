package line

// Noop is the transform noop or passthrough
func Noop(in <-chan interface{}, out chan<- interface{}, errs chan<- error) {
	for m := range in {
		out <- m
	}
}

// NoopC is the noop consumer
func NoopC(in <-chan interface{}, errs chan<- error) {
	for range in {
	}
}
