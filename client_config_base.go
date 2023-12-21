package fastshot

import (
	"github.com/opus-domini/fast-shot/constant/method"
	"net/http"
	"net/url"
	"sync/atomic"
)

type (
	// ClientConfigBase serves as the main entry point for configuring HTTP clients.
	ClientConfigBase struct {
		httpClient  *http.Client
		httpHeader  *http.Header
		httpCookies []*http.Cookie
		validations []error
		ConfigBaseURL
	}

	// DefaultBaseURL implements ConfigBaseURL interface and provides a single base URL.
	DefaultBaseURL struct {
		baseURL *url.URL
	}

	// BalancedBaseURL implements ConfigBaseURL interface and provides load balancing.
	BalancedBaseURL struct {
		baseURLs       []*url.URL
		currentBaseURL uint32
	}
)

// BaseURL for DefaultBaseURL returns the base URL.
func (c *DefaultBaseURL) BaseURL() *url.URL {
	return c.baseURL
}

// BaseURL for BalancedBaseURL returns the next base URL in the list.
func (c *BalancedBaseURL) BaseURL() *url.URL {
	currentIndex := atomic.LoadUint32(&c.currentBaseURL)
	atomic.AddUint32(&c.currentBaseURL, 1)
	c.currentBaseURL = c.currentBaseURL % uint32(len(c.baseURLs))
	return c.baseURLs[currentIndex]
}

// HttpClient for ClientConfigBase returns the HTTP client.
func (c *ClientConfigBase) HttpClient() *http.Client {
	return c.httpClient
}

// HttpHeader for ClientConfigBase returns the HTTP header.
func (c *ClientConfigBase) HttpHeader() *http.Header {
	return c.httpHeader
}

// SetHttpCookie for ClientConfigBase sets the HTTP cookie.
func (c *ClientConfigBase) SetHttpCookie(cookie *http.Cookie) {
	c.httpCookies = append(c.httpCookies, cookie)
}

// HttpCookies for ClientConfigBase returns the HTTP cookies.
func (c *ClientConfigBase) HttpCookies() []*http.Cookie {
	return c.httpCookies
}

// SetValidation for ClientConfigBase sets the validation.
func (c *ClientConfigBase) SetValidation(validations error) {
	c.validations = append(c.validations, validations)
}

// Validations for ClientConfigBase returns the validations.
func (c *ClientConfigBase) Validations() []error {
	return c.validations
}

// GET is a shortcut for NewRequest(c, method.GET, path).
func (c *ClientConfigBase) GET(path string) *RequestBuilder {
	return newRequest(c, method.GET, path)
}

// POST is a shortcut for NewRequest(c, method.POST, path).
func (c *ClientConfigBase) POST(path string) *RequestBuilder {
	return newRequest(c, method.POST, path)
}

// PUT is a shortcut for NewRequest(c, method.PUT, path).
func (c *ClientConfigBase) PUT(path string) *RequestBuilder {
	return newRequest(c, method.PUT, path)
}

// DELETE is a shortcut for NewRequest(c, method.DELETE, path).
func (c *ClientConfigBase) DELETE(path string) *RequestBuilder {
	return newRequest(c, method.DELETE, path)
}

// PATCH is a shortcut for NewRequest(c, method.PATCH, path).
func (c *ClientConfigBase) PATCH(path string) *RequestBuilder {
	return newRequest(c, method.PATCH, path)
}

// HEAD is a shortcut for NewRequest(c, method.HEAD, path).
func (c *ClientConfigBase) HEAD(path string) *RequestBuilder {
	return newRequest(c, method.HEAD, path)
}

// CONNECT is a shortcut for NewRequest(c, method.CONNECT, path).
func (c *ClientConfigBase) CONNECT(path string) *RequestBuilder {
	return newRequest(c, method.CONNECT, path)
}

// OPTIONS is a shortcut for NewRequest(c, method.OPTIONS, path).
func (c *ClientConfigBase) OPTIONS(path string) *RequestBuilder {
	return newRequest(c, method.OPTIONS, path)
}

// TRACE is a shortcut for NewRequest(c, method.TRACE, path).
func (c *ClientConfigBase) TRACE(path string) *RequestBuilder {
	return newRequest(c, method.TRACE, path)
}
