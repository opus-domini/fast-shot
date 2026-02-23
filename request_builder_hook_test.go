package fastshot

import (
	"net/http"
	"testing"
)

func TestRequestHookBuilder(t *testing.T) {
	tests := []struct {
		name       string
		setup      func(*RequestBuilder) *RequestBuilder
		assertFunc func(*testing.T, *RequestBuilder)
	}{
		{
			name: "Add before request hook",
			setup: func(rb *RequestBuilder) *RequestBuilder {
				return rb.Hook().OnBeforeRequest(func(req *http.Request) error {
					return nil
				})
			},
			assertFunc: func(t *testing.T, rb *RequestBuilder) {
				if got := len(rb.request.config.BeforeRequestHooks()); got != 1 {
					t.Errorf("before request hooks count got %d, want 1", got)
				}
			},
		},
		{
			name: "Add after response hook",
			setup: func(rb *RequestBuilder) *RequestBuilder {
				return rb.Hook().OnAfterResponse(func(req *http.Request, resp *http.Response) {})
			},
			assertFunc: func(t *testing.T, rb *RequestBuilder) {
				//nolint:bodyclose // False positive: reading hook slice length, not handling a response body.
				if got := len(rb.request.config.AfterResponseHooks()); got != 1 {
					t.Errorf("after response hooks count got %d, want 1", got)
				}
			},
		},
		{
			name: "Add multiple hooks",
			setup: func(rb *RequestBuilder) *RequestBuilder {
				return rb.
					Hook().OnBeforeRequest(func(req *http.Request) error { return nil }).
					Hook().OnBeforeRequest(func(req *http.Request) error { return nil }).
					Hook().OnAfterResponse(func(req *http.Request, resp *http.Response) {}).
					Hook().OnAfterResponse(func(req *http.Request, resp *http.Response) {})
			},
			assertFunc: func(t *testing.T, rb *RequestBuilder) {
				if got := len(rb.request.config.BeforeRequestHooks()); got != 2 {
					t.Errorf("before request hooks count got %d, want 2", got)
				}
				//nolint:bodyclose // False positive: reading hook slice length, not handling a response body.
				if got := len(rb.request.config.AfterResponseHooks()); got != 2 {
					t.Errorf("after response hooks count got %d, want 2", got)
				}
			},
		},
		{
			name: "No hooks by default",
			setup: func(rb *RequestBuilder) *RequestBuilder {
				return rb
			},
			assertFunc: func(t *testing.T, rb *RequestBuilder) {
				if got := len(rb.request.config.BeforeRequestHooks()); got != 0 {
					t.Errorf("before request hooks count got %d, want 0", got)
				}
				//nolint:bodyclose // False positive: reading hook slice length, not handling a response body.
				if got := len(rb.request.config.AfterResponseHooks()); got != 0 {
					t.Errorf("after response hooks count got %d, want 0", got)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			rb := &RequestBuilder{
				request: &Request{
					config: newRequestConfigBase("", ""),
				},
			}

			// Act
			result := tt.setup(rb)

			// Assert
			if result != rb {
				t.Errorf("got different builder, want same")
			}
			tt.assertFunc(t, rb)
		})
	}
}
