package fastshot

import (
	"encoding/base64"
	"testing"

	"github.com/opus-domini/fast-shot/constant/header"
)

func TestAuthBuilder_Request(t *testing.T) {
	tests := []struct {
		name           string
		method         func(*AuthBuilder[*RequestBuilder]) *RequestBuilder
		expectedHeader string
	}{
		{
			name: "Set custom auth",
			method: func(ab *AuthBuilder[*RequestBuilder]) *RequestBuilder {
				return ab.Set("Custom auth-token")
			},
			expectedHeader: "Custom auth-token",
		},
		{
			name: "Set bearer token",
			method: func(ab *AuthBuilder[*RequestBuilder]) *RequestBuilder {
				return ab.BearerToken("my-token")
			},
			expectedHeader: "Bearer my-token",
		},
		{
			name: "Set basic auth",
			method: func(ab *AuthBuilder[*RequestBuilder]) *RequestBuilder {
				return ab.BasicAuth("username", "password")
			},
			expectedHeader: "Basic " + base64.StdEncoding.EncodeToString([]byte("username:password")),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			rb := &RequestBuilder{
				request: &Request{
					config: newRequestConfigBase("", "", DefaultJSONCodec()),
				},
			}

			// Act
			result := tt.method(rb.Auth())

			// Assert
			if result != rb {
				t.Errorf("got different builder, want same")
			}
			if got := rb.request.config.Header().Get(header.Authorization); got != tt.expectedHeader {
				t.Errorf("Authorization got %q, want %q", got, tt.expectedHeader)
			}
		})
	}
}

func TestAuthBuilder_Client(t *testing.T) {
	tests := []struct {
		name           string
		method         func(*AuthBuilder[*ClientBuilder]) *ClientBuilder
		expectedHeader string
	}{
		{
			name: "Set custom auth",
			method: func(ab *AuthBuilder[*ClientBuilder]) *ClientBuilder {
				return ab.Set("Custom auth-token")
			},
			expectedHeader: "Custom auth-token",
		},
		{
			name: "Set bearer token",
			method: func(ab *AuthBuilder[*ClientBuilder]) *ClientBuilder {
				return ab.BearerToken("my-token")
			},
			expectedHeader: "Bearer my-token",
		},
		{
			name: "Set basic auth",
			method: func(ab *AuthBuilder[*ClientBuilder]) *ClientBuilder {
				return ab.BasicAuth("username", "password")
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
			result := tt.method(cb.Auth())

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
