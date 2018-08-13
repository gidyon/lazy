// Package lazy lazily generates values in a given sequence

package lazy

import (
	"reflect"
	"testing"
)

func Test_buildLazyEvaluator(t *testing.T) {
	type args struct {
		eval      EvalFunc
		initState Any
	}
	tests := []struct {
		name string
		args args
		want func(...int) Any
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := buildLazyEvaluator(tt.args.eval, tt.args.initState); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("buildLazyEvaluator() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestInts(t *testing.T) {
	type args struct {
		eval      EvalFunc
		initState Any
	}
	tests := []struct {
		name string
		args args
		want func(...int) int
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Ints(tt.args.eval, tt.args.initState); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Ints() = %v, want %v", got, tt.want)
			}
		})
	}
}
