package fastshot

import (
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
	Header() http.Header
	Cookies() []*http.Cookie
	Validations() []error
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
// Example usage:
//
//	customClient := &http.Client{
//		Timeout: 30 * time.Second,
//	}
//
//	client := fastshot.NewClient("https://api.example.com").
//		Config().SetCustomHttpClient(fastshot.NewHttpClientComponent(customClient)).
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

// NewHttpClientComponent adapts an *http.Client to the HttpClientComponent interface.
func NewHttpClientComponent(client *http.Client) HttpClientComponent {
	return &DefaultHttpClient{client: client}
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

// Compile-time checks for the interfaces kept in this refactor.
var (
	_ Client              = (*ClientConfigBase)(nil)
	_ ConfigBaseURL       = (*DefaultBaseURL)(nil)
	_ ConfigBaseURL       = (*BalancedBaseURL)(nil)
	_ HttpClientComponent = (*DefaultHttpClient)(nil)
	_ BodyWrapper         = (*BufferedBody)(nil)
	_ BodyWrapper         = (*UnbufferedBody)(nil)

	// header.Type and mime.Type keep their meaning in the fluent builders.
	_ header.Type = header.Accept
	_ mime.Type   = mime.JSON
)
