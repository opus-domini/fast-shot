package fastshot

import (
	"context"
	"io"
	"net/http"
	"time"
)

// Header is the interface that wraps the basic methods for setting custom HTTP Header.
type Header[T any] interface {
	Add(key, value string) *T
	AddAll(headers map[string]string) *T
	Set(key, value string) *T
	SetAll(headers map[string]string) *T
	AddAccept(value string) *T
	AddContentType(value string) *T
	AddUserAgent(value string) *T
}

// Cookie is the interface that wraps the basic methods for setting HTTP Cookies.
type Cookie[T any] interface {
	Add(cookie *http.Cookie) *T
}

// Auth is the interface that wraps the basic methods for setting HTTP Authentication.
type Auth[T any] interface {
	Set(value string) *T
	BearerToken(token string) *T
	BasicAuth(username, password string) *T
}

// ClientConfig is the interface that wraps the basic methods for setting HTTP Client configuration.
type ClientConfig[T any] interface {
	SetTimeout(duration time.Duration) *T
	SetFollowRedirects(follow bool) *T
	SetCustomTransport(transport http.RoundTripper) *T
}

// RequestContext is the interface that wraps the basic methods for setting HTTP RequestContext.
type RequestContext[T any] interface {
	Set(ctx context.Context) *T
}

// RequestBody is the interface that wraps the basic methods for setting HTTP RequestBody.
type RequestBody[T any] interface {
	AsReader(body io.Reader) *T
	AsString(body string) *T
	AsJSON(obj interface{}) *T
}

// RequestQuery is the interface that wraps the basic methods for setting HTTP RequestQuery.
type RequestQuery[T any] interface {
	AddParam(param, value string) *T
	AddParams(params map[string]string) *T
	SetParam(param, value string) *T
	SetParams(params map[string]string) *T
	SetRawString(query string) *T
}

// RequestRetry is the interface that wraps the basic methods for setting HTTP RequestRetry.
type RequestRetry[T any] interface {
	Set(retries int, retryInterval time.Duration) *T
}
