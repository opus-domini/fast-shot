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
	if requestBuilder.request.config.Cookies().Count() != 1 || requestBuilder.request.config.Cookies().Get(0).Name != "name" {
		t.Errorf("BuilderCookie not set correctly")
	}
}
