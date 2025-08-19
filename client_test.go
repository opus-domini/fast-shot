package fastshot

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/opus-domini/fast-shot/constant/method"
	"github.com/stretchr/testify/assert"
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
			assert.Equal(
				t,
				!clientBuilder.client.Validations().IsEmpty(),
				tt.expectError,
				"NewClient() error = %v, expectError %v",
				!clientBuilder.client.Validations().IsEmpty(),
				tt.expectError,
			)
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
			if (clientBuilder.client.Validations().Count() > 0) != tt.expectError {
				t.Errorf("NewClientLoadBalancer() error = %v, expectError %v", clientBuilder.client.Validations().Count() > 0, tt.expectError)
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
			methodType method.Type
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
			{method.Parse("TRACE"), func(url string) *RequestBuilder { return client.client.TRACE(url) }},
		}

		for _, tt := range tests {
			t.Run(client.name+" "+tt.methodType.String(), func(t *testing.T) {
				req := tt.methodFunc("/")
				resp, _ := req.Send()
				if resp.Status().IsError() {
					t.Errorf("Expected 200, got %d", resp.Status().Code())
				}
			})
		}
	}
}
