package fastshot

import (
	"net/http"
)

// BuilderCookie is the interface that wraps the basic methods for setting HTTP CookiesWrapper.
var _ BuilderCookie[RequestBuilder] = (*RequestCookieBuilder)(nil)

// RequestCookieBuilder serves as the main entry point for configuring HTTP CookiesWrapper.
type RequestCookieBuilder struct {
	parentBuilder *RequestBuilder
	requestConfig *RequestConfigBase
}

// Cookie returns a new RequestCookieBuilder for setting custom HTTP CookiesWrapper.
func (b *RequestBuilder) Cookie() *RequestCookieBuilder {
	return &RequestCookieBuilder{
		parentBuilder: b,
		requestConfig: b.request.config,
	}
}

// Add adds a custom cookie to the HTTP client.
func (b *RequestCookieBuilder) Add(cookie *http.Cookie) *RequestBuilder {
	b.requestConfig.Cookies().Add(cookie)
	return b.parentBuilder
}
