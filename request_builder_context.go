package fastshot

import (
	"context"
)

// RequestContextBuilder serves as the main entry point for configuring the request context.
type RequestContextBuilder struct {
	parentBuilder *RequestBuilder
	requestConfig *RequestConfigBase
}

// Context returns a new RequestContextBuilder for setting the request Context.
func (b *RequestBuilder) Context() *RequestContextBuilder {
	return &RequestContextBuilder{
		parentBuilder: b,
		requestConfig: b.request.config,
	}
}

// Set sets the Context.
func (b *RequestContextBuilder) Set(ctx context.Context) *RequestBuilder {
	if ctx != nil {
		b.requestConfig.Context().Set(ctx)
	}
	return b.parentBuilder
}
