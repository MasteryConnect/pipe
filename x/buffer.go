package x

// Buffer will create a buffer of Size to help "drain" a previous step.
type Buffer struct {
	N int
}

// T is the Tfunc for Buffer.
func (b Buffer) T(in <-chan interface{}, out chan<- interface{}, errs chan<- error) {
	buf := make(chan interface{}, b.N)

	go func() {
		defer close(buf)
		for msg := range in {
			buf <- msg
		}
	}()

	for msg := range buf {
		out <- msg
	}
}
