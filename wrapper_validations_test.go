package fastshot

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
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
			assert.Nil(t, v.Get(tt.index))
		})
	}
}

func TestDefaultValidations_Get_Valid(t *testing.T) {
	err := errors.New("test error")
	v := newDefaultValidations([]error{err})

	assert.Equal(t, err, v.Get(0))
}
