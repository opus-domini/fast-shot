package fastshot

import (
	"context"
	"io"
	"net/http"
	"net/url"
	"time"
)

// HeaderWrapper is the interface that wraps the basic methods for setting HTTP Headers.
type HeaderWrapper interface {
	Unwrap() *http.Header
	Get(key string) string
	Add(key, value string)
	Set(key, value string)
}

// CookiesWrapper is the interface that wraps the basic methods for setting HTTP CookiesWrapper.
type CookiesWrapper interface {
	Unwrap() []*http.Cookie
	Get(index int) *http.Cookie
	Count() int
	Add(cookie *http.Cookie)
}

// ValidationsWrapper is the interface that wraps the basic methods for setting HTTP ValidationsWrapper.
type ValidationsWrapper interface {
	Unwrap() []error
	Get(index int) error
	IsEmpty() bool
	Count() int
	Add(err error)
}

// Client is the interface that wraps the basic methods for setting HTTP Client.
type Client interface {
	ClientConfig
	ClientHttpMethods
}

// ClientConfig is the interface that wraps the basic methods for setting HTTP ClientConfig.
type ClientConfig interface {
	ConfigHttpClient
	Header() HeaderWrapper
	Cookies() CookiesWrapper
	Validations() ValidationsWrapper
	ConfigBaseURL
}

// ConfigHttpClient is the interface that wraps the basic methods for setting HTTP HttpClient.
type ConfigHttpClient interface {
	SetHttpClient(httpClient HttpClientComponent)
	HttpClient() HttpClientComponent
}

// HttpClientComponent is the interface that wraps the basic methods for setting HTTP HttpClientComponent.
type HttpClientComponent interface {
	Do(req *http.Request) (*http.Response, error)
	Transport() http.RoundTripper
	SetTransport(http.RoundTripper)
	Timeout() time.Duration
	SetTimeout(time.Duration)
	SetFollowRedirects(follow bool)
}

// ConfigBaseURL is the interface that wraps the basic methods for setting HTTP BaseURL.
type ConfigBaseURL interface {
	BaseURL() *url.URL
}

// ClientHttpMethods is the interface that wraps the basic methods for setting HTTP ClientHttpMethods.
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

// BuilderCookie is the interface that wraps the basic methods for setting HTTP CookiesWrapper.
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
	SetCustomHttpClient(httpClient HttpClientComponent) *T
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
	SetConstantBackoff(interval time.Duration, maxAttempts uint) *T
	SetConstantBackoffWithJitter(interval time.Duration, maxAttempts uint) *T
	SetExponentialBackoff(interval time.Duration, maxAttempts uint, backoffRate float64) *T
	SetExponentialBackoffWithJitter(interval time.Duration, maxAttempts uint, backoffRate float64) *T
	WithRetryCondition(shouldRetry func(response Response) bool) *T
	WithMaxDelay(duration time.Duration) *T
}
