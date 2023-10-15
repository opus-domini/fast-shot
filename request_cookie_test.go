package fastshot

import (
	"net/http"
	"testing"
)

func TestRequestCookieBuilder_Add(t *testing.T) {
	// Arrange
	builder := DefaultClient("https://example.com")
	// Act
	requestBuilder := builder.GET("/test").
		Cookie().Add(&http.Cookie{Name: "name", Value: "value"})
	// Assert
	if len(requestBuilder.request.httpCookies) != 1 || requestBuilder.request.httpCookies[0].Name != "name" {
		t.Errorf("Cookie not set correctly")
	}
}
