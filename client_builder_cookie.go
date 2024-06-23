package fastshot

import (
	"net/http"
)

// BuilderCookie is the interface that wraps the basic methods for setting HTTP CookiesWrapper.
var _ BuilderCookie[ClientBuilder] = (*ClientCookieBuilder)(nil)

// ClientCookieBuilder allows for setting custom HTTP CookiesWrapper.
type ClientCookieBuilder struct {
	parentBuilder *ClientBuilder
}

// Cookie returns a new ClientCookieBuilder for setting custom HTTP CookiesWrapper.
func (b *ClientBuilder) Cookie() *ClientCookieBuilder {
	return &ClientCookieBuilder{parentBuilder: b}
}

// Add adds a custom cookie to the HTTP client.
func (b *ClientCookieBuilder) Add(cookie *http.Cookie) *ClientBuilder {
	b.parentBuilder.client.Cookies().Add(cookie)
	return b.parentBuilder
}
