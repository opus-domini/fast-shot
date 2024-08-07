package fastshot

import (
	"context"
)

// BuilderRequestContext is the interface that wraps the basic methods for setting custom HTTP Context.
var _ BuilderRequestContext[RequestBuilder] = (*RequestContextBuilder)(nil)

// RequestContextBuilder serves as the main entry point for configuring Request Context.
type RequestContextBuilder struct {
	parentBuilder *RequestBuilder
	requestConfig *RequestConfigBase
}

// Context returns a new RequestContextBuilder for setting custom Context.
func (b *RequestBuilder) Context() *RequestContextBuilder {
	return &RequestContextBuilder{
		parentBuilder: b,
		requestConfig: b.request.config,
	}
}

// Set sets the Context.
func (b *RequestContextBuilder) Set(ctx context.Context) *RequestBuilder {
	b.requestConfig.Context().Set(ctx)
	return b.parentBuilder
}
