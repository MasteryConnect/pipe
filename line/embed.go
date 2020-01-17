package line

// Embed runs the whole pipeline as a transformer of a parent pipeline.
func (l *Line) Embed(parentIn <-chan interface{}, parentOut chan<- interface{}, parentErrs chan<- error) {
	embedP := func(out chan<- interface{}, errs chan<- error) {
		for msg := range parentIn {
			out <- msg
		}
	}
	embedC := func(in <-chan interface{}, errs chan<- error) {
		for msg := range in {
			parentOut <- msg
		}
	}

	l.p = embedP
	l.c = embedC

	l.SetErrs(parentErrs)

	err := l.Run()
	if err != nil {
		parentErrs <- err
	}
}
