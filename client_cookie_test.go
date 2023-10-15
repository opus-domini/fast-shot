package fastshot

import (
	"net/http"
	"testing"
)

func TestClientCookieBuilder_Add(t *testing.T) {
	// Arrange
	builder := NewClient("https://example.com")
	// Act
	builder.Cookie().Add(&http.Cookie{Name: "name", Value: "value"})
	// Assert
	if len(builder.client.httpCookies) != 1 || builder.client.httpCookies[0].Name != "name" {
		t.Errorf("Cookie not set correctly")
	}
}
