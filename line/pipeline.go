package line

import "context"

// Pfunc is the function signature for a producer func.
type Pfunc func(chan<- interface{}, chan<- error)

// PfuncContext is the function signature for a producer func.
type PfuncContext func(context.Context, chan<- interface{}, chan<- error)

// Tfunc is the function signature for a transformer func.
type Tfunc func(<-chan interface{}, chan<- interface{}, chan<- error)

// Cfunc is the function signature for a Consumer jfunc.
type Cfunc func(<-chan interface{}, chan<- error)

// Pipeline defines what it takes to be a pipeline.
// This means you could write your own implementation
// of a pipeline (say a distributed one) and still be able
// to use all of the producers, consumers, and transformers
// that match these interfaces.
type Pipeline interface {
	SetP(Pfunc) Pipeline
	SetPContext(PfuncContext) Pipeline
	Add(...Tfunc) Pipeline
	SetC(Cfunc) Pipeline
	SetErrs(chan<- error) Pipeline
	Run() error
	RunContext(context.Context) error
	Embed(<-chan interface{}, chan<- interface{}, chan<- error) // act as a Tfunc
}

// Acker is something that can be "Ack"ed.
type Acker interface {
	Ack()
}
