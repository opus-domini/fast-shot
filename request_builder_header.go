package fastshot

import (
	"github.com/opus-domini/fast-shot/constant/header"
)

// BuilderHeader is the interface that wraps the basic methods for setting custom HTTP BuilderHeader.
var _ BuilderHeader[RequestBuilder] = (*RequestHeaderBuilder)(nil)

// RequestHeaderBuilder is a builder for setting custom HTTP BuilderHeader.
type RequestHeaderBuilder struct {
	parentBuilder *RequestBuilder
	requestConfig *RequestConfigBase
}

// Header returns a new RequestHeaderBuilder for setting custom HTTP BuilderHeader.
func (b *RequestBuilder) Header() *RequestHeaderBuilder {
	return &RequestHeaderBuilder{
		parentBuilder: b,
		requestConfig: b.request.config,
	}
}

// Add adds a custom header to the HTTP request. If header already exists, it will be appended.
func (b *RequestHeaderBuilder) Add(key, value string) *RequestBuilder {
	b.requestConfig.httpHeader.Add(key, value)
	return b.parentBuilder
}

// AddAll adds custom headers to the HTTP request. If header already exists, it will be appended.
func (b *RequestHeaderBuilder) AddAll(headers map[string]string) *RequestBuilder {
	for key, value := range headers {
		b.Add(key, value)
	}
	return b.parentBuilder
}

// Set sets a custom header to the HTTP request. If header already exists, it will be overwritten.
func (b *RequestHeaderBuilder) Set(key, value string) *RequestBuilder {
	b.requestConfig.httpHeader.Set(key, value)
	return b.parentBuilder
}

// SetAll sets custom headers to the HTTP request. If header already exists, it will be overwritten.
func (b *RequestHeaderBuilder) SetAll(headers map[string]string) *RequestBuilder {
	for key, value := range headers {
		b.Set(key, value)
	}
	return b.parentBuilder
}

// AddAccept sets the Accept header. If header already exists, it will be appended.
func (b *RequestHeaderBuilder) AddAccept(value string) *RequestBuilder {
	b.Add(header.Accept, value)
	return b.parentBuilder
}

// AddContentType sets the Content-Type header. If header already exists, it will be appended.
func (b *RequestHeaderBuilder) AddContentType(value string) *RequestBuilder {
	b.Add(header.ContentType, value)
	return b.parentBuilder
}

// AddUserAgent sets the User-Agent header. If header already exists, it will be appended.
func (b *RequestHeaderBuilder) AddUserAgent(value string) *RequestBuilder {
	b.Add(header.UserAgent, value)
	return b.parentBuilder
}
