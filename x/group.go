package x

import (
	"sync"
	"time"

	"github.com/masteryconnect/pipe/message"
)

// GroupByFunc is the function that is used to get the group names
// to group a message by.
type GroupByFunc func(msg interface{}) (groups []string)

// ReduceFunc defines the reducer you can optionaly specify.
// The purpose is to define how a group of messages is collapsed
// down into a single message.
type ReduceFunc func(memo, msg interface{}) (newmemo interface{})

type groupAddCloser interface {
	Add(m interface{})
	Close()
}

// GroupMsg is the message to send down stream.
type GroupMsg struct {
	message.Batch
	Name string
}

// In implements Inner for message.Get
func (gm *GroupMsg) In() interface{} {
	return gm.Batch
}

// Group is the grouping transformer for pipe/line.
type Group struct {
	By   GroupByFunc
	Size int

	Reduce ReduceFunc

	groups *sync.Map
	wg     sync.WaitGroup
}

// NewGroup makes a new Group.
func NewGroup(size int, by GroupByFunc) *Group {
	return &Group{Size: size, By: by, groups: &sync.Map{}}
}

// NewReduceGroup makes a new Group with a reducer
func NewReduceGroup(by GroupByFunc, reduce ReduceFunc) *Group {
	return &Group{By: by, Reduce: reduce, groups: &sync.Map{}}
}

// T is the transform function.
func (g *Group) T(in <-chan interface{}, out chan<- interface{}, errs chan<- error) {
	if g.groups == nil {
		g.groups = &sync.Map{}
	}

	for m := range in {
		for _, groupName := range g.By(m) {
			g.addToGroup(m, groupName, out, errs)
		}
	}

	// send anything that's left
	g.groups.Range(func(_, v interface{}) bool {
		v.(groupAddCloser).Close()
		return true
	})
	g.wg.Wait() // wait for all the Out channels for the groupBatches to finish
}

func (g *Group) addToGroup(m interface{}, groupName string, out chan<- interface{}, errs chan<- error) {
	// see if there is already a group
	if b, ok := g.groups.Load(groupName); ok {
		b.(groupAddCloser).Add(m)
		return
	}

	// group is new
	var newGroup groupAddCloser
	if g.Reduce == nil {
		newGroup = g.newGroupBatch(groupName, out, errs)
	} else {
		newGroup = g.newGroupReducer(out, errs)
	}

	g.groups.Store(groupName, newGroup)
	newGroup.Add(m)
}

//
// groupReducer
//

type groupReducer struct {
	reduce func(memo, val interface{}) (newmemo interface{})
	out    chan<- interface{}
	memo   interface{}
	wg     *sync.WaitGroup
}

func (g *Group) newGroupReducer(out chan<- interface{}, errs chan<- error) *groupReducer {
	g.wg.Add(1)
	return &groupReducer{out: out, reduce: g.Reduce, wg: &g.wg}
}

func (gr *groupReducer) Add(m interface{}) {
	if gr.memo == nil {
		gr.memo = m
		return
	}
	gr.memo = gr.reduce(gr.memo, m)
}

func (gr *groupReducer) Close() {
	gr.out <- gr.memo
	gr.wg.Done()
}

//
// groupBatch
//

// groupBatch is the wrapper around the batch
type groupBatch struct {
	Batch Batch
	Name  string
	In    chan interface{}
	Out   chan interface{}
	wg    *sync.WaitGroup
	timer *time.Timer
}

func (g *Group) newGroupBatch(name string, out chan<- interface{}, errs chan<- error) *groupBatch {
	g.wg.Add(1)
	gb := &groupBatch{
		Name:  name,
		In:    make(chan interface{}),
		Out:   make(chan interface{}),
		Batch: CloseableBatch(g.Size, 0, 0),
		wg:    &g.wg,
	}
	go func() {
		defer close(gb.Out)
		gb.Batch.T(gb.In, gb.Out, errs)
	}()
	go gb.Run(out)
	return gb
}

func (gb *groupBatch) Add(m interface{}) {
	gb.In <- m
}

func (gb *groupBatch) Run(out chan<- interface{}) {
	for b := range gb.Out {
		out <- &GroupMsg{Batch: b.(message.Batch), Name: gb.Name}
	}
	gb.wg.Done() // we are finally done with this groupBatch
}

func (gb *groupBatch) Close() {
	close(gb.In) // don't take any more in messages
}
