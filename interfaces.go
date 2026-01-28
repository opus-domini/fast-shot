package fastshot

import (
	"context"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/opus-domini/fast-shot/constant/header"
	"github.com/opus-domini/fast-shot/constant/mime"
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
//	fmt.Println(response.StatusCode())
//
// The Client interface allows for a fluent, builder-style API that makes it easy to configure
// and send HTTP requests with minimal boilerplate code.
type Client interface {
	ClientConfig
	ClientHttpMethods
}

// ClientConfig is the interface that wraps the basic methods for configuring an HTTP client.
//
// This interface is crucial as it provides a centralized way to configure various aspects of
// the HTTP client, including headers, cookies, validations, and the underlying HTTP client itself.
// It serves as a bridge between the high-level client configuration and the low-level HTTP operations.
//
// Example usage:
//
//	client := fastshot.NewClient("https://api.example.com").
//		Header().Set(header.UserAgent, "MyApp/1.0").
//		Cookie().Add(&http.Cookie{Name: "session", Value: "abc123"}).
//		Config().SetTimeout(10 * time.Second).
//		Build()
//
// The ClientConfig interface allows for a fluent, builder-style API that makes it easy to configure
// all aspects of the HTTP client in a readable and maintainable way. This design promotes
// consistency in client configuration across an application.
type ClientConfig interface {
	ConfigHttpClient
	Header() HeaderWrapper
	Cookies() CookiesWrapper
	Validations() ValidationsWrapper
	ConfigBaseURL
}

// ConfigHttpClient is the interface that wraps the basic methods for configuring the underlying HTTP client.
//
// This interface is essential for providing fine-grained control over the HTTP client used for
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
//		Config().SetHttpClient(customClient).
//		Build()
//
// By providing this level of control, the library can cater to a wide range of use cases,
// from simple API calls to complex scenarios requiring custom HTTP client configurations.
type ConfigHttpClient interface {
	SetHttpClient(httpClient HttpClientComponent)
	HttpClient() HttpClientComponent
}

// HttpClientComponent is the interface that wraps the basic methods for executing HTTP requests.
//
// This interface is crucial as it abstracts the actual HTTP client implementation, allowing
// for easy substitution of the underlying HTTP client (e.g., for testing or using a custom
// implementation). It provides methods to configure key aspects of HTTP communication.
//
// Example usage:
//
//	type CustomClient struct {
//		// Custom implementation
//	}
//
//	func (c *CustomClient) Do(req *http.Request) (*http.Response, error) {
//		// Custom request execution logic
//	}
//
//	... implement other methods ...
//
//	client.Config().SetCustomHttpClient(&CustomClient{})
//
// By abstracting the HTTP client, this interface allows for greater flexibility and
// testability in the library's usage.
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
// This interface is crucial for managing the base URL of the client, which is a fundamental
// part of constructing request URLs. It supports both single base URL configurations and
// load-balanced configurations with multiple base URLs.
//
// Example usage for default client:
//
//	client := fastshot.NewClient("https://api.example.com").Build()
//	fmt.Println(client.BaseURL()) // Output: https://api.example.com
//
// Example usage for load-balanced client:
//
//	client := fastshot.NewClientLoadBalancer([]string{
//		"https://api1.example.com",
//		"https://api2.example.com",
//	}).Build()
//	fmt.Println(client.BaseURL()) // Output: One of the provided URLs, rotating on each call
//
// The ConfigBaseURL interface allows the library to support different base URL strategies,
// enabling features like automatic load balancing or failover between multiple API endpoints.
type ConfigBaseURL interface {
	BaseURL() *url.URL
}

// ClientHttpMethods is the interface that wraps the basic HTTP methods for making requests.
//
// This interface is fundamental to the library as it provides a clean, method-based API for
// initiating different types of HTTP requests. It abstracts away the complexities of constructing
// HTTP requests, allowing users to focus on the specific HTTP method they need.
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
//
//	// PUT request with custom header
//	putResp, err := client.PUT("/users/123").
//		Header().Set(header.ContentType, mime.ApplicationJSON).
//		Body().AsString(`{"name": "Jane Doe"}`).
//		Send()
//
// By providing separate methods for each HTTP verb, this interface makes the API more
// intuitive and less error-prone. It also allows for method-specific optimizations or
// behaviors if needed in the future.
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

// BuilderHeader is the interface that wraps the basic methods for setting custom HTTP headers.
//
// This interface is crucial for customizing request headers, which is a common requirement
// in many API interactions. It provides a fluent API for adding or setting headers individually
// or in bulk, with special methods for common headers like Accept and Content-Type.
//
// Example usage:
//
//	client.
//		Header().Set(header.ContentType, mime.ApplicationJSON).
//		Header().Add(header.UserAgent, "MyApp/1.0").
//		Header().AddAccept(mime.ApplicationJSON)
//
// The generic type parameter T allows this interface to be used with both Client and RequestBuilder,
// enabling header configuration at both the client and request level.
type BuilderHeader[T any] interface {
	Add(key header.Type, value string) *T
	AddAll(headers map[header.Type]string) *T
	Set(key header.Type, value string) *T
	SetAll(headers map[header.Type]string) *T
	AddAccept(value mime.Type) *T
	AddContentType(value mime.Type) *T
	AddUserAgent(value string) *T
}

// BuilderCookie is the interface that wraps the basic method for adding HTTP cookies.
//
// This interface is crucial for managing cookies in HTTP requests, which is essential for
// maintaining session state and other cookie-based authentication or tracking mechanisms.
// It provides a simple way to add cookies to either the client (for all requests) or to
// individual requests.
//
// Example usage for client-level cookies:
//
//	client := fastshot.NewClient("https://api.example.com").
//		Cookie().Add(&http.Cookie{Name: "session", Value: "abc123"}).
//		Build()
//
// Example usage for request-level cookies:
//
//	response, err := client.GET("/protected-resource").
//		Cookie().Add(&http.Cookie{Name: "csrf_token", Value: "xyz789"}).
//		Send()
//
// The BuilderCookie interface allows for easy management of cookies, supporting both
// persistent cookies across all requests and one-time cookies for specific requests.
type BuilderCookie[T any] interface {
	Add(cookie *http.Cookie) *T
}

// BuilderAuth is the interface that wraps the basic methods for setting HTTP authentication.
//
// This interface is essential for implementing various authentication schemes in HTTP requests.
// It provides methods for setting custom authentication headers, as well as convenience methods
// for common auth types like Bearer token and Basic auth.
//
// Example usage for Bearer token auth:
//
//	client := fastshot.NewClient("https://api.example.com").
//		Auth().BearerToken("my-access-token").
//		Build()
//
// Example usage for Basic auth:
//
//	response, err := client.GET("/protected-resource").
//		Auth().BasicAuth("username", "password").
//		Send()
//
// The BuilderAuth interface simplifies the process of adding authentication to requests,
// reducing the likelihood of errors in implementing common auth schemes.
type BuilderAuth[T any] interface {
	Set(value string) *T
	BearerToken(token string) *T
	BasicAuth(username, password string) *T
}

// BuilderHttpClientConfig is the interface that wraps the basic methods for configuring the HTTP client.
//
// This interface is crucial for fine-tuning the behavior of the underlying HTTP client.
// It provides methods to set custom HTTP clients, transports, timeouts, redirect behavior,
// and proxy settings, allowing for advanced customization of the HTTP communication.
//
// Example usage:
//
//	client := fastshot.NewClient("https://api.example.com").
//		Config().SetTimeout(30 * time.Second).
//		Config().SetFollowRedirects(false).
//		Config().SetProxy("http://proxy.example.com:8080").
//		Build()
//
// The BuilderHttpClientConfig interface enables users to adapt the HTTP client to various
// network conditions and security requirements, enhancing the library's flexibility.
type BuilderHttpClientConfig[T any] interface {
	SetCustomHttpClient(httpClient HttpClientComponent) *T
	SetCustomTransport(transport http.RoundTripper) *T
	SetTimeout(duration time.Duration) *T
	SetFollowRedirects(follow bool) *T
	SetProxy(proxyURL string) *T
}

// BuilderRequestContext is the interface that wraps the basic method for setting the request context.
//
// This interface is essential for managing request-specific contexts, which are crucial for
// implementing timeouts, cancellations, and passing request-scoped values. It allows users
// to set a custom context for individual requests.
//
// Example usage:
//
//	ctx, cancel := context.WithTimeout(context.Background(), 5 * time.Second)
//	defer cancel()
//
//	response, err := client.GET("/long-running-operation").
//		Context().Set(ctx).
//		Send()
//
// The BuilderRequestContext interface enables fine-grained control over request lifecycle,
// improving the ability to manage long-running requests or implement advanced cancellation logic.
type BuilderRequestContext[T any] interface {
	Set(ctx context.Context) *T
}

// BuilderRequestBody is the interface that wraps the basic methods for setting the request body.
//
// This interface is essential for sending data in requests (e.g., POST, PUT). It provides
// flexibility in how the body can be set, supporting raw io.Reader, string, JSON serialization,
// and multipart/form-data.
//
// Example usage:
//
//	type User struct {
//		Name  string `json:"name"`
//		Email string `json:"email"`
//	}
//
//	user := User{Name: "Fulano", Email: "fulano@example.com"}
//	response, err := client.POST("/users").
//		Body().AsJSON(user).
//		Send()
//
// Example usage with form data:
//
//	response, err := client.POST("/login").
//		Body().AsFormData(map[string]string{
//			"username": "user",
//			"password": "pass",
//		}).
//		Send()
//
// The ability to set the body as JSON or form data directly is particularly useful for API interactions,
// reducing boilerplate code for serialization.
type BuilderRequestBody[T any] interface {
	AsReader(body io.Reader) *T
	AsString(body string) *T
	AsJSON(obj interface{}) *T
	AsXML(obj interface{}) *T
	AsFormData(fields map[string]string) *T
}

// BuilderRequestQuery is the interface that wraps the basic methods for setting query parameters.
//
// This interface is crucial for constructing URL query strings in a clean and type-safe manner.
// It provides methods to add and set individual query parameters, as well as methods to
// set multiple parameters at once or from a raw query string.
//
// Example usage:
//
//	response, err := client.GET("/search").
//		Query().AddParam("q", "golang").
//		Query().AddParam("sort", "relevance").
//		Query().SetParam("page", "1").
//		Send()
//
// Example usage with raw query string:
//
//	response, err := client.GET("/search").
//		Query().SetRawString("q=golang&sort=relevance&page=1").
//		Send()
//
// The BuilderRequestQuery interface simplifies the process of building complex query strings,
// reducing errors and improving readability when working with URL parameters.
type BuilderRequestQuery[T any] interface {
	AddParam(param, value string) *T
	AddParams(params map[string]string) *T
	SetParam(param, value string) *T
	SetParams(params map[string]string) *T
	SetRawString(query string) *T
}

// BuilderRequestRetry is the interface that wraps the basic methods for configuring request retries.
//
// Retry functionality is crucial for building robust HTTP clients that can handle transient
// failures. This interface provides a variety of retry strategies, including constant and
// exponential backoff, with optional jitter for avoiding thundering herd problems.
//
// Example usage:
//
//	 interval := 100 * time.Millisecond
//	 maxAttempts := 3
//	 backoffRate := 2.0
//
//		response, err := client.GET("/users").
//			Retry().SetExponentialBackoffWithJitter(interval, maxAttempts, backoffRate).
//			Retry().WithRetryCondition(func(resp Response) bool {
//				return resp.Is5xxServerError()
//			}).
//			Send()
//
// The flexibility in retry configuration allows users to fine-tune the retry behavior
// to their specific needs, improving the reliability of their HTTP requests.
type BuilderRequestRetry[T any] interface {
	SetConstantBackoff(interval time.Duration, maxAttempts uint) *T
	SetConstantBackoffWithJitter(interval time.Duration, maxAttempts uint) *T
	SetExponentialBackoff(interval time.Duration, maxAttempts uint, backoffRate float64) *T
	SetExponentialBackoffWithJitter(interval time.Duration, maxAttempts uint, backoffRate float64) *T
	WithRetryCondition(shouldRetry func(response *Response) bool) *T
	WithMaxDelay(duration time.Duration) *T
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

// BodyWrapper is the interface that wraps the basic methods for handling request and response bodies.
//
// This interface provides a unified way to read, write, and manipulate body content
// in various formats (e.g., raw bytes, string, JSON). It abstracts the underlying
// implementation details, allowing for different body handling strategies (e.g., buffered
// or streaming) without changing the public API.
//
// Example:
//
//	body := newBufferedBody()
//	err := body.WriteAsJSON(map[string]string{"key": "value"})
//	if err != nil {
//	    // Handle error
//	}
//
//	var result map[string]interface{}
//	err = body.ReadAsJSON(&result)
//	if err != nil {
//	    // Handle error
//	}
//
// The BodyWrapper interface allows for efficient and flexible handling of HTTP request
// and response bodies, supporting various content types and processing requirements.
type BodyWrapper interface {
	io.ReadCloser
	ReadAsJSON(obj interface{}) error
	WriteAsJSON(obj interface{}) error
	ReadAsXML(obj interface{}) error
	WriteAsXML(obj interface{}) error
	ReadAsString() (string, error)
	WriteAsString(body string) error
	Set(body io.Reader) error
	Unwrap() io.Reader
}
