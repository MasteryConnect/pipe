package x

import (
	"sync"

	l "github.com/MasteryConnect/pipe/line"
)

// IfFunc is a func passed in to determin if the message should be used
type IfFunc func(interface{}) bool

// If only allows the first Limit number of messages through then stops the pipeline.
type If struct {
	Check  IfFunc
	Tfunc  l.Tfunc
	Else   l.Tfunc
	OnlyIf bool
}

// IF applies an IfFunc to see if the Tfunc should get the message
// if not, let the message pass through
func IF(t l.Tfunc, check IfFunc) l.Tfunc {
	return (&If{Check: check, Tfunc: t}).T
}

// OnlyIF applies an IfFunc to see if the Tfunc should get the message
// if not, ignore the message
func OnlyIF(t l.Tfunc, check IfFunc) l.Tfunc {
	return (&If{Check: check, Tfunc: t, OnlyIf: true}).T
}

// IFElse applies an IfFunc to see if the Tfunc should run of Else func
func IFElse(t, e l.Tfunc, check IfFunc) l.Tfunc {
	return (&If{Check: check, Tfunc: t, Else: e}).T
}

// T implements the pipeline interface.
func (i *If) T(inMsgs <-chan interface{}, outMsgs chan<- interface{}, errs chan<- error) {
	proxyIn := make(chan interface{})
	proxyOut := make(chan interface{})
	elseIn := make(chan interface{})
	elseOut := make(chan interface{})

	// run the positive Tfunc
	go func() {
		defer close(proxyOut)
		i.Tfunc(proxyIn, proxyOut, errs)
	}()

	// run the else Tfunc
	if i.Else != nil {
		go func() {
			defer close(elseOut)
			i.Else(elseIn, elseOut, errs)
		}()
	} else {
		close(elseOut)
	}

	go func() {
		defer close(proxyIn)
		defer close(elseIn)
		for msg := range inMsgs {
			if i.Check(msg) {
				proxyIn <- msg
			} else if i.Else != nil {
				elseIn <- msg
			} else if !i.OnlyIf {
				outMsgs <- msg
			}
		}
	}()

	var wg sync.WaitGroup
	wg.Add(1)

	// drain proxy output
	go func() {
		defer wg.Done()
		for msg := range proxyOut {
			outMsgs <- msg
		}
	}()

	// drain else output
	if i.Else != nil {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for msg := range elseOut {
				outMsgs <- msg
			}
		}()
	}

	// wait for everything to drain
	wg.Wait()
}
