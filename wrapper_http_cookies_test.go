package fastshot

import (
	"net/http"
	"testing"
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
			if got := c.Get(tt.index); got != nil {
				t.Errorf("got %v, want nil", got)
			}
		})
	}
}

func TestDefaultHttpCookies_Get_Valid(t *testing.T) {
	c := newDefaultHttpCookies()
	cookie := &http.Cookie{Name: "test", Value: "value"}
	c.Add(cookie)

	if got := c.Get(0); got != cookie {
		t.Errorf("got %v, want %v", got, cookie)
	}
}
