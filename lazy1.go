// Package lazy lazily generates values in a given sequence
package lazy

import "time"

// Any is a type that will later be unexported
type Any interface{}

type nextVal chan Any

// EvalFunc is another one
type EvalFunc func(Any) (Any, Any)

type SequenceGenerator struct {
	quit      chan struct{}
	generator EvalFunc
	limit     uint64
	timeout   time.Duration
}

// NewLazySeq returns a generator
func NewLazySeq(gen EvalFunc) *SequenceGenerator {
	return &SequenceGenerator{
		quit:      make(chan struct{}),
		generator: gen,
		limit:     10000000,
		timeout:   time.Second * 60 * 60 * 24,
	}
}

// Generate returns a func that when called returns the next value in a sequence
func (sq *SequenceGenerator) Generate(initState Any) <-chan Any {
	val := sq.buildLazyGenerator(initState)
	return val
}

func (sq *SequenceGenerator) buildLazyGenerator(initState Any) chan Any {
	retValChan := make(chan Any)
	go func() {
		var actState = initState
		var retVal Any
		for i := 0; i < int(sq.limit); i++ {
			retVal, actState = sq.generator(actState)
			select {
			case <-sq.quit:
				sq.cleanup(retValChan)
				return
			case retValChan <- retVal:
			}
		}
		sq.cleanup(retValChan)
	}()
	return retValChan
}

// Stop sends a signal to free resources
func (sq *SequenceGenerator) Stop() {
	sq.quit <- struct{}{}
}

func (sq *SequenceGenerator) cleanup(ch chan Any) {
	close(ch)
	close(sq.quit)
}

// Limit generate the values of the sequence for a specified number of times
func (sq *SequenceGenerator) Limit(n uint64) *SequenceGenerator {
	sq.limit = n
	return sq
}

// Timeout will stop the genration of values after specified amount of time
func (sq *SequenceGenerator) Timeout(d time.Duration) *SequenceGenerator {
	sq.timeout = d
	return sq
}

// BuildLazyEvaluator is a higher-order function that lazily evaluates next value in a sequence
func buildLazyEvaluator(eval EvalFunc, initState Any) func() Any {
	retValChan := make(chan Any)
	quitChan := make(chan struct{})
	retValFunc := func() Any {
		return <-retValChan
	}
	go func() {
		var actState = initState
		var retVal Any
		for {
			retVal, actState = eval(actState)
			select {
			case <-quitChan:
				close(quitChan)
				close(retValChan)
				return
			case retValChan <- retVal:
			}
		}
	}()
	return retValFunc
}

// Ints is a wrapper to BuildLazyEvaluator that lazily generate integer values in sequence
func Ints(eval EvalFunc, initState Any) func() int {
	ef := buildLazyEvaluator(eval, initState)
	return func() int {
		return ef().(int)
	}
}
