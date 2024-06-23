package fastshot

import (
	"strings"
	"testing"
)

func TestRequestHeaderBuilder_Add(t *testing.T) {
	// Arrange
	builder := DefaultClient("https://example.com")
	// Act
	requestBuilder := builder.GET("/test").
		Header().Add("key", "value").
		Header().Add("key", "value2")
	// Assert
	if !strings.Contains(requestBuilder.request.config.httpHeader.Get("key"), "value") {
		t.Errorf("BuilderHeader not set correctly")
	}
}

func TestRequestHeaderBuilder_AddAll(t *testing.T) {
	// Arrange
	builder := DefaultClient("https://example.com")
	// Act
	requestBuilder := builder.GET("/test").
		Header().AddAll(map[string]string{"key1": "value1", "key2": "value2"})
	// Assert
	if !strings.Contains(requestBuilder.request.config.httpHeader.Get("key2"), "value2") {
		t.Errorf("BuilderHeader not set correctly")
	}
}

func TestRequestHeaderBuilder_Set(t *testing.T) {
	// Arrange
	builder := DefaultClient("https://example.com")
	// Act
	requestBuilder := builder.GET("/test").
		Header().Set("key", "value").
		Header().Set("key", "value2")

	// Assert
	if requestBuilder.request.config.httpHeader.Get("key") != "value2" {
		t.Errorf("BuilderHeader not set correctly")
	}
}

func TestRequestHeaderBuilder_SetAll(t *testing.T) {
	// Arrange
	builder := DefaultClient("https://example.com")
	// Act
	requestBuilder := builder.GET("/test").
		Header().SetAll(map[string]string{"key1": "value1", "key2": "value2"})
	// Assert
	if !strings.Contains(requestBuilder.request.config.httpHeader.Get("key2"), "value2") {
		t.Errorf("BuilderHeader not set correctly")
	}
}

func TestRequestHeaderBuilder_AddAccept(t *testing.T) {
	// Arrange
	builder := DefaultClient("https://example.com")
	valueToFind := "application/xml"
	// Act
	headerBuilder := builder.GET("/test").
		Header().AddAccept("application/json").
		Header().AddAccept(valueToFind)
	// Assert
	values := headerBuilder.request.config.Header().Unwrap().Values("Accept")
	valueFound := false
	for _, value := range values {
		if value == valueToFind {
			valueFound = true
			break
		}
	}
	if !valueFound {
		t.Errorf("BuilderHeader not set correctly")
	}
}

func TestRequestHeaderBuilder_AddUserAgent(t *testing.T) {
	// Arrange
	builder := DefaultClient("https://example.com")
	valueToFind := "chrome"
	// Act
	headerBuilder := builder.GET("/test/").
		Header().AddUserAgent("mobile").
		Header().AddUserAgent(valueToFind).
		Header().AddUserAgent("firefox")
	// Assert
	values := headerBuilder.request.config.Header().Unwrap().Values("User-Agent")
	valueFound := false
	for _, value := range values {
		if value == valueToFind {
			valueFound = true
			break
		}
	}
	if !valueFound {
		t.Errorf("BuilderHeader not set correctly")
	}
}

func TestRequestHeaderBuilder_AddContentType(t *testing.T) {
	// Arrange
	builder := DefaultClient("https://example.com")
	valueToFind := "multipart/form-data; boundary=something"
	// Act
	requestBuilder := builder.GET("/test").
		Header().AddContentType("text/html; charset=utf-8").
		Header().AddContentType(valueToFind)
	// Assert
	values := requestBuilder.request.config.Header().Unwrap().Values("Content-Type")
	valueFound := false
	for _, value := range values {
		if value == valueToFind {
			valueFound = true
			break
		}
	}
	if !valueFound {
		t.Errorf("BuilderHeader not set correctly")
	}
}
