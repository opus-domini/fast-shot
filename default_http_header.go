package fastshot

import (
	"net/http"
)

// DefaultHttpHeader implements HeaderWrapper interface and provides a default HTTP header.
var _ HeaderWrapper = (*DefaultHttpHeader)(nil)

// DefaultHttpHeader implements HeaderWrapper interface and provides a default HTTP header.
type DefaultHttpHeader struct {
	header *http.Header
}

// Unwrap will return the underlying header
func (c *DefaultHttpHeader) Unwrap() *http.Header {
	return c.header
}

// Get will return the value of the specified key
func (c *DefaultHttpHeader) Get(key string) string {
	return c.header.Get(key)
}

// Add will append a new key value pair to the underlying header
func (c *DefaultHttpHeader) Add(key, value string) {
	c.header.Add(key, value)
}

// Set will set the value of the specified key
func (c *DefaultHttpHeader) Set(key, value string) {
	c.header.Set(key, value)
}
