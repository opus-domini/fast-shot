package fastshot

import (
	"net/http"

	"github.com/opus-domini/fast-shot/constant/header"
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
func (c *DefaultHttpHeader) Get(key header.Type) string {
	return c.header.Get(key.String())
}

// Add will append a new key value pair to the underlying header
func (c *DefaultHttpHeader) Add(key header.Type, value string) {
	c.header.Add(key.String(), value)
}

// Set will set the value of the specified key
func (c *DefaultHttpHeader) Set(key header.Type, value string) {
	c.header.Set(key.String(), value)
}

// newDefaultHttpHeader initializes a new DefaultHttpHeader with a given header.
func newDefaultHttpHeader() *DefaultHttpHeader {
	return &DefaultHttpHeader{
		header: &http.Header{},
	}
}
