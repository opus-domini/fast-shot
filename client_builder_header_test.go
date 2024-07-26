package fastshot

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/opus-domini/fast-shot/constant/header"
	"github.com/opus-domini/fast-shot/constant/mime"
)

func TestClientHeaderBuilder(t *testing.T) {
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
				assert.Equal(t, mime.JSON.String(), cb.client.Header().Get(header.ContentType))
			},
		},
		{
			name: "Add multiple headers",
			setup: func(cb *ClientBuilder) {
				cb.Header().AddAll(map[header.Type]string{
					header.ContentType: mime.JSON.String(),
					header.UserAgent:   "TestAgent",
				})
			},
			assertFunc: func(t *testing.T, cb *ClientBuilder) {
				assert.Equal(t, mime.JSON.String(), cb.client.Header().Get(header.ContentType))
				assert.Equal(t, "TestAgent", cb.client.Header().Get(header.UserAgent))
			},
		},
		{
			name: "Set single header",
			setup: func(cb *ClientBuilder) {
				cb.Header().Set(header.ContentType, mime.JSON.String())
			},
			assertFunc: func(t *testing.T, cb *ClientBuilder) {
				assert.Equal(t, mime.JSON.String(), cb.client.Header().Get(header.ContentType))
			},
		},
		{
			name: "Set multiple headers",
			setup: func(cb *ClientBuilder) {
				cb.Header().SetAll(map[header.Type]string{
					header.ContentType: mime.JSON.String(),
					header.UserAgent:   "TestAgent",
				})
			},
			assertFunc: func(t *testing.T, cb *ClientBuilder) {
				assert.Equal(t, mime.JSON.String(), cb.client.Header().Get(header.ContentType))
				assert.Equal(t, "TestAgent", cb.client.Header().Get(header.UserAgent))
			},
		},
		{
			name: "Add Accept header",
			setup: func(cb *ClientBuilder) {
				cb.Header().AddAccept(mime.JSON)
			},
			assertFunc: func(t *testing.T, cb *ClientBuilder) {
				assert.Equal(t, mime.JSON.String(), cb.client.Header().Get(header.Accept))
			},
		},
		{
			name: "Add Content-Type header",
			setup: func(cb *ClientBuilder) {
				cb.Header().AddContentType(mime.JSON)
			},
			assertFunc: func(t *testing.T, cb *ClientBuilder) {
				assert.Equal(t, mime.JSON.String(), cb.client.Header().Get(header.ContentType))
			},
		},
		{
			name: "Add User-Agent header",
			setup: func(cb *ClientBuilder) {
				cb.Header().AddUserAgent("TestAgent")
			},
			assertFunc: func(t *testing.T, cb *ClientBuilder) {
				assert.Equal(t, "TestAgent", cb.client.Header().Get(header.UserAgent))
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
