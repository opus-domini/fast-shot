package fastshot

import (
	"encoding/base64"
	"testing"

	"github.com/opus-domini/fast-shot/constant/header"
)

func TestRequestAuthBuilder(t *testing.T) {
	tests := []struct {
		name           string
		method         func(*RequestBuilder) *RequestBuilder
		expectedHeader string
	}{
		{
			name: "Set custom auth",
			method: func(rb *RequestBuilder) *RequestBuilder {
				return rb.Auth().Set("Custom auth-token")
			},
			expectedHeader: "Custom auth-token",
		},
		{
			name: "Set bearer token",
			method: func(rb *RequestBuilder) *RequestBuilder {
				return rb.Auth().BearerToken("my-token")
			},
			expectedHeader: "Bearer my-token",
		},
		{
			name: "Set basic auth",
			method: func(rb *RequestBuilder) *RequestBuilder {
				return rb.Auth().BasicAuth("username", "password")
			},
			expectedHeader: "Basic " + base64.StdEncoding.EncodeToString([]byte("username:password")),
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
			result := tt.method(rb)

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
