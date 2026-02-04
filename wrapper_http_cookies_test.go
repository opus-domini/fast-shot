package fastshot

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDefaultHttpCookies_Get_OutOfBounds(t *testing.T) {
	tests := []struct {
		name  string
		index int
	}{
		{"negative index", -1},
		{"zero on empty", 0},
		{"large index", 100},
	}

	c := newDefaultHttpCookies()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Nil(t, c.Get(tt.index))
		})
	}
}

func TestDefaultHttpCookies_Get_Valid(t *testing.T) {
	c := newDefaultHttpCookies()
	cookie := &http.Cookie{Name: "test", Value: "value"}
	c.Add(cookie)

	assert.Equal(t, cookie, c.Get(0))
}
