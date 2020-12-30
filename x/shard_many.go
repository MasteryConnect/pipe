package x

import (
	"errors"
	"fmt"
	"hash/crc32"
	"sync"

	l "github.com/MasteryConnect/pipe/line"
)

// Shard messages by a key e.g. two messages with the same key
// will go down the same go channel, and hence be processed in order they are
// received.
// This is similar to the Many processor except that the Many processor will
// process messages in any order i.e. two messages with the same key could
// process out of order.

type ShardMany struct {
	concurrency  int
	tfunc        l.Tfunc
	keyFunc      ShardManyKeyFunc
	shardInChans []chan interface{}
	wg           *sync.WaitGroup
}

// Given a message return a key for the message. Multiple messages can return
// the same key if they should be processed in order they come down the in
// channel
type ShardManyKeyFunc func(msg interface{}) (key []byte)

// Create a new ShardMany instance
// concurrency - the number of go routines (shards)
// tfunc - the pipeline/function to process each message
// shardManyKeyFunc - when passed a message this function with return a key
func NewShardMany(concurrency int, tfunc l.Tfunc, shardManyKeyFunc ShardManyKeyFunc) (*ShardMany, error) {
	if concurrency < 1 {
		return nil, fmt.Errorf("Concurrency should be greater than 1, actual value is: %d", concurrency)
	}
	if tfunc == nil {
		return nil, errors.New("tfunc cannot be null")
	}
	if shardManyKeyFunc == nil {
		return nil, errors.New("shardManyKeyFunc cannot be null")
	}
	return &ShardMany{
		concurrency: concurrency,
		tfunc:       tfunc,
		keyFunc:     shardManyKeyFunc,
		wg:          &sync.WaitGroup{},
	}, nil
}

// The transformer function
// Each message pulled from the in channel will be:
// 1) Passed to the ShardManyKeyFunc to get the messages key
// 2) Sent to the shard processor (channel) for messages with the key from
//    step 1
// To determine which shard to pass the message to, the key is turned into an
// int, and then the int of the key is mapped to a shard by using the modulus
// operator with the number of shards (concurrency).
func (sm *ShardMany) T(in <-chan interface{}, out chan<- interface{}, errs chan<- error) {
	sm.initShards(out, errs)
	shardCount := uint32(sm.concurrency)
	for msg := range in {
		key := sm.keyFunc(msg)
		shardIndex := crc32.ChecksumIEEE(key) % shardCount
		shardChan := sm.shardInChans[shardIndex]
		shardChan <- msg
	}

	// Close each FanoutTfunc's in channel, signaling to them that no more
	// messages are coming, so they can exit
	sm.shutdown()
	// Wait for each FanoutTfunc to finish processing and exist
	sm.wg.Wait()
}

// The in chan was closed, so stop the shards by closing their channels
func (sm *ShardMany) shutdown() {
	for _, inChan := range sm.shardInChans {
		close(inChan)
	}
}

// Create a in channel per tfunc,  and goroutine the tfunc with the created
// channel, so we can pass along messages to each tfunc to process
// concurrently
func (sm *ShardMany) initShards(out chan<- interface{}, errs chan<- error) {
	// We will wait for all FanoutTfunc's to finish
	sm.wg.Add(sm.concurrency)
	for i := 0; i < sm.concurrency; i++ {
		inChan := make(chan interface{})
		sm.shardInChans = append(sm.shardInChans, inChan)

		go sm.wgT(inChan, out, errs)
	}
}

// A Tfunc wrapper to handle the WaitGroup.Done() call to tell ShardMany when
// a ShardMany is finished
func (sm *ShardMany) wgT(in <-chan interface{}, out chan<- interface{}, errs chan<- error) {
	defer sm.wg.Done()
	sm.tfunc(in, out, errs)
}
