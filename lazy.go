// Package lazy lazily generates values in a given sequence
package lazy

type Any interface{}
type EvalFunc func(Any) (Any, Any)

// BuildLazyEvaluator is a higher-order function that lazily evaluates next value in a sequence
func buildLazyEvaluator(eval EvalFunc, initState Any) func(...int) Any {
	retValChan := make(chan Any)
	quitChan := make(chan struct{})
	retValFunc := func(sig ...int) Any {
		if len(sig) > 0 && sig[0] < 0 {
			defer func() { quitChan <- struct{}{} }()
		}
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
func Ints(eval EvalFunc, initState Any) func(...int) int {
	ef := buildLazyEvaluator(eval, initState)
	return func(sig ...int) int {
		val := ef(sig...)
		if val == nil {
			return 0
		}
		return val.(int)
	}
}
