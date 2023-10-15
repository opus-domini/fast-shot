package fastshot

import (
	"context"
	"encoding/json"
	"github.com/opus-domini/fast-shot/constant"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

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
		Query().AddParams(
		map[string]string{
			"key1": "value1",
			"key2": "value2",
		})

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
			expectedErrMsg: constant.ErrMsgParseURL,
		},
		{
			name:           "Error Creating HTTP Request",
			clientBaseURL:  "https://example.com",
			requestPath:    " ",
			ctx:            context.Background(),
			expectError:    true,
			expectedErrMsg: constant.ErrMsgParseURL,
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
				Cookie().Add(&http.Cookie{Name: "client", Value: "test-value"}).
				Header().AddAccept("application/json").
				Build()

			httpReq, err := client.GET(tt.requestPath).
				Context().Set(tt.ctx).
				Cookie().Add(&http.Cookie{Name: "request", Value: "test-value"}).
				Header().Set("request-header", "test-value").
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
		configure      func(client *Client) *RequestBuilder
		expectedResult map[string]string
		expectedError  string
	}{
		{
			name: "Successful Request",
			configure: func(client *Client) *RequestBuilder {
				return client.GET("/test")
			},
			expectedResult: map[string]string{"message": "Success!"},
		},
		{
			name: "JSON Marshalling Error",
			configure: func(client *Client) *RequestBuilder {
				invalidObject := func() {}
				return client.GET("/test").Body().AsJSON(invalidObject)
			},
			expectedError: constant.ErrMsgValidation,
		},
		{
			name: "URL Error",
			configure: func(client *Client) *RequestBuilder {
				return client.GET("%ˆ&ˆ")
			},
			expectedError: constant.ErrMsgParseURL,
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
