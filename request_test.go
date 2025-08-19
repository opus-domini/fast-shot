package fastshot

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/opus-domini/fast-shot/constant"
	"github.com/opus-domini/fast-shot/constant/header"
	"github.com/opus-domini/fast-shot/constant/method"
	"github.com/opus-domini/fast-shot/mock"
	"github.com/stretchr/testify/assert"
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
		name          string
		baseURL       string
		path          string
		method        method.Type
		ctx           context.Context
		clientCookie  *http.Cookie
		requestCookie *http.Cookie
		clientHeader  map[header.Type]string
		requestHeader map[header.Type]string
		expectError   error
		setupMock     func(*mock.HttpClientComponent)
	}{
		{
			name:          "Successful Request Creation",
			baseURL:       "https://example.com",
			path:          "/test",
			method:        method.GET,
			ctx:           context.Background(),
			clientCookie:  &http.Cookie{Name: "client", Value: "test-value"},
			requestCookie: &http.Cookie{Name: "request", Value: "test-value"},
			clientHeader:  map[header.Type]string{"Accept": "application/json"},
			requestHeader: map[header.Type]string{"X-Request-ID": "123"},
		},
		{
			name:        "Invalid URL",
			baseURL:     "://invalid-url",
			path:        "/test",
			method:      method.GET,
			ctx:         context.Background(),
			expectError: errors.New("invalid URL"),
		},
		{
			name:        "Request Creation Failure",
			baseURL:     "https://example.com",
			path:        "/test",
			method:      method.Parse(":%^:"),
			ctx:         context.Background(),
			expectError: errors.New("net/http: invalid method \":%^:\""),
			setupMock:   func(m *mock.HttpClientComponent) {},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			var err error

			defer func() {
				if r := recover(); r != nil {
					err = fmt.Errorf("panic occurred: %v", r)
				}
			}()

			clientBuilder := NewClient(tt.baseURL)

			if tt.setupMock != nil {
				mockClient := new(mock.HttpClientComponent)
				tt.setupMock(mockClient)
				clientBuilder.Config().SetCustomHttpClient(mockClient)
			}

			if tt.clientCookie != nil {
				clientBuilder.Cookie().Add(tt.clientCookie)
			}

			clientBuilder.Header().AddAll(tt.clientHeader)

			rb := newRequest(clientBuilder.client, tt.method, tt.path)

			if tt.ctx != nil {
				rb.Context().Set(tt.ctx)
			}

			if tt.requestCookie != nil {
				rb.Cookie().Add(tt.requestCookie)
			}

			rb.Header().AddAll(tt.requestHeader)

			// Act
			var httpReq *http.Request
			if rb != nil {
				httpReq, err = rb.createHTTPRequest()
			}

			// Assert
			if tt.expectError != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectError, err)
				assert.Nil(t, httpReq)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, httpReq)
				assert.Equal(t, tt.method.String(), httpReq.Method)
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

func TestRequest_execute_Error(t *testing.T) {
	tests := []struct {
		name          string
		setupClient   func() ClientHttpMethods
		expectedError string
	}{
		{
			name: "Network error",
			setupClient: func() ClientHttpMethods {
				// non-existent URL to simulate a network error
				return DefaultClient("http://localhost:12345")
			},
			expectedError: "connection refused",
		},
		{
			name: "Context canceled",
			setupClient: func() ClientHttpMethods {
				server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					time.Sleep(100 * time.Millisecond)
					w.WriteHeader(http.StatusOK)
				}))
				return DefaultClient(server.URL)
			},
			expectedError: "context canceled",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			client := tt.setupClient()

			var req *RequestBuilder
			if tt.name == "Context canceled" {
				ctx, cancel := context.WithCancel(context.Background())
				req = client.GET("/").Context().Set(ctx)
				go func() {
					time.Sleep(50 * time.Millisecond)
					cancel()
				}()
			} else {
				req = client.GET("/")
			}

			// Act
			resp, err := req.Send()

			// Assert
			assert.Error(t, err)
			assert.Nil(t, resp)
			assert.Contains(t, err.Error(), tt.expectedError)
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
		configureRequest func(client ClientHttpMethods) *RequestBuilder
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
		{
			name: "Request Creation Error",
			setupClient: func() ClientHttpMethods {
				return DefaultClient(server.URL)
			},
			configureRequest: func(client ClientHttpMethods) *RequestBuilder {
				return newRequest(
					newClientConfigBase(server.URL),
					method.Parse(":%^:"),
					"/test",
				)
			},
			expectedError: constant.ErrMsgCreateRequest,
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
				assert.Nil(t, resp)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, resp)
				var result map[string]string
				assert.NoError(t, resp.Body().AsJSON(&result))
				assert.Equal(t, tt.expectedResult, result)
			}
		})
	}
}

func TestRequest_WithLoadBalancer(t *testing.T) {
	tests := []struct {
		name                 string
		numRequests          int
		expectedServer1Count int
		expectedServer2Count int
	}{
		{
			name:                 "Load balancing with 10 requests",
			numRequests:          10,
			expectedServer1Count: 5,
			expectedServer2Count: 5,
		},
		{
			name:                 "Load balancing with 11 requests",
			numRequests:          11,
			expectedServer1Count: 6,
			expectedServer2Count: 5,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
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

			// Act
			responses := make([]string, tt.numRequests)
			for i := 0; i < tt.numRequests; i++ {
				resp, err := client.GET("/test").Send()
				assert.NoError(t, err)

				body, err := resp.Body().AsString()
				assert.NoError(t, err)
				responses[i] = body
			}

			// Assert
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

			assert.Equal(t, tt.numRequests, server1Count+server2Count)
			assert.Equal(t, tt.expectedServer1Count, server1Count)
			assert.Equal(t, tt.expectedServer2Count, server2Count)
		})
	}
}

func TestRequest_Retry(t *testing.T) {
	tests := []struct {
		name             string
		setupClient      func(string) ClientHttpMethods
		request          func(ClientHttpMethods) *RequestBuilder
		serverResponses  []int
		expectedAttempts int
		expectedStatus   int
		expectedBody     string
		expectError      bool
	}{
		{
			name: "Success on third attempt (GET)",
			setupClient: func(url string) ClientHttpMethods {
				return DefaultClient(url)
			},
			request: func(client ClientHttpMethods) *RequestBuilder {
				return client.GET("/test").
					Retry().SetExponentialBackoff(10*time.Millisecond, 5, 2.0)
			},
			serverResponses: []int{
				http.StatusInternalServerError,
				http.StatusInternalServerError,
				http.StatusOK,
			},
			expectedAttempts: 3,
			expectedStatus:   http.StatusOK,
			expectedBody:     "Success",
		},
		{
			name: "Failure after max attempts (GET)",
			setupClient: func(url string) ClientHttpMethods {
				return DefaultClient(url)
			},
			request: func(client ClientHttpMethods) *RequestBuilder {
				return client.GET("/test").
					Retry().SetConstantBackoff(10*time.Millisecond, 3)
			},
			serverResponses: []int{
				http.StatusInternalServerError,
				http.StatusInternalServerError,
				http.StatusInternalServerError,
				http.StatusInternalServerError,
			},
			expectedAttempts: 3,
			expectError:      true,
		},
		{
			name: "When max attempts is 0, no retries are made (GET)",
			setupClient: func(url string) ClientHttpMethods {
				return DefaultClient(url)
			},
			request: func(client ClientHttpMethods) *RequestBuilder {
				return client.GET("/test").
					Retry().SetConstantBackoff(10*time.Millisecond, 0)
			},
			serverResponses: []int{
				http.StatusInternalServerError,
			},
			expectedAttempts: 1,
			expectedStatus:   http.StatusInternalServerError,
		},
		{
			name: "Success on second attempt (POST with body)",
			setupClient: func(url string) ClientHttpMethods {
				return DefaultClient(url)
			},
			request: func(client ClientHttpMethods) *RequestBuilder {
				return client.POST("/test").
					Body().AsJSON(map[string]string{"key": "value"}).
					Retry().SetConstantBackoffWithJitter(10*time.Millisecond, 3)
			},
			serverResponses: []int{
				http.StatusInternalServerError,
				http.StatusOK,
			},
			expectedAttempts: 2,
			expectedStatus:   http.StatusOK,
			expectedBody:     "Success",
		},
		{
			name: "Custom retry condition (retry on 404)",
			setupClient: func(url string) ClientHttpMethods {
				return DefaultClient(url)
			},
			request: func(client ClientHttpMethods) *RequestBuilder {
				return client.GET("/test").
					Retry().SetConstantBackoff(10*time.Millisecond, 3).
					Retry().WithRetryCondition(
					func(resp *Response) bool {
						return resp.Status().Code() == 404
					})
			},
			serverResponses: []int{
				http.StatusNotFound,
				http.StatusNotFound,
				http.StatusOK,
			},
			expectedAttempts: 3,
			expectedStatus:   http.StatusOK,
			expectedBody:     "Success",
		},
		{
			name: "Max delay reached using exponential backoff",
			setupClient: func(url string) ClientHttpMethods {
				return DefaultClient(url)
			},
			request: func(client ClientHttpMethods) *RequestBuilder {
				return client.GET("/test").
					Retry().SetExponentialBackoff(10*time.Millisecond, 5, 2.0).
					Retry().WithMaxDelay(15 * time.Millisecond)
			},
			serverResponses: []int{
				http.StatusInternalServerError,
				http.StatusInternalServerError,
				http.StatusInternalServerError,
				http.StatusOK,
			},
			expectedAttempts: 4,
			expectedStatus:   http.StatusOK,
			expectedBody:     "Success",
		},
		{
			name: "Max delay reached using constant backoff",
			setupClient: func(url string) ClientHttpMethods {
				return DefaultClient(url)
			},
			request: func(client ClientHttpMethods) *RequestBuilder {
				return client.GET("/test").
					Retry().SetConstantBackoff(10*time.Millisecond, 5).
					Retry().WithMaxDelay(15 * time.Millisecond)
			},
			serverResponses: []int{
				http.StatusInternalServerError,
				http.StatusInternalServerError,
				http.StatusInternalServerError,
				http.StatusOK,
			},
			expectedAttempts: 4,
			expectedStatus:   http.StatusOK,
			expectedBody:     "Success",
		},
		{
			name: "Success on first attempt with no retry",
			setupClient: func(url string) ClientHttpMethods {
				return DefaultClient(url)
			},
			request: func(client ClientHttpMethods) *RequestBuilder {
				return client.GET("/test")
			},
			serverResponses: []int{
				http.StatusOK,
			},
			expectedAttempts: 1,
			expectedStatus:   http.StatusOK,
			expectedBody:     "Success",
		},
		{
			name: "Success after 3 retries with exponential backoff and jitter",
			setupClient: func(url string) ClientHttpMethods {
				return DefaultClient(url)
			},
			request: func(client ClientHttpMethods) *RequestBuilder {
				return client.GET("/test").
					Retry().SetExponentialBackoffWithJitter(10*time.Millisecond, 5, 2.0)
			},
			serverResponses: []int{
				http.StatusInternalServerError,
				http.StatusInternalServerError,
				http.StatusInternalServerError,
				http.StatusOK,
			},
			expectedAttempts: 4,
			expectedStatus:   http.StatusOK,
			expectedBody:     "Success",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			attemptCount := 0
			requestBody := ""
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				defer func() { attemptCount++ }()
				if attemptCount >= len(tt.serverResponses) {
					w.WriteHeader(http.StatusInternalServerError)
					return
				}

				statusCode := tt.serverResponses[attemptCount]
				w.WriteHeader(statusCode)

				if r.Method == "POST" {
					body, _ := io.ReadAll(r.Body)
					requestBody = string(body)
				}

				if statusCode == http.StatusOK {
					_, _ = w.Write([]byte("Success"))
				}
			}))
			defer server.Close()

			client := tt.setupClient(server.URL)

			// Act
			resp, err := tt.request(client).Send()

			// Assert
			if tt.expectError {
				assert.Error(t, err)
				assert.Nil(t, resp)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, resp)
				assert.Equal(t, tt.expectedStatus, resp.Status().Code())
				body, _ := resp.Body().AsString()
				assert.Equal(t, tt.expectedBody, body)
			}
			assert.Equal(t, tt.expectedAttempts, attemptCount)

			if tt.request(client).request.config.Method() == method.POST {
				assert.NotEmpty(t, requestBody)
				var bodyMap map[string]string
				err := json.Unmarshal([]byte(requestBody), &bodyMap)
				assert.NoError(t, err)
				assert.Equal(t, "value", bodyMap["key"])
			}
		})
	}
}
