package fastshot

import (
	"github.com/opus-domini/fast-shot/constant/header"
)

// RequestHeaderBuilder is a builder for setting custom HTTP Header.
type RequestHeaderBuilder struct {
	parentBuilder *RequestBuilder
}

// Header returns a new RequestHeaderBuilder for setting custom HTTP Header.
func (b *RequestBuilder) Header() *RequestHeaderBuilder {
	return &RequestHeaderBuilder{parentBuilder: b}
}

// Add adds a custom header to the HTTP request. If header already exists, it will be appended.
func (b *RequestHeaderBuilder) Add(key, value string) *RequestBuilder {
	b.parentBuilder.request.httpHeader.Add(key, value)
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
	b.parentBuilder.request.httpHeader.Set(key, value)
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
