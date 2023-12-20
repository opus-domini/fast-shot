package fastshot

import (
	"testing"
)

func TestRequest_AddQueryParam(t *testing.T) {
	// Arrange
	client := DefaultClient("https://example.com")
	// Act
	builder := client.GET("/test").
		Query().AddParam("key", "value")
	// Assert
	if builder.request.queryParams == nil || builder.request.queryParams["key"][0] != "value" {
		t.Errorf("AddQueryParam not set correctly")
	}
}

func TestRequest_SetQueryParam(t *testing.T) {
	// Arrange
	client := DefaultClient("https://example.com")
	// Act
	builder := client.GET("/test").
		Query().SetParam("key", "value")
	// Assert
	if builder.request.queryParams == nil || builder.request.queryParams["key"][0] != "value" {
		t.Errorf("SetQueryParam not set correctly")
	}
}

func TestRequest_SetQueryParams(t *testing.T) {
	// Arrange
	client := DefaultClient("https://example.com")
	// Act
	builder := client.GET("/test").
		Query().SetParams(map[string]string{"key": "value"})
	// Assert
	if builder.request.queryParams == nil || builder.request.queryParams["key"][0] != "value" {
		t.Errorf("SetQueryParams not set correctly")
	}
}

func TestRequest_SetQueryString(t *testing.T) {
	// Arrange
	client := DefaultClient("https://example.com")
	// Act
	builder := client.GET("/test").
		Query().SetRawString("key1=value1&key2=value2")
	// Assert
	if builder.request.queryParams.Get("key1") != "value1" || builder.request.queryParams.Get("key2") != "value2" {
		t.Errorf("SetQueryString failed to set query parameters correctly")
	}
}

func TestRequest_SetQueryString_InvalidQuery(t *testing.T) {
	// Arrange
	client := DefaultClient("https://example.com")
	// Act
	builder := client.GET("/test").
		Query().SetRawString("%")
	// Assert
	if len(builder.request.validations) == 0 {
		t.Errorf("SetQueryString should append error for invalid query string")
	}
}
