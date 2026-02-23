package fastshot

import (
	"net/http"
	"testing"
)

func TestClientHookBuilder(t *testing.T) {
	tests := []struct {
		name       string
		setup      func(*ClientBuilder) *ClientBuilder
		assertFunc func(*testing.T, *ClientBuilder)
	}{
		{
			name: "Add before request hook",
			setup: func(cb *ClientBuilder) *ClientBuilder {
				return cb.Hook().OnBeforeRequest(func(req *http.Request) error {
					return nil
				})
			},
			assertFunc: func(t *testing.T, cb *ClientBuilder) {
				if got := len(cb.client.BeforeRequestHooks()); got != 1 {
					t.Errorf("before request hooks count got %d, want 1", got)
				}
			},
		},
		{
			name: "Add after response hook",
			setup: func(cb *ClientBuilder) *ClientBuilder {
				return cb.Hook().OnAfterResponse(func(req *http.Request, resp *http.Response) {})
			},
			assertFunc: func(t *testing.T, cb *ClientBuilder) {
				//nolint:bodyclose // False positive: reading hook slice length, not handling a response body.
				if got := len(cb.client.AfterResponseHooks()); got != 1 {
					t.Errorf("after response hooks count got %d, want 1", got)
				}
			},
		},
		{
			name: "Add multiple hooks",
			setup: func(cb *ClientBuilder) *ClientBuilder {
				return cb.
					Hook().OnBeforeRequest(func(req *http.Request) error { return nil }).
					Hook().OnBeforeRequest(func(req *http.Request) error { return nil }).
					Hook().OnAfterResponse(func(req *http.Request, resp *http.Response) {}).
					Hook().OnAfterResponse(func(req *http.Request, resp *http.Response) {})
			},
			assertFunc: func(t *testing.T, cb *ClientBuilder) {
				if got := len(cb.client.BeforeRequestHooks()); got != 2 {
					t.Errorf("before request hooks count got %d, want 2", got)
				}
				//nolint:bodyclose // False positive: reading hook slice length, not handling a response body.
				if got := len(cb.client.AfterResponseHooks()); got != 2 {
					t.Errorf("after response hooks count got %d, want 2", got)
				}
			},
		},
		{
			name: "No hooks by default",
			setup: func(cb *ClientBuilder) *ClientBuilder {
				return cb
			},
			assertFunc: func(t *testing.T, cb *ClientBuilder) {
				if got := len(cb.client.BeforeRequestHooks()); got != 0 {
					t.Errorf("before request hooks count got %d, want 0", got)
				}
				//nolint:bodyclose // False positive: reading hook slice length, not handling a response body.
				if got := len(cb.client.AfterResponseHooks()); got != 0 {
					t.Errorf("after response hooks count got %d, want 0", got)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			cb := NewClient("https://example.com")

			// Act
			result := tt.setup(cb)

			// Assert
			if result != cb {
				t.Errorf("got different builder, want same")
			}
			tt.assertFunc(t, cb)
		})
	}
}
