package fastshot

import (
	"net/http"
)

// BuilderCookie is the interface that wraps the basic methods for setting HTTP Cookies.
var _ BuilderCookie[ClientBuilder] = (*ClientCookieBuilder)(nil)

// ClientCookieBuilder allows for setting custom HTTP Cookies.
type ClientCookieBuilder struct {
	parentBuilder *ClientBuilder
}

// BuilderCookie returns a new ClientCookieBuilder for setting custom HTTP Cookies.
func (b *ClientBuilder) Cookie() *ClientCookieBuilder {
	return &ClientCookieBuilder{parentBuilder: b}
}

// Add adds a custom cookie to the HTTP client.
func (b *ClientCookieBuilder) Add(cookie *http.Cookie) *ClientBuilder {
	b.parentBuilder.client.SetHttpCookie(cookie)
	return b.parentBuilder
}
