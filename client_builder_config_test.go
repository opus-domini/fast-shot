package fastshot

import (
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

func TestClientConfigBuilder_SetTimeout(t *testing.T) {
	// Arrange
	builder := NewClient("https://example.com")
	// Act
	builder.Config().SetTimeout(1)
	// Assert
	if builder.client.HttpClient().Timeout != 1 {
		t.Errorf("Timeout not set correctly")
	}
}

func TestClientConfigBuilder_SetFollowRedirects(t *testing.T) {
	// Arrange
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/redirect" {
			http.Redirect(w, r, "/target", http.StatusFound) // StatusFound 302
			return
		}
		_, _ = w.Write([]byte("OK"))
	}))
	defer server.Close()

	tests := []struct {
		name            string
		followRedirects bool
		wantFinalURL    string
	}{
		{
			name:            "Follow Redirects",
			followRedirects: true,
			wantFinalURL:    server.URL + "/target",
		},
		{
			name:            "Do Not Follow Redirects",
			followRedirects: false,
			wantFinalURL:    server.URL + "/redirect",
		},
	}

	// Act
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := NewClient(server.URL).
				Config().SetFollowRedirects(tt.followRedirects).
				Build()

			resp, err := client.GET("/redirect").Send()
			if err != nil {
				t.Fatalf("Failed: %v", err)
			}

			defer func(Body io.ReadCloser) {
				_ = Body.Close()
			}(resp.RawBody())

			// Assert
			if resp.RawResponse.Request.URL.String() != tt.wantFinalURL {
				t.Errorf("Final URL = %v, want %v", resp.RawResponse.Request.URL, tt.wantFinalURL)
			}
		})
	}
}

func TestClientConfigBuilder_SetCustomTransport(t *testing.T) {
	// Arrange
	builder := NewClient("https://example.com")
	// Act
	builder.Config().SetCustomTransport(&http.Transport{})
	// Assert
	if builder.client.HttpClient().Transport == nil {
		t.Errorf("Transport not set correctly")
	}
}

func TestClientConfigBuilder_SetProxy(t *testing.T) {
	// Arrange
	builder := NewClient("https://example.com")
	// Act
	builder.Config().SetProxy("http://localhost:8080")
	// Assert
	if builder.client.HttpClient().Transport == nil {
		t.Errorf("Transport not set correctly")
	}
}

func TestClientConfigBuilder_SetProxy_WithCustomTransport(t *testing.T) {
	// Arrange
	builder := NewClient("https://example.com")
	// Act
	builder.Config().SetCustomTransport(
		&http.Transport{
			Proxy: http.ProxyURL(&url.URL{
				Scheme: "http",
				Host:   "localhost:9090",
			}),
		}).
		Config().SetProxy("http://localhost:8080")
	// Assert
	if builder.client.HttpClient().Transport == nil {
		t.Errorf("Transport not set correctly")
	}
}

func TestClientConfigBuilder_SetProxy_WithProxyURLParserError(t *testing.T) {
	// Arrange
	builder := NewClient("https://example.com")
	// Act
	builder.Config().SetProxy(":%^:")
	// Assert
	if builder.client.HttpClient().Transport != nil {
		t.Errorf("Transport should not be set")
	}
	if len(builder.client.Validations()) != 1 {
		t.Errorf("Validation for proxy URL parser error not set correctly")
	}
}
