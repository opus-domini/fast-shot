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
	parent P
	add    func(*http.Cookie)
}

// Cookie returns a CookieBuilder for setting custom HTTP cookies on the request.
func (b *RequestBuilder) Cookie() *CookieBuilder[*RequestBuilder] {
	return &CookieBuilder[*RequestBuilder]{
		parent: b,
		add:    b.request.config.addCookie,
	}
}

// Cookie returns a CookieBuilder for setting custom HTTP cookies on the client.
func (b *ClientBuilder) Cookie() *CookieBuilder[*ClientBuilder] {
	return &CookieBuilder[*ClientBuilder]{
		parent: b,
		add:    b.client.addCookie,
	}
}

// Add adds a custom cookie.
func (b *CookieBuilder[P]) Add(cookie *http.Cookie) P {
	b.add(cookie)
	return b.parent
}
