package line

// InlineTfunc is the function signature for a transformer func.
// The idea is to make writing an anonymous func right inline
// with the definition of the pipeline easier. It also provides
// a way to boil the essense of most transformers down to a pure
// single input single output func. That makes testing a lot easier
// as well.
type InlineTfunc func(interface{}) (interface{}, error)

// Inline wraps an InlineTfunc and returns a Tfunc.
// Most of the time, transformers will just range over
// the in channel and do stuff inside of the range and then send any
// errors off to the error channel. This func does that
// for you so you can just write a simpler transformer func.
// The parameter is the incoming message.
// The resulting interface{} is the outgoing message to be
// sent downstream. If nil is passed, no message will be sent
// downstream. If and error is returned, it will be sent
// down the errror channel.
func Inline(it InlineTfunc) Tfunc {
	return func(in <-chan interface{}, out chan<- interface{}, errs chan<- error) {
		for msg := range in {
			newMsg, err := it(msg)
			if err != nil {
				errs <- err
			}
			if newMsg != nil {
				out <- newMsg
			}
		}
	}
}

// I is a convenience wrapper around Inline
func I(it InlineTfunc) Tfunc {
	return Inline(it)
}

// ForEach is a convenience wrapper around Inline
// and actually makes more sense as the func name
func ForEach(it InlineTfunc) Tfunc {
	return Inline(it)
}
