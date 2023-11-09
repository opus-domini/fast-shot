package fastshot

import (
	"net/http"
)

// Client encapsulates HTTP client, HTTP Header, HTTP Cookies, and baseURL.
type Client struct {
	httpClient  *http.Client
	httpHeader  *http.Header
	httpCookies []*http.Cookie
	baseURL     string
	validations []error
}

// ClientBuilder serves as the main entry point for configuring HTTP clients.
type ClientBuilder struct {
	client *Client
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

// Build finalizes the ClientBuilder configurations and returns a new Client.
func (b *ClientBuilder) Build() *Client {
	return b.client
}
