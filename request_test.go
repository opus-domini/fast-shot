package fastshot

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/opus-domini/fast-shot/constant/header"
	"github.com/stretchr/testify/assert"

	"github.com/opus-domini/fast-shot/constant"
)

func TestRequest_createFullURL(t *testing.T) {
	tests := []struct {
		name           string
		baseURL        string
		path           string
		queryParams    map[string]string
		expectedURLStr string
	}{
		{
			name:           "Base URL with path and query params",
			baseURL:        "https://example.com",
			path:           "/path",
			queryParams:    map[string]string{"key1": "value1", "key2": "value2"},
			expectedURLStr: "https://example.com/path?key1=value1&key2=value2",
		},
		{
			name:           "Base URL with path and no query params",
			baseURL:        "https://example.com",
			path:           "/path",
			queryParams:    map[string]string{},
			expectedURLStr: "https://example.com/path",
		},
		{
			name:           "Base URL with no path and query params",
			baseURL:        "https://example.com",
			path:           "",
			queryParams:    map[string]string{"key": "value"},
			expectedURLStr: "https://example.com?key=value",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			client := DefaultClient(tt.baseURL)
			rb := client.GET(tt.path)
			for k, v := range tt.queryParams {
				rb.Query().AddParam(k, v)
			}

			// Act
			fullURL := rb.createFullURL()

			// Assert
			assert.Equal(t, tt.expectedURLStr, fullURL.String())
		})
	}
}

func TestRequest_createHTTPRequest(t *testing.T) {
	tests := []struct {
		name           string
		baseURL        string
		path           string
		method         string
		ctx            context.Context
		clientCookie   *http.Cookie
		requestCookie  *http.Cookie
		clientHeader   map[header.Type]string
		requestHeader  map[header.Type]string
		expectError    bool
		expectedErrMsg string
	}{
		{
			name:          "Successful Request Creation",
			baseURL:       "https://example.com",
			path:          "/test",
			method:        "GET",
			ctx:           context.Background(),
			clientCookie:  &http.Cookie{Name: "client", Value: "test-value"},
			requestCookie: &http.Cookie{Name: "request", Value: "test-value"},
			clientHeader:  map[header.Type]string{"Accept": "application/json"},
			requestHeader: map[header.Type]string{"X-Request-ID": "123"},
			expectError:   false,
		},
		{
			name:           "Invalid URL",
			baseURL:        "://invalid-url",
			path:           "/test",
			method:         "GET",
			ctx:            context.Background(),
			expectError:    true,
			expectedErrMsg: "invalid URL",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			var client ClientHttpMethods
			var err error

			defer func() {
				if r := recover(); r != nil {
					err = fmt.Errorf("panic occurred: %v", r)
				}
			}()

			clientBuilder := NewClient(tt.baseURL)

			if tt.clientCookie != nil {
				clientBuilder.Cookie().Add(tt.clientCookie)
			}

			for k, v := range tt.clientHeader {
				clientBuilder.Header().Add(k, v)
			}

			client = clientBuilder.Build()

			rb := client.GET(tt.path)

			if tt.ctx != nil {
				rb.Context().Set(tt.ctx)
			}

			if tt.requestCookie != nil {
				rb.Cookie().Add(tt.requestCookie)
			}

			for k, v := range tt.requestHeader {
				rb.Header().Set(k, v)
			}

			// Act
			httpReq, err := rb.createHTTPRequest()

			// Assert
			if tt.expectError {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedErrMsg)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, httpReq)
				assert.Equal(t, tt.method, httpReq.Method)
				assert.Equal(t, tt.baseURL+tt.path, httpReq.URL.String())

				// Check cookies
				cookies := httpReq.Cookies()
				cookieMap := make(map[string]string)
				for _, cookie := range cookies {
					cookieMap[cookie.Name] = cookie.Value
				}

				if tt.clientCookie != nil {
					assert.Equal(t, tt.clientCookie.Value, cookieMap[tt.clientCookie.Name])
				}
				if tt.requestCookie != nil {
					assert.Equal(t, tt.requestCookie.Value, cookieMap[tt.requestCookie.Name])
				}

				// Check headers
				for k, v := range tt.requestHeader {
					assert.Equal(t, v, httpReq.Header.Get(k.String()))
				}
				for k, v := range tt.clientHeader {
					assert.Equal(t, v, httpReq.Header.Get(k.String()))
				}
			}
		})
	}
}

func TestRequest_Send(t *testing.T) {
	// Arrange
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(map[string]string{"message": "Success!"})
	}))
	defer server.Close()

	tests := []struct {
		name             string
		setupClient      func() ClientHttpMethods
		configureRequest func(ClientHttpMethods) *RequestBuilder
		expectedResult   map[string]string
		expectedError    string
	}{
		{
			name: "Successful Request",
			setupClient: func() ClientHttpMethods {
				return DefaultClient(server.URL)
			},
			configureRequest: func(client ClientHttpMethods) *RequestBuilder {
				return client.GET("/test")
			},
			expectedResult: map[string]string{"message": "Success!"},
		},
		{
			name: "ClientConfig Proxy URL Parser Error",
			setupClient: func() ClientHttpMethods {
				return NewClient("https://example.com").
					Config().SetProxy(":%^:").
					Build()
			},
			configureRequest: func(client ClientHttpMethods) *RequestBuilder {
				return client.GET("/test")
			},
			expectedError: constant.ErrMsgClientValidation,
		},
		{
			name: "JSON Marshalling Error",
			setupClient: func() ClientHttpMethods {
				return DefaultClient(server.URL)
			},
			configureRequest: func(client ClientHttpMethods) *RequestBuilder {
				return client.GET("/test").Body().AsJSON(func() {})
			},
			expectedError: constant.ErrMsgRequestValidation,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			client := tt.setupClient()
			req := tt.configureRequest(client)

			// Act
			resp, err := req.Send()

			// Assert
			if tt.expectedError != "" {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedError)
			} else {
				assert.NoError(t, err)
				var result map[string]string
				assert.NoError(t, resp.Body().AsJSON(&result))
				assert.Equal(t, tt.expectedResult, result)
			}
		})
	}
}

func TestRequest_WithLoadBalancer(t *testing.T) {
	// Arrange
	server1 := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("Server 1"))
	}))
	defer server1.Close()

	server2 := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("Server 2"))
	}))
	defer server2.Close()

	client := NewClientLoadBalancer([]string{server1.URL, server2.URL}).Build()

	// Act & Assert
	numRequests := 10
	responses := make([]string, numRequests)
	for i := 0; i < numRequests; i++ {
		resp, err := client.GET("/test").Send()
		assert.NoError(t, err)

		body, err := resp.Body().AsString()
		assert.NoError(t, err)
		responses[i] = body
	}

	// Assert load balancing
	server1Count := 0
	server2Count := 0
	for _, response := range responses {
		switch response {
		case "Server 1":
			server1Count++
		case "Server 2":
			server2Count++
		default:
			t.Errorf("Unexpected response: %s", response)
		}
	}

	assert.Equal(t, numRequests, server1Count+server2Count)
	assert.True(t, server1Count > 0 && server2Count > 0, "Load balancing not working as expected")
}

func TestRequest_Retry(t *testing.T) {
	// Arrange
	attemptCount := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		attemptCount++
		if attemptCount < 3 {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("Success"))
	}))
	defer server.Close()

	client := DefaultClient(server.URL)

	// Act
	resp, err := client.GET("/test").
		Retry().SetExponentialBackoff(10*time.Millisecond, 5, 2.0).
		Send()

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.Status().Code())
	body, err := resp.Body().AsString()
	assert.NoError(t, err)
	assert.Equal(t, "Success", body)
	assert.Equal(t, 3, attemptCount)
}
