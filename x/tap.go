package x

// Tap will send message on through the pipe/line as well as to another channel.
func Tap(otherOut chan<- interface{}) func(<-chan interface{}, chan<- interface{}, chan<- error) {
	return func(in <-chan interface{}, out chan<- interface{}, errs chan<- error) {
		for m := range in {
			otherOut <- m
			out <- m
		}
	}
}
