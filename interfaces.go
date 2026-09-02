package fastshot

import (
	"context"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/opus-domini/fast-shot/constant/header"
)

// Client is the interface that wraps the basic methods for configuring and executing HTTP requests.
//
// It combines ClientConfig for setup and ClientHttpMethods for executing requests, providing
// a complete HTTP client solution. This interface is the main entry point for users of the library.
//
// Example usage:
//
//	client := fastshot.NewClient("https://api.example.com").Build()
//	response, err := client.GET("/users").Send()
//	if err != nil {
//		log.Fatal(err)
//	}
//	fmt.Println(response.Status().Code())
type Client interface {
	ClientConfig
	ClientHttpMethods
}

// ClientConfig is the interface that wraps the basic methods for configuring an HTTP client.
//
// Example usage:
//
//	client := fastshot.NewClient("https://api.example.com").
//		Header().Add(header.UserAgent, "MyApp/1.0").
//		Cookie().Add(&http.Cookie{Name: "session", Value: "abc123"}).
//		Config().SetTimeout(10 * time.Second).
//		Build()
type ClientConfig interface {
	ConfigHttpClient
	Header() HeaderWrapper
	Cookies() CookiesWrapper
	Validations() ValidationsWrapper
	ConfigBaseURL
	JSONCodec() JSONCodec
	SetJSONCodec(JSONCodec)
	BeforeRequestHooks() []func(*http.Request) error
	AfterResponseHooks() []func(*http.Request, *http.Response)
	AddBeforeRequestHook(func(*http.Request) error)
	AddAfterResponseHook(func(*http.Request, *http.Response))
}

// ConfigHttpClient is the interface that wraps the basic methods for configuring the underlying HTTP client.
//
// It is essential for providing fine-grained control over the HTTP client used for
// making requests. It allows users to set a custom HTTP client or retrieve the current one,
// enabling advanced use cases such as custom transport layers or connection pooling.
//
// Example usage:
//
//	customClient := &http.Client{
//		Timeout: 30 * time.Second,
//		Transport: &http.Transport{
//			MaxIdleConns: 100,
//			IdleConnTimeout: 90 * time.Second,
//		},
//	}
//
//	client := fastshot.NewClient("https://api.example.com").
//		Config().SetCustomHttpClient(myHttpClientComponent).
//		Build()
type ConfigHttpClient interface {
	SetHttpClient(httpClient HttpClientComponent)
	HttpClient() HttpClientComponent
}

// HttpClientComponent is the interface that wraps the basic methods for executing HTTP requests.
//
// It abstracts the actual HTTP client implementation, allowing for easy substitution
// of the underlying client (e.g., for testing or using a custom implementation).
type HttpClientComponent interface {
	Do(req *http.Request) (*http.Response, error)
	Transport() http.RoundTripper
	SetTransport(http.RoundTripper)
	Timeout() time.Duration
	SetTimeout(time.Duration)
	SetFollowRedirects(follow bool)
}

// ConfigBaseURL is the interface that wraps the basic method for retrieving the base URL.
//
// It supports both single base URL configurations and load-balanced configurations
// with multiple base URLs.
type ConfigBaseURL interface {
	BaseURL() *url.URL
}

// ClientHttpMethods is the interface that wraps the basic HTTP methods for making requests.
//
// Example usage:
//
//	client := fastshot.NewClient("https://api.example.com").Build()
//
//	// GET request
//	getResp, err := client.GET("/users").Send()
//
//	// POST request with JSON body
//	postResp, err := client.POST("/users").
//		Body().AsJSON(map[string]string{"name": "John Doe"}).
//		Send()
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

// BodyWrapper is the interface that wraps the basic methods for handling request and response bodies.
//
// Implementations are not safe for concurrent use, matching the semantics of
// io.ReadCloser and http.Response.Body.
//
// Example:
//
//	body := newBufferedBody(DefaultJSONCodec())
//	err := body.WriteAsJSON(map[string]string{"key": "value"})
//
//	var result map[string]any
//	err = body.ReadAsJSON(&result)
type BodyWrapper interface {
	io.ReadCloser
	ReadAsJSON(obj any) error
	WriteAsJSON(obj any) error
	ReadAsXML(obj any) error
	WriteAsXML(obj any) error
	ReadAsString() (string, error)
	WriteAsString(body string) error
	WriteAsFormData(fields map[string]string) (contentType string, err error)
	Set(body io.Reader) error
	Unwrap() io.Reader
}

// HeaderWrapper is the interface that wraps the basic methods for managing HTTP headers.
//
// This wrapper provides an abstraction layer over the standard http.Header type,
// allowing for type-safe header manipulation and potential future enhancements without
// changing the public API.
//
// It enables the library to implement custom header handling logic, such as
// case-insensitive header matching or header-specific validations, while maintaining
// a consistent interface for both internal use and potential extension points.
//
// Example (for library developers):
//
//	type CustomHeaderWrapper struct {
//		header http.Header
//	}
//
//	func (w *CustomHeaderWrapper) Set(key header.Type, value string) {
//		w.header.Set(string(key), value)
//		// Custom logic, e.g., logging or validation
//	}
type HeaderWrapper interface {
	Unwrap() *http.Header
	Get(key header.Type) string
	Add(key header.Type, value string)
	Set(key header.Type, value string)
}

// CookiesWrapper is the interface that wraps the basic methods for managing HTTP cookies.
//
// This wrapper provides a unified interface for cookie management, abstracting
// away the details of cookie storage and retrieval.
//
// It allows the library to implement different cookie storage strategies
// (e.g., in-memory, persistent storage) without affecting the public API. It also
// facilitates easier testing and mocking of cookie-related functionality.
//
// Example (for library developers):
//
//	type PersistentCookieWrapper struct {
//		storage CookieStorage
//	}
//
//	func (w *PersistentCookieWrapper) Add(cookie *http.Cookie) {
//		w.storage.Save(cookie)
//		// Additional logic, e.g., expiration handling
//	}
type CookiesWrapper interface {
	Unwrap() []*http.Cookie
	Get(index int) *http.Cookie
	Count() int
	Add(cookie *http.Cookie)
}

// ValidationsWrapper is the interface that wraps the basic methods for managing HTTP request validations.
//
// This wrapper centralizes the handling of validation errors, providing a
// consistent way to accumulate and access errors throughout the request building process.
//
// It allows for more complex validation scenarios, such as conditional validations
// or aggregating errors from multiple sources, while keeping the public API clean and simple.
//
// Example:
//
//	type EnhancedValidationsWrapper struct {
//		errors []error
//		warnings []string
//	}
//
//	func (w *EnhancedValidationsWrapper) AddWarning(warning string) {
//		w.warnings = append(w.warnings, warning)
//	}
type ValidationsWrapper interface {
	Unwrap() []error
	Get(index int) error
	IsEmpty() bool
	Count() int
	Add(err error)
}

// ContextWrapper is the interface that wraps the basic methods for managing HTTP request context.
//
// This wrapper provides a layer of abstraction over the standard context.Context,
// allowing for potential enhancements to context handling without affecting the public API.
//
// It enables the library to implement custom context-related features, such as
// automatic context propagation or context-based tracing, while maintaining a simple interface.
//
// Example:
//
//	type TracingContextWrapper struct {
//		ctx context.Context
//		tracer Tracer
//	}
//
//	func (w *TracingContextWrapper) Unwrap() context.Context {
//		return w.tracer.ContextWithSpan(w.ctx)
//	}
type ContextWrapper interface {
	Unwrap() context.Context
	Set(ctx context.Context)
}

// Compile-time checks.
var (
	_ Client              = (*ClientConfigBase)(nil)
	_ ConfigBaseURL       = (*DefaultBaseURL)(nil)
	_ ConfigBaseURL       = (*BalancedBaseURL)(nil)
	_ HttpClientComponent = (*DefaultHttpClient)(nil)
	_ HeaderWrapper       = (*DefaultHttpHeader)(nil)
	_ CookiesWrapper      = (*DefaultHttpCookies)(nil)
	_ ValidationsWrapper  = (*DefaultValidations)(nil)
	_ ContextWrapper      = (*DefaultContext)(nil)
	_ BodyWrapper         = (*BufferedBody)(nil)
	_ BodyWrapper         = (*UnbufferedBody)(nil)
)
