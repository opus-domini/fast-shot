package fastshot

import (
	"net/http"
	"time"
)

// Compile-time check that DefaultHttpClient implements HttpClientComponent.
var _ HttpClientComponent = (*DefaultHttpClient)(nil)

// DefaultHttpClient implements HttpClientComponent interface and provides a default HTTP client.
type DefaultHttpClient struct {
	client *http.Client
}

// Do will execute the *http.Client Do method
func (c *DefaultHttpClient) Do(req *http.Request) (*http.Response, error) {
	return c.client.Do(req)
}

// Transport will return the underlying transport type
func (c *DefaultHttpClient) Transport() http.RoundTripper {
	return c.client.Transport
}

// SetTransport sets the Transport field on the underlying http.Client type
func (c *DefaultHttpClient) SetTransport(transport http.RoundTripper) {
	c.client.Transport = transport
}

// Timeout will return the underlying timeout value
func (c *DefaultHttpClient) Timeout() time.Duration {
	return c.client.Timeout
}

// SetTimeout sets the Timeout field on the underlying http.Client type
func (c *DefaultHttpClient) SetTimeout(duration time.Duration) {
	c.client.Timeout = duration
}

// SetFollowRedirects sets the CheckRedirect field on the underlying http.Client type
func (c *DefaultHttpClient) SetFollowRedirects(follow bool) {
	if !follow {
		c.client.CheckRedirect = func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		}
	}
}

// newDefaultHttpClient initializes a new DefaultHttpClient.
func newDefaultHttpClient() *DefaultHttpClient {
	return &DefaultHttpClient{
		client: &http.Client{},
	}
}
