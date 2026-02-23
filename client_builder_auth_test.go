package fastshot

import (
	"encoding/base64"
	"testing"

	"github.com/opus-domini/fast-shot/constant/header"
)

func TestClientAuthBuilder(t *testing.T) {
	tests := []struct {
		name           string
		method         func(*ClientBuilder) *ClientBuilder
		expectedHeader string
	}{
		{
			name: "Set custom auth",
			method: func(cb *ClientBuilder) *ClientBuilder {
				return cb.Auth().Set("Custom auth-token")
			},
			expectedHeader: "Custom auth-token",
		},
		{
			name: "Set bearer token",
			method: func(cb *ClientBuilder) *ClientBuilder {
				return cb.Auth().BearerToken("my-token")
			},
			expectedHeader: "Bearer my-token",
		},
		{
			name: "Set basic auth",
			method: func(cb *ClientBuilder) *ClientBuilder {
				return cb.Auth().BasicAuth("username", "password")
			},
			expectedHeader: "Basic " + base64.StdEncoding.EncodeToString([]byte("username:password")),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			cb := &ClientBuilder{
				client: newClientConfigBase("https://api.example.com"),
			}

			// Act
			result := tt.method(cb)

			// Assert
			if result != cb {
				t.Errorf("got different builder, want same")
			}
			if got := cb.client.Header().Get(header.Authorization); got != tt.expectedHeader {
				t.Errorf("Authorization got %q, want %q", got, tt.expectedHeader)
			}
		})
	}
}
