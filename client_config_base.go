package fastshot

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"

	"github.com/opus-domini/fast-shot/constant"
	"github.com/opus-domini/fast-shot/constant/method"
)

// ClientConfigBase serves as the main entry point for configuring HTTP clients.
type ClientConfigBase struct {
	httpClient    HttpClientComponent
	httpHeader    HeaderWrapper
	httpCookies   CookiesWrapper
	validations   ValidationsWrapper
	beforeRequest []func(*http.Request) error
	afterResponse []func(*http.Request, *http.Response)
	ConfigBaseURL
}

// HttpClient for ClientConfigBase returns the HTTP client.
func (c *ClientConfigBase) HttpClient() HttpClientComponent {
	return c.httpClient
}

// SetHttpClient for ClientConfigBase sets the HttpClientComponent.
func (c *ClientConfigBase) SetHttpClient(httpClient HttpClientComponent) {
	c.httpClient = httpClient
}

// Header for ClientConfigBase returns the HeaderWrapper.
func (c *ClientConfigBase) Header() HeaderWrapper {
	return c.httpHeader
}

// Cookies for ClientConfigBase returns the CookiesWrapper.
func (c *ClientConfigBase) Cookies() CookiesWrapper {
	return c.httpCookies
}

// Validations for ClientConfigBase returns the ValidationsWrapper.
func (c *ClientConfigBase) Validations() ValidationsWrapper {
	return c.validations
}

// BeforeRequestHooks returns the before-request hooks.
func (c *ClientConfigBase) BeforeRequestHooks() []func(*http.Request) error {
	return c.beforeRequest
}

// AfterResponseHooks returns the after-response hooks.
func (c *ClientConfigBase) AfterResponseHooks() []func(*http.Request, *http.Response) {
	return c.afterResponse
}

// AddBeforeRequestHook appends a before-request hook.
func (c *ClientConfigBase) AddBeforeRequestHook(hook func(*http.Request) error) {
	c.beforeRequest = append(c.beforeRequest, hook)
}

// AddAfterResponseHook appends an after-response hook.
func (c *ClientConfigBase) AddAfterResponseHook(hook func(*http.Request, *http.Response)) {
	//nolint:bodyclose // False positive: appending a hook function, not handling a response body.
	c.afterResponse = append(c.afterResponse, hook)
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

// newClientConfigBase initializes a new ClientConfigBase with a given baseURL.
func newClientConfigBase(baseURL string) *ClientConfigBase {
	var validations []error

	if baseURL == "" {
		validations = append(validations, errors.New(constant.ErrMsgEmptyBaseURL))
	}

	parsedURL, err := url.Parse(baseURL)
	if err != nil {
		validations = append(validations, errors.Join(errors.New(constant.ErrMsgParseURL), err))
	}

	return &ClientConfigBase{
		httpClient:    newDefaultHttpClient(),
		httpHeader:    newDefaultHttpHeader(),
		httpCookies:   newDefaultHttpCookies(),
		validations:   newDefaultValidations(validations),
		ConfigBaseURL: newDefaultBaseURL(parsedURL),
	}
}

// newBalancedClientConfigBase initializes a new ClientConfigBase with a given baseURLs.
func newBalancedClientConfigBase(baseURLs []string) *ClientConfigBase {
	var validations []error

	var parsedURLs []*url.URL
	for index, baseURL := range baseURLs {
		if baseURL == "" {
			validations = append(validations, fmt.Errorf("base URL %d: %s", index, constant.ErrMsgEmptyBaseURL))
			continue
		}

		parsedURL, err := url.Parse(baseURL)
		if err != nil {
			validations = append(validations, errors.Join(errors.New(constant.ErrMsgParseURL), err))
		}
		parsedURLs = append(parsedURLs, parsedURL)
	}

	if len(parsedURLs) == 0 {
		validations = append(validations, errors.New(constant.ErrMsgEmptyBaseURL))
	}

	return &ClientConfigBase{
		httpClient:    newDefaultHttpClient(),
		httpHeader:    newDefaultHttpHeader(),
		httpCookies:   newDefaultHttpCookies(),
		validations:   newDefaultValidations(validations),
		ConfigBaseURL: newBalancedBaseURL(parsedURLs),
	}
}
