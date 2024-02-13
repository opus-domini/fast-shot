package fastshot

import (
	"net/http"
	"time"
)

// HTTPClient decorates an *http.Client types to provide a RawClient interface
type HTTPClient struct {
	client *http.Client
}

var _ RawClient = &HTTPClient{}

// Do will execute the *http.Client Do method
func (c *HTTPClient) Do(req *http.Request) (*http.Response, error) {
	return c.client.Do(req)
}

// SetTransport sets the Transport field on the underlying http.Client type
func (c *HTTPClient) SetTransport(transport http.RoundTripper) {
	c.client.Transport = transport
}

// Transport will return the underlying transport type
func (c *HTTPClient) Transport() http.RoundTripper {
	return c.client.Transport
}

// SetTimeout sets the Timeout field on the underlying http.Client type
func (c *HTTPClient) SetTimeout(duration time.Duration) {
	c.client.Timeout = duration
}

// Timeout will return the underlying timeout value
func (c *HTTPClient) Timeout() time.Duration {
	return c.client.Timeout
}

// SetCheckRedirect sets the CheckRedirect field on the underlying http.Client type
func (c *HTTPClient) SetCheckRedirect(f func(*http.Request, []*http.Request) error) {
	c.client.CheckRedirect = f
}

// NewHTTPClient will create a wrapped *http.Client that confirms to the RawClient interface
func NewHTTPClient(client *http.Client) *HTTPClient {
	return &HTTPClient{client: client}
}
