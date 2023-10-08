package fastshot

import (
	"encoding/base64"
	"net/http"
	"time"
)

// Client encapsulates HTTP client, HTTP Header, HTTP Cookies, and baseURL.
type Client struct {
	httpClient  *http.Client
	httpHeader  *http.Header
	httpCookies []*http.Cookie
	baseURL     string
}

// ClientBuilder serves as the main entry point for configuring HTTP clients.
type ClientBuilder struct {
	client *Client
}

// ClientAuthBuilder allows for setting authentication configurations.
type ClientAuthBuilder struct {
	parentBuilder *ClientBuilder
}

// ClientHeaderBuilder allows for setting custom HTTP Header.
type ClientHeaderBuilder struct {
	parentBuilder *ClientBuilder
}

// ClientCookieBuilder allows for setting custom HTTP Cookies.
type ClientCookieBuilder struct {
	parentBuilder *ClientBuilder
}

// ClientConfigBuilder allows for setting other client configurations.
type ClientConfigBuilder struct {
	parentBuilder *ClientBuilder
}

// NewClient initializes a new ClientBuilder with a given baseURL.
func NewClient(baseURL string) *ClientBuilder {
	return &ClientBuilder{
		client: &Client{
			httpClient:  &http.Client{},
			httpHeader:  &http.Header{},
			httpCookies: []*http.Cookie{},
			baseURL:     baseURL,
		},
	}
}

// DefaultClient initializes a new default Client with a given baseURL.
func DefaultClient(baseURL string) *Client {
	return NewClient(baseURL).Build()
}

// Auth returns a new ClientAuthBuilder for setting authentication options.
func (b *ClientBuilder) Auth() *ClientAuthBuilder {
	return &ClientAuthBuilder{parentBuilder: b}
}

// Set sets the Authorization header for custom authentication.
func (b *ClientAuthBuilder) Set(value string) *ClientAuthBuilder {
	b.parentBuilder.client.httpHeader.Set("Authorization", value)
	return b
}

// BearerToken sets the Authorization header for Bearer token authentication.
func (b *ClientAuthBuilder) BearerToken(token string) *ClientAuthBuilder {
	b.parentBuilder.client.httpHeader.Set("Authorization", "Bearer "+token)
	return b
}

// BasicAuth sets the Authorization header for Basic authentication.
func (b *ClientAuthBuilder) BasicAuth(username, password string) *ClientAuthBuilder {
	encoded := base64.StdEncoding.EncodeToString([]byte(username + ":" + password))
	b.parentBuilder.client.httpHeader.Set("Authorization", "Basic "+encoded)
	return b
}

// End returns the parent ClientBuilder.
func (b *ClientAuthBuilder) End() *ClientBuilder {
	return b.parentBuilder
}

// Header returns a new ClientHeaderBuilder for setting custom HTTP Header.
func (b *ClientBuilder) Header() *ClientHeaderBuilder {
	return &ClientHeaderBuilder{parentBuilder: b}
}

// Add adds a custom header to the HTTP client. If header already exists, it will be appended.
func (b *ClientHeaderBuilder) Add(key, value string) *ClientHeaderBuilder {
	b.parentBuilder.client.httpHeader.Add(key, value)
	return b
}

// Set sets a custom header to the HTTP client. If header already exists, it will be overwritten.
func (b *ClientHeaderBuilder) Set(key, value string) *ClientHeaderBuilder {
	b.parentBuilder.client.httpHeader.Set(key, value)
	return b
}

// AddAccept sets the Accept header. If header already exists, it will be appended.
func (b *ClientHeaderBuilder) AddAccept(value string) *ClientHeaderBuilder {
	b.parentBuilder.client.httpHeader.Add("Accept", value)
	return b
}

// AddContentType sets the Content-Type header. If header already exists, it will be appended.
func (b *ClientHeaderBuilder) AddContentType(value string) *ClientHeaderBuilder {
	b.parentBuilder.client.httpHeader.Add("Content-Type", value)
	return b
}

// AddUserAgent sets the User-Agent header. If header already exists, it will be appended.
func (b *ClientHeaderBuilder) AddUserAgent(value string) *ClientHeaderBuilder {
	b.parentBuilder.client.httpHeader.Add("User-Agent", value)
	return b
}

// End returns the parent ClientBuilder.
func (b *ClientHeaderBuilder) End() *ClientBuilder {
	return b.parentBuilder
}

// Cookie returns a new ClientCookieBuilder for setting custom HTTP Cookies.
func (b *ClientBuilder) Cookie() *ClientCookieBuilder {
	return &ClientCookieBuilder{parentBuilder: b}
}

// Add adds a custom cookie to the HTTP client.
func (b *ClientCookieBuilder) Add(cookie *http.Cookie) *ClientCookieBuilder {
	b.parentBuilder.client.httpCookies = append(b.parentBuilder.client.httpCookies, cookie)
	return b
}

// End returns the parent ClientBuilder.
func (b *ClientCookieBuilder) End() *ClientBuilder {
	return b.parentBuilder
}

// Config returns a new ClientConfigBuilder for setting custom client configurations.
func (b *ClientBuilder) Config() *ClientConfigBuilder {
	return &ClientConfigBuilder{parentBuilder: b}
}

// SetTimeout sets the timeout for the HTTP client.
func (b *ClientConfigBuilder) SetTimeout(duration time.Duration) *ClientConfigBuilder {
	b.parentBuilder.client.httpClient.Timeout = duration
	return b
}

// SetFollowRedirects controls whether the HTTP client should follow redirects.
func (b *ClientConfigBuilder) SetFollowRedirects(follow bool) *ClientConfigBuilder {
	if !follow {
		b.parentBuilder.client.httpClient.CheckRedirect = func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		}
	}
	return b
}

// SetCustomTransport sets custom transport for the HTTP client.
func (b *ClientConfigBuilder) SetCustomTransport(transport http.RoundTripper) *ClientConfigBuilder {
	b.parentBuilder.client.httpClient.Transport = transport
	return b
}

// End returns the parent ClientBuilder.
func (b *ClientConfigBuilder) End() *ClientBuilder {
	return b.parentBuilder
}

// Build finalizes the ClientBuilder configurations and returns a new Client.
func (b *ClientBuilder) Build() *Client {
	return b.client
}
