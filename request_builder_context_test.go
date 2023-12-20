package fastshot

import (
	"context"
	"testing"
)

func TestRequestContextBuilder_Set(t *testing.T) {
	// Arrange
	client := DefaultClient("https://example.com")
	// Act
	builder := client.GET("/test").
		Context().Set(context.Background())
	// Assert
	if builder.request.ctx == nil {
		t.Errorf("Set Context not set correctly")
	}
}
