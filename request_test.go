package fastshot

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestRequest_SetContext(t *testing.T) {
	// Arrange
	client := DefaultClient("https://example.com")
	// Act
	req := client.GET("/test").
		SetContext(nil)
	// Assert
	if req.ctx != nil {
		t.Errorf("SetContext not set correctly")
	}
}

func TestRequest_SetHeader(t *testing.T) {
	// Arrange
	client := DefaultClient("https://example.com")
	// Act
	req := client.GET("/test").
		SetHeader("key", "value")
	// Assert
	if req.httpHeader == nil || req.httpHeader.Get("key") != "value" {
		t.Errorf("SetHeader not set correctly")
	}
}

func TestRequest_SetHeaders(t *testing.T) {
	// Arrange
	client := DefaultClient("https://example.com")
	// Act
	req := client.GET("/test").
		SetHeaders(map[string]string{"key": "value"})
	// Assert
	if req.httpHeader == nil || req.httpHeader.Get("Key") != "value" {
		t.Errorf("SetHeaders not set correctly")
	}
}

func TestRequest_AddQueryParam(t *testing.T) {
	// Arrange
	client := DefaultClient("https://example.com")
	// Act
	req := client.GET("/test").
		AddQueryParam("key", "value")
	// Assert
	if req.queryParams == nil || req.queryParams["key"][0] != "value" {
		t.Errorf("AddQueryParam not set correctly")
	}
}

func TestRequest_SetQueryParam(t *testing.T) {
	// Arrange
	client := DefaultClient("https://example.com")
	// Act
	req := client.GET("/test").
		SetQueryParam("key", "value")
	// Assert
	if req.queryParams == nil || req.queryParams["key"][0] != "value" {
		t.Errorf("SetQueryParam not set correctly")
	}
}

func TestRequest_SetQueryParams(t *testing.T) {
	// Arrange
	client := DefaultClient("https://example.com")
	// Act
	req := client.GET("/test").
		SetQueryParams(map[string]string{"key": "value"})
	// Assert
	if req.queryParams == nil || req.queryParams["key"][0] != "value" {
		t.Errorf("SetQueryParams not set correctly")
	}
}

func TestRequest_SetQueryString(t *testing.T) {
	// Arrange
	client := DefaultClient("https://example.com")
	// Act
	req := client.GET("/test").
		SetQueryString("key1=value1&key2=value2")
	// Assert
	if req.queryParams.Get("key1") != "value1" || req.queryParams.Get("key2") != "value2" {
		t.Errorf("SetQueryString failed to set query parameters correctly")
	}
}

func TestRequest_SetQueryString_InvalidQuery(t *testing.T) {
	// Arrange
	client := DefaultClient("https://example.com")
	// Act
	req := client.GET("/test").
		SetQueryString("%")
	// Assert
	if len(req.validations) == 0 {
		t.Errorf("SetQueryString should append error for invalid query string")
	}
}

func TestRequest_Body(t *testing.T) {
	// Arrange
	client := DefaultClient("https://example.com")
	body := bytes.NewBuffer([]byte("test body"))
	// Act
	req := client.POST("/test").
		Body(body)
	// Assert
	if req.body == nil {
		t.Errorf("Body not set correctly")
	}
}

func TestRequest_BodyJSON(t *testing.T) {
	// Arrange
	client := DefaultClient("https://example.com")
	body := map[string]string{"key": "value"}
	// Act
	req := client.POST("/test").
		BodyJSON(body)
	// Assert
	if req.body == nil {
		t.Errorf("BodyJSON not set correctly")
	}
}

func TestRequest_BodyJSON_Error(t *testing.T) {
	// Arrange
	client := DefaultClient("https://example.com")
	body := func() {}
	// Act
	r := client.POST("/path").
		BodyJSON(body)
	// Assert
	if r.validate() == nil || !strings.Contains(r.validate().Error(), "failed to marshal JSON") {
		t.Errorf("BodyJSON didn't capture the marshaling error")
	}
}

func TestRequest_createFullURL_Error(t *testing.T) {
	// Arrange
	client := DefaultClient(":%^:")
	r := client.GET("/path")
	// Act
	_, err := r.createFullURL()
	// Assert
	if err == nil {
		t.Errorf("createFullURL did not return an error for invalid baseURL")
	}
}

func TestRequest_createFullURL_WithQueryParams(t *testing.T) {
	// Arrange
	client := DefaultClient("https://example.com")
	r := client.GET("/path").
		AddQueryParam("key1", "value1").
		AddQueryParam("key2", "value2")

	// Act
	fullURL, err := r.createFullURL()

	// Assert
	if err != nil {
		t.Errorf("createFullURL returned an error: %v", err)
		return
	}

	expectedURL := "https://example.com/path?key1=value1&key2=value2"
	if fullURL.String() != expectedURL {
		t.Errorf("createFullURL returned wrong URL, got: %s, want: %s", fullURL, expectedURL)
	}
}

func TestRequest_Send(t *testing.T) {
	// Arrange
	server := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			_ = json.NewEncoder(w).
				Encode(map[string]string{
					"message": "Success!",
				})
		}))
	defer server.Close()

	client := DefaultClient(server.URL)

	// Act
	resp, err := client.GET("/test").Send()
	if err != nil {
		t.Errorf("Execute method failed: %v", err)
		return
	}

	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(resp.RawBody())

	body, _ := io.ReadAll(resp.RawBody())

	var result map[string]string
	_ = json.Unmarshal(body, &result)

	// Assert
	if result["message"] != "Success!" {
		t.Errorf(
			"Unexpected response: got %v, want 'success'",
			result["message"],
		)
	}
}

func TestRequest_Send_Error(t *testing.T) {
	// Arrange
	client := DefaultClient("https://api.example.com")
	// Act
	_, err := client.GET("invalid path").Send()
	// Assert
	if err == nil {
		t.Errorf("Send should return error for invalid URL")
	}
}
