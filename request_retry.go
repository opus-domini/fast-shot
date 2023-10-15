package fastshot

import (
	"time"
)

// RequestRetryBuilder serves as the main entry point for configuring Request retries.
type RequestRetryBuilder struct {
	parentBuilder *RequestBuilder
}

// Retry returns a new RequestRetryBuilder for setting custom HTTP Cookies.
func (b *RequestBuilder) Retry() *RequestRetryBuilder {
	return &RequestRetryBuilder{parentBuilder: b}
}

// Set sets the number of retries and the retry interval.
func (b *RequestRetryBuilder) Set(retries int, retryInterval time.Duration) *RequestBuilder {
	b.parentBuilder.request.retries = retries
	b.parentBuilder.request.retryInterval = retryInterval
	return b.parentBuilder
}
