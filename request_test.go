package fastshot

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestRequest_SetContext(t *testing.T) {
	// Arrange
	client := DefaultClient("https://example.com")
	// Act
	req := client.GET("/test").
		SetContext(context.Background())
	// Assert
	if req.ctx == nil {
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
	if len(r.validations) != 1 || !strings.Contains(r.validations[0].Error(), ErrMsgMarshalJSON) {
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

func TestRequest_createHTTPRequest(t *testing.T) {
	tests := []struct {
		name           string
		clientBaseURL  string
		requestPath    string
		ctx            context.Context
		expectError    bool
		expectedErrMsg string
	}{
		{
			name:           "Successful Request Creation",
			clientBaseURL:  "https://example.com",
			requestPath:    "/test",
			ctx:            context.Background(),
			expectError:    false,
			expectedErrMsg: "",
		},
		{
			name:           "Error Parsing URL",
			clientBaseURL:  ":%^:",
			requestPath:    "/test",
			ctx:            context.Background(),
			expectError:    true,
			expectedErrMsg: ErrMsgParseURL,
		},
		{
			name:           "Error Creating HTTP Request",
			clientBaseURL:  "https://example.com",
			requestPath:    " ",
			ctx:            context.Background(),
			expectError:    true,
			expectedErrMsg: ErrMsgParseURL,
		},

		{
			name:           "Nil Context",
			clientBaseURL:  "https://example.com",
			requestPath:    "/test",
			ctx:            nil,
			expectError:    true,
			expectedErrMsg: "net/http: nil Context",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := NewClient(tt.clientBaseURL).
				Cookie().Add(&http.Cookie{Name: "client", Value: "test-value"}).End().
				Header().AddAccept("application/json").End().
				Build()

			httpReq, err := client.GET(tt.requestPath).
				SetContext(tt.ctx).
				AddCookie(&http.Cookie{Name: "request", Value: "test-value"}).
				SetHeader("request-header", "test-value").
				createHTTPRequest()

			if (err != nil) != tt.expectError {
				t.Errorf("createHTTPRequest() error = %v, expectError %v", err, tt.expectError)
				return
			}

			if err != nil && !strings.Contains(err.Error(), tt.expectedErrMsg) {
				t.Errorf("createHTTPRequest() error = %v, expectedErrMsg %v", err, tt.expectedErrMsg)
			}

			if err == nil && httpReq == nil {
				t.Error("Expected a non-nil http.Request object, got nil")
			}
		})
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

	tests := []struct {
		name           string
		configure      func(client *Client) *Request
		expectedResult map[string]string
		expectedError  string
	}{
		{
			name: "Successful Request",
			configure: func(client *Client) *Request {
				return client.GET("/test")
			},
			expectedResult: map[string]string{"message": "Success!"},
		},
		{
			name: "JSON Marshalling Error",
			configure: func(client *Client) *Request {
				invalidObject := func() {}
				return client.GET("/test").BodyJSON(invalidObject)
			},
			expectedError: ErrMsgValidation,
		},
		{
			name: "URL Error",
			configure: func(client *Client) *Request {
				return client.GET("%ˆ&ˆ")
			},
			expectedError: ErrMsgParseURL,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := DefaultClient(server.URL)
			req := tt.configure(client)

			resp, err := req.Send()
			if err != nil && tt.expectedError == "" {
				t.Errorf("unexpected error: %v", err)
				return
			} else if err == nil && tt.expectedError != "" {
				t.Errorf("expected error: %v, got nil", tt.expectedError)
				return
			} else if err != nil && !strings.Contains(err.Error(), tt.expectedError) {
				t.Errorf("expected error: %v, got: %v", tt.expectedError, err)
				return
			}

			if tt.expectedResult != nil {
				defer func(Body io.ReadCloser) {
					_ = Body.Close()
				}(resp.RawBody())

				body, _ := io.ReadAll(resp.RawBody())

				var result map[string]string
				_ = json.Unmarshal(body, &result)

				if result["message"] != tt.expectedResult["message"] {
					t.Errorf(
						"Unexpected response: got %v, want %v",
						result["message"], tt.expectedResult["message"],
					)
				}
			}
		})
	}
}

func TestRequest_Send_Retry(t *testing.T) {
	type retryConfig struct {
		retries       int
		retryInterval time.Duration
	}

	tests := []struct {
		name        string
		retryConfig retryConfig
		serverFunc  func() http.HandlerFunc
		expectError bool
		expectCount int
	}{
		{
			name: "Retry Successful",
			retryConfig: retryConfig{
				retries:       3,
				retryInterval: time.Millisecond,
			},
			serverFunc: func() http.HandlerFunc {
				retryCount := 0
				return func(w http.ResponseWriter, r *http.Request) {
					retryCount++
					if retryCount < 3 {
						w.WriteHeader(http.StatusInternalServerError)
						return
					}
					w.WriteHeader(http.StatusOK)
					_ = json.NewEncoder(w).Encode(map[string]string{"message": "Success!"})
				}
			},
			expectError: false,
			expectCount: 3,
		},
		{
			name: "Retry Unsuccessful",
			retryConfig: retryConfig{
				retries:       3,
				retryInterval: time.Millisecond,
			},
			serverFunc: func() http.HandlerFunc {
				retryCount := 0
				return func(w http.ResponseWriter, r *http.Request) {
					retryCount++
					w.WriteHeader(http.StatusInternalServerError)
				}
			},
			expectError: true,
			expectCount: 4, // 1 initial attempt + 3 retries
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			server := httptest.NewServer(tt.serverFunc())
			defer server.Close()

			client := DefaultClient(server.URL)

			// Act
			resp, err := client.GET("/test").SetRetry(tt.retryConfig.retries, tt.retryConfig.retryInterval).Send()
			if err != nil && !tt.expectError {
				t.Errorf("Execute method failed: %v", err)
				return
			}

			if err == nil && tt.expectError {
				t.Errorf("Expected error, but got nil")
				return
			}

			if !tt.expectError {
				defer func(Body io.ReadCloser) {
					_ = Body.Close()
				}(resp.RawBody())

				body, _ := io.ReadAll(resp.RawBody())

				var result map[string]string
				_ = json.Unmarshal(body, &result)

				// Assert
				if result["message"] != "Success!" {
					t.Errorf(
						"Unexpected response: got %v, want 'Success!'",
						result["message"],
					)
				}
			} else {
				// Assert
				if resp.RawResponse.StatusCode != http.StatusInternalServerError {
					t.Errorf(
						"Unexpected status code: got %v, want %v",
						resp.RawResponse.StatusCode, http.StatusInternalServerError,
					)
				}
			}
		})
	}
}
