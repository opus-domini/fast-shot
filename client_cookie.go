package fastshot

import (
	"net/http"
)

// ClientCookieBuilder allows for setting custom HTTP Cookies.
type ClientCookieBuilder struct {
	parentBuilder *ClientBuilder
}

// Cookie returns a new ClientCookieBuilder for setting custom HTTP Cookies.
func (b *ClientBuilder) Cookie() *ClientCookieBuilder {
	return &ClientCookieBuilder{parentBuilder: b}
}

// Add adds a custom cookie to the HTTP client.
func (b *ClientCookieBuilder) Add(cookie *http.Cookie) *ClientBuilder {
	b.parentBuilder.client.httpCookies = append(b.parentBuilder.client.httpCookies, cookie)
	return b.parentBuilder
}
