package fastshot

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/opus-domini/fast-shot/constant"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

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
	fullURL := r.createFullURL()

	// Assert
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
			name:          "Creating HTTP Request with space in path",
			clientBaseURL: "https://example.com",
			requestPath:   " ",
			ctx:           context.Background(),
			expectError:   false,
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

// noins
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
		client         ClientHttpMethods
		configure      func(client ClientHttpMethods) *RequestBuilder
		expectedResult map[string]string
		expectedError  string
	}{
		{
			name: "Successful Request",
			configure: func(client ClientHttpMethods) *RequestBuilder {
				return client.GET("/test")
			},
			expectedResult: map[string]string{"message": "Success!"},
		},
		{
			name: "ClientConfig Proxy URL Parser Error",
			client: NewClient("https://example.com").
				Config().SetProxy(":%^:").
				Build(),
			configure: func(client ClientHttpMethods) *RequestBuilder {
				return client.GET("/test")
			},
			expectedError: constant.ErrMsgClientValidation,
		},
		{
			name:   "Request set with nil Context",
			client: NewClient("https://example.com").Build(),
			configure: func(client ClientHttpMethods) *RequestBuilder {
				//nolint:staticcheck
				return client.GET("/test").
					Context().Set(nil)
			},
			expectedError: constant.ErrMsgCreateRequest,
		},
		{
			name: "JSON Marshalling Error",
			configure: func(client ClientHttpMethods) *RequestBuilder {
				invalidObject := func() {}
				return client.GET("/test").Body().AsJSON(invalidObject)
			},
			expectedError: constant.ErrMsgRequestValidation,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.client == nil {
				tt.client = DefaultClient(server.URL)
			}

			req := tt.configure(tt.client)

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

func TestRequest_WithLoadBalancer(t *testing.T) {
	// Arrange
	server1 := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte("Server 1"))
		}))
	defer server1.Close()

	server2 := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte("Server 2"))
		}))
	defer server2.Close()

	client := NewClientLoadBalancer([]string{server1.URL, server2.URL}).Build()

	// Act
	numRequests := 10
	responses := make([]string, numRequests)
	for i := 0; i < numRequests; i++ {
		resp, err := client.GET("/test").Send()
		if err != nil {
			t.Errorf("unexpected error: %v", err)
			return
		}

		//defer func(Body io.ReadCloser) {
		//	_ = Body.Close()
		//}(resp.RawBody())

		body, _ := io.ReadAll(resp.RawBody())

		// Close the response body at the end of each iteration
		err = resp.RawBody().Close()
		if err != nil {
			t.Errorf("error closing the response body: %v", err)
		}

		responses[i] = string(body)
	}

	// Assert
	for i, response := range responses {
		expectedServer := fmt.Sprintf("Server %d", (i%2)+1)
		if response != expectedServer {
			t.Errorf("unexpected response for request %d: got %v, want %v", i+1, response, expectedServer)
		}
	}
}
