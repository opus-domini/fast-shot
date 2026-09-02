package fastshot

// ClientBuilder serves as the main entry point for configuring HTTP clients.
type ClientBuilder struct {
	client Client
}

// NewClient initializes a new ClientBuilder with a given baseURL.
func NewClient(baseURL string) *ClientBuilder {
	return &ClientBuilder{
		client: newClientConfigBase(baseURL),
	}
}

// NewClientLoadBalancer initializes a new ClientBuilder with a given baseURLs.
func NewClientLoadBalancer(baseURLs []string) *ClientBuilder {
	return &ClientBuilder{
		client: newBalancedClientConfigBase(baseURLs),
	}
}

// DefaultClient initializes a new default ClientConfig with a given baseURL.
func DefaultClient(baseURL string) ClientHttpMethods {
	return NewClient(baseURL).Build()
}

// DefaultClientLoadBalancer initializes a new default ClientConfig with a given baseURLs.
func DefaultClientLoadBalancer(baseURLs []string) ClientHttpMethods {
	return NewClientLoadBalancer(baseURLs).Build()
}

// Build finalizes the ClientBuilder configurations and returns a new ClientConfig.
func (b *ClientBuilder) Build() ClientHttpMethods {
	return b.client
}
