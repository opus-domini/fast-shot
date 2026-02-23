package fastshot

import (
	"errors"
	"testing"
)

func TestDefaultValidations_Get_OutOfBounds(t *testing.T) {
	tests := []struct {
		name  string
		index int
	}{
		{"negative index", -1},
		{"zero on empty", 0},
		{"large index", 100},
	}

	v := newDefaultValidations(nil)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := v.Get(tt.index); got != nil {
				t.Errorf("got %v, want nil", got)
			}
		})
	}
}

func TestDefaultValidations_Get_Valid(t *testing.T) {
	err := errors.New("test error")
	v := newDefaultValidations([]error{err})

	if got := v.Get(0); got != err {
		t.Errorf("got %v, want %v", got, err)
	}
}
