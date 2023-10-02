package fastshot

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
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
	if req.headers == nil || req.headers.Get("key") != "value" {
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
	if req.headers == nil || req.headers.Get("Key") != "value" {
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

func TestRequest_Send(t *testing.T) {
	// Arrange
	server := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			_ = json.NewEncoder(w).
				Encode(map[string]string{
					"message": "success",
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
	if result["message"] != "success" {
		t.Errorf(
			"Unexpected response: got %v, want 'success'",
			result["message"],
		)
	}
}
