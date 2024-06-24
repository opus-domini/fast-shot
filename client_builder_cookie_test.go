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
	if builder.client.Cookies().Count() != 1 || builder.client.Cookies().Get(0).Name != "name" {
		t.Errorf("BuilderCookie not set correctly")
	}
}
