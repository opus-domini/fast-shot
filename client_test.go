package fastshot

import (
	"github.com/opus-domini/fast-shot/constant/method"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestNewClient(t *testing.T) {
	tests := []struct {
		name        string
		baseURL     string
		expectError bool
	}{
		{
			name:        "Successful Client Creation",
			baseURL:     "https://example.com",
			expectError: false,
		},
		{
			name:        "Error Parsing URL",
			baseURL:     ":%^:",
			expectError: true,
		},
		{
			name:        "Empty URL",
			baseURL:     "",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			clientBuilder := NewClient(tt.baseURL)
			if (len(clientBuilder.client.Validations()) > 0) != tt.expectError {
				t.Errorf("NewClient() error = %v, expectError %v", len(clientBuilder.client.Validations()) > 0, tt.expectError)
			}
		})
	}
}

func TestNewClientLoadBalancer(t *testing.T) {
	tests := []struct {
		name        string
		baseURLs    []string
		expectError bool
	}{
		{
			name:        "Successful Client Load Balancer Creation",
			baseURLs:    []string{"https://example1.com", "https://example2.com"},
			expectError: false,
		},
		{
			name:        "Error Parsing URL",
			baseURLs:    []string{"https://example1.com", ":%^:"},
			expectError: true,
		},
		{
			name:        "Empty URL",
			baseURLs:    []string{"https://example1.com", ""},
			expectError: true,
		},
		{
			name:        "Empty URL List",
			baseURLs:    []string{},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			clientBuilder := NewClientLoadBalancer(tt.baseURLs)
			if (len(clientBuilder.client.Validations()) > 0) != tt.expectError {
				t.Errorf("NewClientLoadBalancer() error = %v, expectError %v", len(clientBuilder.client.Validations()) > 0, tt.expectError)
			}
		})
	}
}

func TestClientMethods(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, err := w.Write([]byte("OK"))
		if err != nil {
			return
		}
	}))
	defer server.Close()

	clients := []struct {
		name   string
		client ClientHttpMethods
	}{
		{"clientDefault", DefaultClient(server.URL)},
		{"clientLoadBalancer", DefaultClientLoadBalancer([]string{server.URL})},
	}

	for _, client := range clients {
		tests := []struct {
			name       string
			methodFunc func(string) *RequestBuilder
		}{
			{method.CONNECT, func(url string) *RequestBuilder { return client.client.CONNECT(url) }},
			{method.DELETE, func(url string) *RequestBuilder { return client.client.DELETE(url) }},
			{method.GET, func(url string) *RequestBuilder { return client.client.GET(url) }},
			{method.HEAD, func(url string) *RequestBuilder { return client.client.HEAD(url) }},
			{method.OPTIONS, func(url string) *RequestBuilder { return client.client.OPTIONS(url) }},
			{method.PATCH, func(url string) *RequestBuilder { return client.client.PATCH(url) }},
			{method.POST, func(url string) *RequestBuilder { return client.client.POST(url) }},
			{method.PUT, func(url string) *RequestBuilder { return client.client.PUT(url) }},
			{method.TRACE, func(url string) *RequestBuilder { return client.client.TRACE(url) }},
		}

		for _, tt := range tests {
			t.Run(client.name+" "+tt.name, func(t *testing.T) {
				req := tt.methodFunc("/")
				resp, _ := req.Send()
				if resp.IsError() {
					t.Errorf("Expected 200, got %d", resp.StatusCode())
				}
			})
		}
	}
}
