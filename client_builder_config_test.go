package fastshot

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"

	"github.com/opus-domini/fast-shot/mock"
)

func TestClientConfigBuilder(t *testing.T) {
	tests := []struct {
		name           string
		method         func(*ClientBuilder) *ClientBuilder
		setupClient    func(*ClientConfigBase)
		expectedConfig func(*ClientConfigBase) bool
		expectedErrors []error
	}{
		{
			name: "Set custom HTTP client",
			method: func(cb *ClientBuilder) *ClientBuilder {
				mockClient := &mock.HttpClientComponent{}
				return cb.Config().SetCustomHttpClient(mockClient)
			},
			expectedConfig: func(ccb *ClientConfigBase) bool {
				_, ok := ccb.HttpClient().(*mock.HttpClientComponent)
				return ok
			},
		},
		{
			name: "Set custom transport",
			method: func(cb *ClientBuilder) *ClientBuilder {
				transport := &http.Transport{}
				return cb.Config().SetCustomTransport(transport)
			},
			expectedConfig: func(ccb *ClientConfigBase) bool {
				return ccb.HttpClient().Transport() != nil
			},
		},
		{
			name: "Set timeout",
			method: func(cb *ClientBuilder) *ClientBuilder {
				return cb.Config().SetTimeout(5 * time.Second)
			},
			expectedConfig: func(ccb *ClientConfigBase) bool {
				return ccb.HttpClient().Timeout() == 5*time.Second
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			cb := &ClientBuilder{
				client: newClientConfigBase("https://api.example.com"),
			}
			if tt.setupClient != nil {
				tt.setupClient(cb.client.(*ClientConfigBase))
			}

			// Act
			result := tt.method(cb)

			// Assert
			if result != cb {
				t.Errorf("got different builder, want same")
			}
			if !tt.expectedConfig(cb.client.(*ClientConfigBase)) {
				t.Errorf("expectedConfig returned false")
			}

			if tt.expectedErrors != nil {
				if got := len(cb.client.Validations().Unwrap()); got != len(tt.expectedErrors) {
					t.Errorf("validations count got %d, want %d", got, len(tt.expectedErrors))
				}
				for i, expectedErr := range tt.expectedErrors {
					if !errors.Is(cb.client.Validations().Get(i), expectedErr) {
						t.Errorf("validation[%d] got %v, want %v", i, cb.client.Validations().Get(i), expectedErr)
					}
				}
			} else {
				if got := cb.client.Validations().Unwrap(); len(got) != 0 {
					t.Errorf("validations got %v, want empty", got)
				}
			}
		})
	}
}

func TestClientConfigBuilder_SetFollowRedirects(t *testing.T) {
	// Arrange
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/redirect" {
			http.Redirect(w, r, "/target", http.StatusFound)
			return
		}
		_, _ = w.Write([]byte("OK"))
	}))
	defer server.Close()

	tests := []struct {
		name            string
		followRedirects bool
		wantFinalURL    string
		wantStatusCode  int
	}{
		{
			name:            "Follow Redirects",
			followRedirects: true,
			wantFinalURL:    server.URL + "/target",
			wantStatusCode:  http.StatusOK,
		},
		{
			name:            "Do Not Follow Redirects",
			followRedirects: false,
			wantFinalURL:    server.URL + "/redirect",
			wantStatusCode:  http.StatusFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Act
			client := NewClient(server.URL).
				Config().SetFollowRedirects(tt.followRedirects).
				Build()

			resp, err := client.GET("/redirect").Send()

			// Assert
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			defer resp.Body().Close()

			if got := resp.Request().URL(); got != tt.wantFinalURL {
				t.Errorf("URL got %q, want %q", got, tt.wantFinalURL)
			}
			if got := resp.Status().Code(); got != tt.wantStatusCode {
				t.Errorf("StatusCode got %d, want %d", got, tt.wantStatusCode)
			}
		})
	}
}

func TestClientConfigBuilder_SetProxy(t *testing.T) {
	tests := []struct {
		name           string
		baseURL        string
		proxyURL       string
		setupTransport func(*ClientBuilder)
		assertFunc     func(*testing.T, *ClientBuilder)
	}{
		{
			name:     "Set Proxy with Default Transport",
			baseURL:  "https://example.com",
			proxyURL: "http://localhost:8080",
			assertFunc: func(t *testing.T, builder *ClientBuilder) {
				transport, ok := builder.client.HttpClient().Transport().(*http.Transport)
				if !ok {
					t.Fatal("Transport should be of type *http.Transport")
				}
				if transport.Proxy == nil {
					t.Fatal("Proxy function should be set")
				}
			},
		},
		{
			name:     "Set Proxy with Custom Transport",
			baseURL:  "https://example.com",
			proxyURL: "http://localhost:8080",
			setupTransport: func(builder *ClientBuilder) {
				builder.Config().SetCustomTransport(&http.Transport{
					Proxy: http.ProxyURL(&url.URL{
						Scheme: "http",
						Host:   "localhost:9090",
					}),
				})
			},
			assertFunc: func(t *testing.T, builder *ClientBuilder) {
				transport, ok := builder.client.HttpClient().Transport().(*http.Transport)
				if !ok {
					t.Fatal("Transport should be of type *http.Transport")
				}
				if transport.Proxy == nil {
					t.Fatal("Proxy function should be set")
				}

				proxyURL, _ := transport.Proxy(&http.Request{URL: &url.URL{Scheme: "http", Host: "example.com"}})
				if proxyURL.Host != "localhost:8080" {
					t.Errorf("Proxy URL host got %q, want %q", proxyURL.Host, "localhost:8080")
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			builder := NewClient(tt.baseURL)
			if tt.setupTransport != nil {
				tt.setupTransport(builder)
			}

			// Act
			builder.Config().SetProxy(tt.proxyURL)

			// Assert
			tt.assertFunc(t, builder)
		})
	}
}
