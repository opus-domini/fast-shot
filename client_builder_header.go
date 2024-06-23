package fastshot

import (
	"github.com/opus-domini/fast-shot/constant/header"
)

// BuilderHeader is the interface that wraps the basic methods for setting custom HTTP BuilderHeader.
var _ BuilderHeader[ClientBuilder] = (*ClientHeaderBuilder)(nil)

// ClientHeaderBuilder allows for setting custom HTTP BuilderHeader.
type ClientHeaderBuilder struct {
	parentBuilder *ClientBuilder
}

// Header returns a new ClientHeaderBuilder for setting custom HTTP BuilderHeader.
func (b *ClientBuilder) Header() *ClientHeaderBuilder {
	return &ClientHeaderBuilder{parentBuilder: b}
}

// Add adds a custom header to the HTTP client. If header already exists, it will be appended.
func (b *ClientHeaderBuilder) Add(key, value string) *ClientBuilder {
	b.parentBuilder.client.Header().Add(key, value)
	return b.parentBuilder
}

// AddAll adds multiple custom headers to the HTTP client. If header already exists, it will be appended.
func (b *ClientHeaderBuilder) AddAll(headers map[string]string) *ClientBuilder {
	for key, value := range headers {
		b.Add(key, value)
	}
	return b.parentBuilder
}

// Set sets a custom header to the HTTP client. If header already exists, it will be overwritten.
func (b *ClientHeaderBuilder) Set(key, value string) *ClientBuilder {
	b.parentBuilder.client.Header().Set(key, value)
	return b.parentBuilder
}

// SetAll sets multiple custom headers to the HTTP client. If header already exists, it will be overwritten.
func (b *ClientHeaderBuilder) SetAll(headers map[string]string) *ClientBuilder {
	for key, value := range headers {
		b.Set(key, value)
	}
	return b.parentBuilder
}

// AddAccept sets the Accept header. If header already exists, it will be appended.
func (b *ClientHeaderBuilder) AddAccept(value string) *ClientBuilder {
	b.Add(header.Accept, value)
	return b.parentBuilder
}

// AddContentType sets the Content-Type header. If header already exists, it will be appended.
func (b *ClientHeaderBuilder) AddContentType(value string) *ClientBuilder {
	b.Add(header.ContentType, value)
	return b.parentBuilder
}

// AddUserAgent sets the User-Agent header. If header already exists, it will be appended.
func (b *ClientHeaderBuilder) AddUserAgent(value string) *ClientBuilder {
	b.Add(header.UserAgent, value)
	return b.parentBuilder
}
