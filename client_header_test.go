package fastshot

import (
	"strings"
	"testing"
)

func TestClientHeaderBuilder_Add(t *testing.T) {
	// Arrange
	builder := NewClient("https://example.com")
	// Act
	builder.Header().Add("key", "value").
		Header().Add("key", "value2")
	// Assert
	if !strings.Contains(builder.client.httpHeader.Get("key"), "value") {
		t.Errorf("Header not set correctly")
	}
}

func TestClientHeaderBuilder_AddAll(t *testing.T) {
	// Arrange
	builder := NewClient("https://example.com")
	// Act
	builder.Header().AddAll(map[string]string{"key1": "value1", "key2": "value2"})
	// Assert
	if !strings.Contains(builder.client.httpHeader.Get("key2"), "value2") {
		t.Errorf("Header not set correctly")
	}
}

func TestClientHeaderBuilder_Set(t *testing.T) {
	// Arrange
	builder := NewClient("https://example.com")
	// Act
	builder.Header().Set("key", "value").
		Header().Set("key", "value2")

	// Assert
	if builder.client.httpHeader.Get("key") != "value2" {
		t.Errorf("Header not set correctly")
	}
}

func TestClientHeaderBuilder_SetAll(t *testing.T) {
	// Arrange
	builder := NewClient("https://example.com")
	// Act
	builder.Header().SetAll(map[string]string{"key1": "value1", "key2": "value2"})
	// Assert
	if !strings.Contains(builder.client.httpHeader.Get("key2"), "value2") {
		t.Errorf("Header not set correctly")
	}
}

func TestClientHeaderBuilder_AddAccept(t *testing.T) {
	// Arrange
	builder := NewClient("https://example.com")
	valueToFind := "application/xml"
	// Act
	headerBuilder := builder.Header()
	headerBuilder.AddAccept("application/json")
	headerBuilder.AddAccept(valueToFind)
	// Assert
	values := builder.client.httpHeader.Values("Accept")
	valueFound := false
	for _, value := range values {
		if value == valueToFind {
			valueFound = true
			break
		}
	}
	if !valueFound {
		t.Errorf("Header not set correctly")
	}
}

func TestClientHeaderBuilder_AddUserAgent(t *testing.T) {
	// Arrange
	builder := NewClient("https://example.com")
	valueToFind := "chrome"
	// Act
	headerBuilder := builder.Header()
	headerBuilder.AddUserAgent("mobile")
	headerBuilder.AddUserAgent(valueToFind)
	headerBuilder.AddUserAgent("firefox")
	// Assert
	values := builder.client.httpHeader.Values("User-Agent")
	valueFound := false
	for _, value := range values {
		if value == valueToFind {
			valueFound = true
			break
		}
	}
	if !valueFound {
		t.Errorf("Header not set correctly")
	}
}

func TestClientHeaderBuilder_AddContentType(t *testing.T) {
	// Arrange
	builder := NewClient("https://example.com")
	valueToFind := "multipart/form-data; boundary=something"
	// Act
	builder.Header().AddContentType("text/html; charset=utf-8").
		Header().AddContentType(valueToFind)
	// Assert
	values := builder.client.httpHeader.Values("Content-Type")
	valueFound := false
	for _, value := range values {
		if value == valueToFind {
			valueFound = true
			break
		}
	}
	if !valueFound {
		t.Errorf("Header not set correctly")
	}
}
