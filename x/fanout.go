package x

import "sync"

/*
Fanout takes one or more Tfunc's, and when a message is received on the 'in'
channel, sends each Tfunc the message. Each Tfunc gets its own 'in' channel and
reuses the 'out' and 'err' channel passed to Fanout.
*/

const (
	ANY_TYPE = "_any_type_"
)

type Fanout struct {
	tfuncs       []FanoutTfunc
	msgTypesFunc FanoutMsgTypesFunc
	chanLookup   map[string][]chan interface{} // Type to channel list lookup
	allInChans   []chan interface{}
	wg           *sync.WaitGroup
}

// The Tfunc's Fanout fans out messages to.
type FanoutTfunc interface {
	T(<-chan interface{}, chan<- interface{}, chan<- error)
}

// Pass this in when creating Fanout if you want to include certain
// messages sent to a FanoutTfunc. When a message is passed to this function
// it should return the types (names) of the msg. These along with a
// FanoutTfunc's implementation of FanoutIncludeTypes
// will be used to determine the FanoutTfunc's the message
// is sent to.
type FanoutMsgTypesFunc func(msg interface{}) (types []string)

// A FanoutTfunc can implement this interface to tell Fanout that it only
// wants the message types returned by I() sent to it.
type FanoutIncludeTypes interface {
	I() []string
}

// Create a new Fanout with FanoutMsgTypesFunc to get the types (names) of a
// message with the intention of filtering messages to the provided
// FanoutTfunc's. The filtering only works if a FanoutTfunc implements
// FanoutIncludeTypes. If no FanoutMsgTypesFunc is passed in, then
// FanoutTfunc's get all messages
func NewFanout(tfuncs []FanoutTfunc, msgTypes FanoutMsgTypesFunc) *Fanout {
	return &Fanout{
		tfuncs:       tfuncs,
		msgTypesFunc: msgTypes,
		chanLookup:   make(map[string][]chan interface{}),
		wg:           &sync.WaitGroup{}, // Wait for all tfuncs to finish
	}
}

func (f *Fanout) T(in <-chan interface{}, out chan<- interface{}, errs chan<- error) {
	f.initTfuncs(out, errs)

	// Receive in messages and pass along (fanout) to the tfunc's channels
	for msg := range in {
		// Can't filter without knowing the message type, so send to all by default
		if f.msgTypesFunc == nil {
			for _, inChan := range f.allInChans {
				inChan <- msg
			}
		} else {
			anyTypeInChans, hasAnyTypeInChans := f.chanLookup[ANY_TYPE]

			// Fan out to Tfunc's wanting all messages
			for _, anyInChan := range anyTypeInChans {
				anyInChan <- msg
			}

			// If there are Tfunc's wanting filtered messages, check if any want
			// this message
			if (hasAnyTypeInChans && len(f.chanLookup) > 1) || (!hasAnyTypeInChans && len(f.chanLookup) > 0) {
				// Used for lookup of tfunc in channels already sent to
				seenInChans := make(map[chan interface{}]struct{})
				// For each type the message is, find channels interested in those types
				// and send the message to them. Check that the message hasn't already
				// been sent to a channel by checking seenInChans
				for _, msgType := range f.msgTypesFunc(msg) {
					for _, inChan := range f.chanLookup[msgType] {
						if _, seen := seenInChans[inChan]; !seen {
							// simple key lookup with zero byte value
							seenInChans[inChan] = struct{}{}
							inChan <- msg
						}
					}
				}
			}
		}
	}

	// Close each FanoutTfunc's in channel, signaling to them that no more
	// messages are coming, so they can exit
	f.shutdown()
	// Wait for each FanoutTfunc to finish processing and exist
	f.wg.Wait()
}

// Create a in channel per tfunc,  and goroutine the tfunc with the created
// channel, so we can pass along messages to each tfunc to process
// concurrently
func (f *Fanout) initTfuncs(out chan<- interface{}, errs chan<- error) {
	// We will wait for all FanoutTfunc's to finish
	f.wg.Add(len(f.tfuncs))
	// Initialize each FanoutTfunc and goroutine it
	for _, tfunc := range f.tfuncs {
		inChan := make(chan interface{})
		f.allInChans = append(f.allInChans, inChan)
		// Check to see if the tfunc wants only certain message types sent to it
		if fit, ok := tfunc.(FanoutIncludeTypes); ok {
			for _, t := range fit.I() {
				f.initArrayChan(t)
				f.chanLookup[t] = append(f.chanLookup[t], inChan)
			}
		} else {
			f.initArrayChan(ANY_TYPE)
			f.chanLookup[ANY_TYPE] = append(f.chanLookup[ANY_TYPE], inChan)
		}
		go f.wgT(tfunc, inChan, out, errs)
	}
}

// A FanoutTfunc wrapper to handle the WaitGroup.Done() call to tell Fanout
// when a FanoutTfunc is finished
func (f *Fanout) wgT(tf FanoutTfunc, in <-chan interface{}, out chan<- interface{}, errs chan<- error) {
	defer f.wg.Done()
	tf.T(in, out, errs)
}

// The in chan was closed, so stop the tfunc's by closing their channels
func (f *Fanout) shutdown() {
	for _, inChan := range f.allInChans {
		close(inChan)
	}
}

func (f *Fanout) initArrayChan(name string) {
	if _, ok := f.chanLookup[name]; !ok {
		var typeChan []chan interface{}
		f.chanLookup[name] = typeChan
	}
}
