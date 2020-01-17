package x

// Cap will cap off a pipeline with no-op. This can be useful when embedding
// a pipeline and don't want to forward message back to the parent pipeline.
func Cap(in <-chan interface{}, out chan<- interface{}, errs chan<- error) {
	for range in {
		// noop
	}
}
