// Package lazy provides a generic utility for generating the next value in a sequence in a lazily manner
package lazy

import "time"

// Any is a type for the empty interface and all types will satisfy
type Any interface{}

// EvalFunc is a function type used to generate the values in a sequence.
type EvalFunc func(Any) (Any, Any)

// Options provide a way for passing options to use for the Generator
type Options struct {
	Generator EvalFunc
	Limit     uint
	Timeout   time.Duration
	InitState Any
}

type sequenceGenerator struct {
	quit      chan struct{}
	generator EvalFunc
	limit     uint
	timeout   time.Duration
	init      Any
}

// NewLazyGenerator returns a generator object
func NewLazyGenerator(gen EvalFunc, opt ...Options) *sequenceGenerator {
	var (
		timeout      = time.Duration(time.Second * 60 * 60 * 24)
		limit   uint = 1000000
		init    Any
	)
	if len(opt) != 0 {
		if opt[0].Timeout == 0 {
			timeout = time.Second * 60 * 60 * 24
		} else {
			timeout = opt[0].Timeout
		}
		if opt[0].Limit == 0 {
			limit = 10000000
		} else {
			limit = opt[0].Limit
		}
	}
	return &sequenceGenerator{
		quit:      make(chan struct{}),
		generator: gen,
		limit:     limit,
		timeout:   timeout,
		init:      init,
	}
}

// SetTimeout sets the duration of the generator. When not set, default is one year
func (sq *sequenceGenerator) SetTimeout(d time.Duration) *sequenceGenerator {
	sq.timeout = d
	return sq
}

// SetInit sets the initial values for the generator
func (sq *sequenceGenerator) SetInit(initState Any) *sequenceGenerator {
	sq.init = initState
	return sq
}

// SetLimit sets the number of times the generator will generate values. When not set, default is 10000000
func (sq *sequenceGenerator) SetLimit(n uint) *sequenceGenerator {
	sq.limit = n
	return sq
}

// Stop sends a signal to generator to free resources
func (sq *sequenceGenerator) Stop() {
	sq.quit <- struct{}{}
}

// Limit returns the number of times the generator will run.
func (sq *sequenceGenerator) Limit() int {
	return int(sq.limit)
}

func (sq *sequenceGenerator) cleanup(ch chan Any) {
	close(ch)
	close(sq.quit)
}

// Generate returns a generator function that generates the next value in the sequence
func (sq *sequenceGenerator) Generate() func() Any {
	return sq.buildLazyGenerator()
}

func (sq *sequenceGenerator) buildLazyGenerator() func() Any {
	retValChan := make(chan Any)
	go func() {
		var actState = sq.init
		var retVal Any
		timeout := time.After(sq.timeout)
		for i := 0; i < int(sq.limit); i++ {
			retVal, actState = sq.generator(actState)
			select {
			case <-sq.quit:
				sq.cleanup(retValChan)
				return
			case retValChan <- retVal:
			case <-timeout:
				sq.cleanup(retValChan)
				return
			}
		}
		sq.cleanup(retValChan)
	}()
	return func() Any {
		return <-retValChan
	}
}
