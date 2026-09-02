package fastshot

import (
	"net/http"
)

// Compile-time check that DefaultHttpCookies implements CookiesWrapper.
var _ CookiesWrapper = (*DefaultHttpCookies)(nil)

// DefaultHttpCookies implements CookiesWrapper interface and provides a default HTTP cookies.
type DefaultHttpCookies struct {
	cookies []*http.Cookie
}

// Unwrap will return the underlying cookies
func (c *DefaultHttpCookies) Unwrap() []*http.Cookie {
	return c.cookies
}

// Count will return the number of cookies
func (c *DefaultHttpCookies) Count() int {
	return len(c.cookies)
}

// Get will return the cookie at the specified index
func (c *DefaultHttpCookies) Get(index int) *http.Cookie {
	if index < 0 || index >= len(c.cookies) {
		return nil
	}
	return c.cookies[index]
}

// Add will append a new cookie to the underlying cookies
func (c *DefaultHttpCookies) Add(cookie *http.Cookie) {
	c.cookies = append(c.cookies, cookie)
}

// newDefaultHttpCookies initializes a new DefaultHttpCookies.
func newDefaultHttpCookies() *DefaultHttpCookies {
	return &DefaultHttpCookies{
		cookies: []*http.Cookie{},
	}
}
