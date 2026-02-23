package fastshot

import (
	"testing"

	"github.com/opus-domini/fast-shot/constant/header"
	"github.com/opus-domini/fast-shot/constant/mime"
)

func TestClientHeaderBuilder(t *testing.T) {
	const testAgent = "TestAgent"

	tests := []struct {
		name       string
		setup      func(*ClientBuilder)
		assertFunc func(*testing.T, *ClientBuilder)
	}{
		{
			name: "Add single header",
			setup: func(cb *ClientBuilder) {
				cb.Header().Add(header.ContentType, mime.JSON.String())
			},
			assertFunc: func(t *testing.T, cb *ClientBuilder) {
				if got := cb.client.Header().Get(header.ContentType); got != mime.JSON.String() {
					t.Errorf("got %q, want %q", got, mime.JSON.String())
				}
			},
		},
		{
			name: "Add multiple headers",
			setup: func(cb *ClientBuilder) {
				cb.Header().AddAll(map[header.Type]string{
					header.ContentType: mime.JSON.String(),
					header.UserAgent:   testAgent,
				})
			},
			assertFunc: func(t *testing.T, cb *ClientBuilder) {
				if got := cb.client.Header().Get(header.ContentType); got != mime.JSON.String() {
					t.Errorf("ContentType got %q, want %q", got, mime.JSON.String())
				}
				if got := cb.client.Header().Get(header.UserAgent); got != testAgent {
					t.Errorf("UserAgent got %q, want %q", got, testAgent)
				}
			},
		},
		{
			name: "Set single header",
			setup: func(cb *ClientBuilder) {
				cb.Header().Set(header.ContentType, mime.JSON.String())
			},
			assertFunc: func(t *testing.T, cb *ClientBuilder) {
				if got := cb.client.Header().Get(header.ContentType); got != mime.JSON.String() {
					t.Errorf("got %q, want %q", got, mime.JSON.String())
				}
			},
		},
		{
			name: "Set multiple headers",
			setup: func(cb *ClientBuilder) {
				cb.Header().SetAll(map[header.Type]string{
					header.ContentType: mime.JSON.String(),
					header.UserAgent:   testAgent,
				})
			},
			assertFunc: func(t *testing.T, cb *ClientBuilder) {
				if got := cb.client.Header().Get(header.ContentType); got != mime.JSON.String() {
					t.Errorf("ContentType got %q, want %q", got, mime.JSON.String())
				}
				if got := cb.client.Header().Get(header.UserAgent); got != testAgent {
					t.Errorf("UserAgent got %q, want %q", got, testAgent)
				}
			},
		},
		{
			name: "Add Accept header",
			setup: func(cb *ClientBuilder) {
				cb.Header().AddAccept(mime.JSON)
			},
			assertFunc: func(t *testing.T, cb *ClientBuilder) {
				if got := cb.client.Header().Get(header.Accept); got != mime.JSON.String() {
					t.Errorf("got %q, want %q", got, mime.JSON.String())
				}
			},
		},
		{
			name: "Add Content-Type header",
			setup: func(cb *ClientBuilder) {
				cb.Header().AddContentType(mime.JSON)
			},
			assertFunc: func(t *testing.T, cb *ClientBuilder) {
				if got := cb.client.Header().Get(header.ContentType); got != mime.JSON.String() {
					t.Errorf("got %q, want %q", got, mime.JSON.String())
				}
			},
		},
		{
			name: "Add User-Agent header",
			setup: func(cb *ClientBuilder) {
				cb.Header().AddUserAgent(testAgent)
			},
			assertFunc: func(t *testing.T, cb *ClientBuilder) {
				if got := cb.client.Header().Get(header.UserAgent); got != testAgent {
					t.Errorf("got %q, want %q", got, testAgent)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			cb := NewClient("https://example.com")

			// Act
			tt.setup(cb)

			// Assert
			tt.assertFunc(t, cb)
		})
	}
}
