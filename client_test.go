package fastshot

import (
	"testing"
)

func TestClientBuilder_Build(t *testing.T) {
	// Arrange
	builder := NewClient("https://example.com")
	// Act
	client := builder.Build()
	// Assert
	if client.baseURL != "https://example.com" {
		t.Errorf("BaseURL not set correctly")
	}
}

func TestDefaultClient(t *testing.T) {
	// Arrange
	client := DefaultClient("https://api.example.com")
	// Assert
	if client.baseURL != "https://api.example.com" {
		t.Errorf("BaseURL not set correctly")
	}
}
