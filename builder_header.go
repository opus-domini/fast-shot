package fastshot

import (
	"github.com/opus-domini/fast-shot/constant/header"
	"github.com/opus-domini/fast-shot/constant/mime"
)

// HeaderBuilder provides fluent HTTP header configuration.
//
// The generic type parameter P is the parent builder returned by every
// method to keep the fluent chain going, so the same builder serves both
// ClientBuilder and RequestBuilder:
//
//	fastshot.NewClient("https://api.example.com").
//		Header().Add(header.UserAgent, "MyApp/1.0").
//		Build()
//
//	client.GET("/users").
//		Header().AddContentType(mime.JSON).
//		Send()
type HeaderBuilder[P any] struct {
	parent P
	header HeaderWrapper
}

// Header returns a HeaderBuilder for setting custom HTTP headers on the request.
func (b *RequestBuilder) Header() *HeaderBuilder[*RequestBuilder] {
	return &HeaderBuilder[*RequestBuilder]{
		parent: b,
		header: b.request.config.Header(),
	}
}

// Header returns a HeaderBuilder for setting custom HTTP headers on the client.
func (b *ClientBuilder) Header() *HeaderBuilder[*ClientBuilder] {
	return &HeaderBuilder[*ClientBuilder]{
		parent: b,
		header: b.client.Header(),
	}
}

// Add adds a custom header. If the header already exists, it will be appended.
func (b *HeaderBuilder[P]) Add(key header.Type, value string) P {
	b.header.Add(key, value)
	return b.parent
}

// AddAll adds custom headers. If a header already exists, it will be appended.
func (b *HeaderBuilder[P]) AddAll(headers map[header.Type]string) P {
	for key, value := range headers {
		b.header.Add(key, value)
	}
	return b.parent
}

// Set sets a custom header. If the header already exists, it will be overwritten.
func (b *HeaderBuilder[P]) Set(key header.Type, value string) P {
	b.header.Set(key, value)
	return b.parent
}

// SetAll sets custom headers. If a header already exists, it will be overwritten.
func (b *HeaderBuilder[P]) SetAll(headers map[header.Type]string) P {
	for key, value := range headers {
		b.header.Set(key, value)
	}
	return b.parent
}

// AddAccept sets the Accept header. If the header already exists, it will be appended.
func (b *HeaderBuilder[P]) AddAccept(value mime.Type) P {
	return b.Add(header.Accept, value.String())
}

// AddContentType sets the Content-Type header. If the header already exists, it will be appended.
func (b *HeaderBuilder[P]) AddContentType(value mime.Type) P {
	return b.Add(header.ContentType, value.String())
}

// AddUserAgent sets the User-Agent header. If the header already exists, it will be appended.
func (b *HeaderBuilder[P]) AddUserAgent(value string) P {
	return b.Add(header.UserAgent, value)
}
