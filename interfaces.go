package fastshot

import (
	"context"
	"io"
	"net/http"
	"net/url"
	"time"
)

type Client interface {
	ClientConfig
	ClientHttpMethods
}

type ClientConfig interface {
	HttpClient() *http.Client
	HttpHeader() *http.Header
	SetHttpCookie(cookie *http.Cookie)
	HttpCookies() []*http.Cookie
	SetValidation(validation error)
	Validations() []error
	ConfigBaseURL
}

type ConfigBaseURL interface {
	BaseURL() *url.URL
}

type ClientHttpMethods interface {
	GET(path string) *RequestBuilder
	POST(path string) *RequestBuilder
	PUT(path string) *RequestBuilder
	DELETE(path string) *RequestBuilder
	PATCH(path string) *RequestBuilder
	HEAD(path string) *RequestBuilder
	CONNECT(path string) *RequestBuilder
	OPTIONS(path string) *RequestBuilder
	TRACE(path string) *RequestBuilder
}

// BuilderHeader is the interface that wraps the basic methods for setting custom HTTP BuilderHeader.
type BuilderHeader[T any] interface {
	Add(key, value string) *T
	AddAll(headers map[string]string) *T
	Set(key, value string) *T
	SetAll(headers map[string]string) *T
	AddAccept(value string) *T
	AddContentType(value string) *T
	AddUserAgent(value string) *T
}

// BuilderCookie is the interface that wraps the basic methods for setting HTTP Cookies.
type BuilderCookie[T any] interface {
	Add(cookie *http.Cookie) *T
}

// BuilderAuth is the interface that wraps the basic methods for setting HTTP Authentication.
type BuilderAuth[T any] interface {
	Set(value string) *T
	BearerToken(token string) *T
	BasicAuth(username, password string) *T
}

// BuilderHttpClientConfig is the interface that wraps the basic methods for setting HTTP ClientConfig configuration.
type BuilderHttpClientConfig[T any] interface {
	SetCustomTransport(transport http.RoundTripper) *T
	SetTimeout(duration time.Duration) *T
	SetFollowRedirects(follow bool) *T
	SetProxy(proxyURL string) *T
}

// BuilderRequestContext is the interface that wraps the basic methods for setting HTTP BuilderRequestContext.
type BuilderRequestContext[T any] interface {
	Set(ctx context.Context) *T
}

// BuilderRequestBody is the interface that wraps the basic methods for setting HTTP BuilderRequestBody.
type BuilderRequestBody[T any] interface {
	AsReader(body io.Reader) *T
	AsString(body string) *T
	AsJSON(obj interface{}) *T
}

// BuilderRequestQuery is the interface that wraps the basic methods for setting HTTP BuilderRequestQuery.
type BuilderRequestQuery[T any] interface {
	AddParam(param, value string) *T
	AddParams(params map[string]string) *T
	SetParam(param, value string) *T
	SetParams(params map[string]string) *T
	SetRawString(query string) *T
}

// BuilderRequestRetry is the interface that wraps the basic methods for setting HTTP BuilderRequestRetry.
type BuilderRequestRetry[T any] interface {
	Set(retries int, retryInterval time.Duration) *T
}
