package line

// Consumer is the default consumer for the line.
func Consumer(in <-chan interface{}, errs chan<- error) {
	for msg := range in {
		if v, ok := msg.(Acker); ok {
			v.Ack()
		}
	}
}
