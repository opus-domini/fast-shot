package fastshot

import (
	"encoding/base64"
	"testing"
)

func TestClientAuthBuilder_Set(t *testing.T) {
	// Arrange
	builder := NewClient("https://example.com")
	// Act
	authBuilder := builder.Auth()
	authBuilder.Set("value")
	// Assert
	if builder.client.httpHeader.Get("Authorization") != "value" {
		t.Errorf("Authorization header not set correctly")
	}
}

func TestClientAuthBuilder_BasicAuth(t *testing.T) {
	// Arrange
	builder := NewClient("https://example.com")
	// Act
	authBuilder := builder.Auth()
	authBuilder.BasicAuth("username", "password")
	// Assert
	expected := "Basic " + base64.StdEncoding.EncodeToString([]byte("username:password"))
	if builder.client.httpHeader.Get("Authorization") != expected {
		t.Errorf(
			"Header not set correctly, got: %s, want: %s",
			builder.client.httpHeader.Get("Authorization"),
			expected,
		)
	}
}

func TestClientAuthBuilder_BearerToken(t *testing.T) {
	// Arrange
	builder := NewClient("https://example.com")
	// Act
	authBuilder := builder.Auth()
	authBuilder.BearerToken("token")
	// Assert
	if builder.client.httpHeader.Get("Authorization") != "Bearer token" {
		t.Errorf(
			"Header not set correctly, got: %s, want: %s",
			builder.client.httpHeader.Get("Authorization"),
			"Bearer token",
		)
	}
}
