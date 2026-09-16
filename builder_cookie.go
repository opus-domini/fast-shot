package fastshot

import (
	"net/http"
)

// CookieBuilder provides fluent HTTP cookie configuration.
//
// The generic type parameter P is the parent builder returned by every
// method to keep the fluent chain going, so the same builder serves both
// ClientBuilder and RequestBuilder:
//
//	fastshot.NewClient("https://api.example.com").
//		Cookie().Add(&http.Cookie{Name: "session", Value: "abc123"}).
//		Build()
//
//	client.GET("/protected").
//		Cookie().Add(&http.Cookie{Name: "csrf_token", Value: "xyz789"}).
//		Send()
type CookieBuilder[P any] struct {
	parent  P
	cookies CookiesWrapper
}

// Cookie returns a CookieBuilder for setting custom HTTP cookies on the request.
func (b *RequestBuilder) Cookie() *CookieBuilder[*RequestBuilder] {
	return &CookieBuilder[*RequestBuilder]{
		parent:  b,
		cookies: b.request.config.Cookies(),
	}
}

// Cookie returns a CookieBuilder for setting custom HTTP cookies on the client.
func (b *ClientBuilder) Cookie() *CookieBuilder[*ClientBuilder] {
	return &CookieBuilder[*ClientBuilder]{
		parent:  b,
		cookies: b.client.Cookies(),
	}
}

// Add adds a custom cookie.
func (b *CookieBuilder[P]) Add(cookie *http.Cookie) P {
	b.cookies.Add(cookie)
	return b.parent
}
