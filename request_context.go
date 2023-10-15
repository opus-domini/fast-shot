package fastshot

import (
	"context"
)

// RequestContextBuilder serves as the main entry point for configuring Request Context.
type RequestContextBuilder struct {
	parentBuilder *RequestBuilder
}

// Context returns a new RequestContextBuilder for setting custom Context.
func (b *RequestBuilder) Context() *RequestContextBuilder {
	return &RequestContextBuilder{parentBuilder: b}
}

// Set sets the Context.
func (b *RequestContextBuilder) Set(ctx context.Context) *RequestBuilder {
	b.parentBuilder.request.ctx = ctx
	return b.parentBuilder
}
