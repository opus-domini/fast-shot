package fastshot

import (
	"encoding/base64"
	"testing"
)

func TestRequestAuthBuilder_Set(t *testing.T) {
	// Arrange
	client := DefaultClient("https://example.com")
	// Act
	builder := client.GET("/test").
		Auth().Set("value")
	// Assert
	if builder.request.config.httpHeader.Get("Authorization") != "value" {
		t.Errorf("Authorization header not set correctly")
	}
}

func TestRequestAuthBuilder_BasicAuth(t *testing.T) {
	// Arrange
	client := DefaultClient("https://example.com")
	// Act
	builder := client.GET("/test").
		Auth().BasicAuth("username", "password")
	// Assert
	expected := "Basic " + base64.StdEncoding.EncodeToString([]byte("username:password"))
	if builder.request.config.httpHeader.Get("Authorization") != expected {
		t.Errorf(
			"BuilderHeader not set correctly, got: %s, want: %s",
			builder.request.config.httpHeader.Get("Authorization"),
			expected,
		)
	}
}

func TestRequestAuthBuilder_BearerToken(t *testing.T) {
	// Arrange
	client := DefaultClient("https://example.com")
	// Act
	builder := client.GET("/test").
		Auth().BearerToken("token")
	// Assert
	if builder.request.config.httpHeader.Get("Authorization") != "Bearer token" {
		t.Errorf(
			"BuilderHeader not set correctly, got: %s, want: %s",
			builder.request.config.httpHeader.Get("Authorization"),
			"Bearer token",
		)
	}
}
