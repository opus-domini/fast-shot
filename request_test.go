package fastshot

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	"github.com/opus-domini/fast-shot/constant"
	"github.com/opus-domini/fast-shot/constant/header"
	"github.com/opus-domini/fast-shot/constant/method"
	"github.com/opus-domini/fast-shot/mock"
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
			if got := fullURL.String(); got != tt.expectedURLStr {
				t.Errorf("got %q, want %q", got, tt.expectedURLStr)
			}
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
				mockClient := &mock.HttpClientComponent{}
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
				if err == nil {
					t.Error("expected error, got nil")
				} else if err.Error() != tt.expectError.Error() {
					t.Errorf("error got %q, want %q", err.Error(), tt.expectError.Error())
				}
				if httpReq != nil {
					t.Errorf("httpReq got %v, want nil", httpReq)
				}
			} else {
				if err != nil {
					t.Fatalf("unexpected error: %v", err)
				}
				if httpReq == nil {
					t.Fatal("httpReq got nil, want non-nil")
				}
				if got := httpReq.Method; got != tt.method.String() {
					t.Errorf("Method got %q, want %q", got, tt.method.String())
				}
				if got := httpReq.URL.String(); got != tt.baseURL+tt.path {
					t.Errorf("URL got %q, want %q", got, tt.baseURL+tt.path)
				}

				// Check cookies
				cookies := httpReq.Cookies()
				cookieMap := make(map[string]string)
				for _, cookie := range cookies {
					cookieMap[cookie.Name] = cookie.Value
				}

				if tt.clientCookie != nil {
					if got := cookieMap[tt.clientCookie.Name]; got != tt.clientCookie.Value {
						t.Errorf("client cookie %q got %q, want %q", tt.clientCookie.Name, got, tt.clientCookie.Value)
					}
				}
				if tt.requestCookie != nil {
					if got := cookieMap[tt.requestCookie.Name]; got != tt.requestCookie.Value {
						t.Errorf("request cookie %q got %q, want %q", tt.requestCookie.Name, got, tt.requestCookie.Value)
					}
				}

				// Check headers
				for k, v := range tt.requestHeader {
					if got := httpReq.Header.Get(k.String()); got != v {
						t.Errorf("request header %q got %q, want %q", k, got, v)
					}
				}
				for k, v := range tt.clientHeader {
					if got := httpReq.Header.Get(k.String()); got != v {
						t.Errorf("client header %q got %q, want %q", k, got, v)
					}
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
			if err == nil {
				t.Error("expected error, got nil")
			}
			if resp != nil {
				t.Errorf("resp got %v, want nil", resp)
			}
			if !strings.Contains(err.Error(), tt.expectedError) {
				t.Errorf("error %q does not contain %q", err.Error(), tt.expectedError)
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
				if err == nil {
					t.Error("expected error, got nil")
				} else if !strings.Contains(err.Error(), tt.expectedError) {
					t.Errorf("error %q does not contain %q", err.Error(), tt.expectedError)
				}
				if resp != nil {
					t.Errorf("resp got %v, want nil", resp)
				}
			} else {
				if err != nil {
					t.Fatalf("unexpected error: %v", err)
				}
				if resp == nil {
					t.Fatal("resp got nil, want non-nil")
				}
				var result map[string]string
				if err := resp.Body().AsJSON(&result); err != nil {
					t.Fatalf("unexpected error reading body: %v", err)
				}
				if !reflect.DeepEqual(result, tt.expectedResult) {
					t.Errorf("body got %v, want %v", result, tt.expectedResult)
				}
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
			for i := range tt.numRequests {
				resp, err := client.GET("/test").Send()
				if err != nil {
					t.Fatalf("unexpected error: %v", err)
				}

				body, err := resp.Body().AsString()
				if err != nil {
					t.Fatalf("unexpected error reading body: %v", err)
				}
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

			if got := server1Count + server2Count; got != tt.numRequests {
				t.Errorf("total requests got %d, want %d", got, tt.numRequests)
			}
			if server1Count != tt.expectedServer1Count {
				t.Errorf("server1 count got %d, want %d", server1Count, tt.expectedServer1Count)
			}
			if server2Count != tt.expectedServer2Count {
				t.Errorf("server2 count got %d, want %d", server2Count, tt.expectedServer2Count)
			}
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
			name:        "Success on third attempt (GET)",
			setupClient: DefaultClient,
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
			name:        "Failure after max attempts (GET)",
			setupClient: DefaultClient,
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
			name:        "When max attempts is 0, no retries are made (GET)",
			setupClient: DefaultClient,
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
			name:        "Success on second attempt (POST with body)",
			setupClient: DefaultClient,
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
			name:        "Custom retry condition (retry on 404)",
			setupClient: DefaultClient,
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
			name:        "Max delay reached using exponential backoff",
			setupClient: DefaultClient,
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
			name:        "Max delay reached using constant backoff",
			setupClient: DefaultClient,
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
			name:        "Success on first attempt with no retry",
			setupClient: DefaultClient,
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
			name:        "Success after 3 retries with exponential backoff and jitter",
			setupClient: DefaultClient,
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
				if err == nil {
					t.Error("expected error, got nil")
				}
				if resp != nil {
					t.Errorf("resp got %v, want nil", resp)
				}
			} else {
				if err != nil {
					t.Fatalf("unexpected error: %v", err)
				}
				if resp == nil {
					t.Fatal("resp got nil, want non-nil")
				}
				if got := resp.Status().Code(); got != tt.expectedStatus {
					t.Errorf("status got %d, want %d", got, tt.expectedStatus)
				}
				body, _ := resp.Body().AsString()
				if body != tt.expectedBody {
					t.Errorf("body got %q, want %q", body, tt.expectedBody)
				}
			}
			if attemptCount != tt.expectedAttempts {
				t.Errorf("attempts got %d, want %d", attemptCount, tt.expectedAttempts)
			}

			if tt.request(client).request.config.Method() == method.POST {
				if requestBody == "" {
					t.Error("request body is empty, want non-empty")
				}
				var bodyMap map[string]string
				if err := json.Unmarshal([]byte(requestBody), &bodyMap); err != nil {
					t.Fatalf("unexpected error unmarshalling body: %v", err)
				}
				if bodyMap["key"] != "value" {
					t.Errorf("body[key] got %q, want %q", bodyMap["key"], "value")
				}
			}
		})
	}
}

func TestRequest_Hooks(t *testing.T) {
	t.Run("BeforeRequest hook modifies header", func(t *testing.T) {
		// Arrange
		var receivedHeader string
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			receivedHeader = r.Header.Get("X-Hook-Header")
			w.WriteHeader(http.StatusOK)
		}))
		defer server.Close()

		client := NewClient(server.URL).
			Hook().OnBeforeRequest(func(req *http.Request) error {
			req.Header.Set("X-Hook-Header", "injected-value")
			return nil
		}).
			Build()

		// Act
		_, err := client.GET("/test").Send()

		// Assert
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if receivedHeader != "injected-value" {
			t.Errorf("header got %q, want %q", receivedHeader, "injected-value")
		}
	})

	t.Run("BeforeRequest hook error aborts request", func(t *testing.T) {
		// Arrange
		serverCalled := false
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			serverCalled = true
			w.WriteHeader(http.StatusOK)
		}))
		defer server.Close()

		hookErr := errors.New("hook rejected request")
		client := NewClient(server.URL).
			Hook().OnBeforeRequest(func(req *http.Request) error {
			return hookErr
		}).
			Build()

		// Act
		resp, err := client.GET("/test").Send()

		// Assert
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		if !strings.Contains(err.Error(), constant.ErrMsgBeforeRequestHook) {
			t.Errorf("error %q does not contain %q", err.Error(), constant.ErrMsgBeforeRequestHook)
		}
		if !strings.Contains(err.Error(), hookErr.Error()) {
			t.Errorf("error %q does not contain %q", err.Error(), hookErr.Error())
		}
		if resp != nil {
			t.Errorf("resp got %v, want nil", resp)
		}
		if serverCalled {
			t.Error("server was called, want not called")
		}
	})

	t.Run("AfterResponse hook records status code", func(t *testing.T) {
		// Arrange
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusTeapot)
		}))
		defer server.Close()

		var capturedStatus int
		client := NewClient(server.URL).
			Hook().OnAfterResponse(func(req *http.Request, resp *http.Response) {
			capturedStatus = resp.StatusCode
		}).
			Build()

		// Act
		_, err := client.GET("/test").Send()

		// Assert
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if capturedStatus != http.StatusTeapot {
			t.Errorf("captured status got %d, want %d", capturedStatus, http.StatusTeapot)
		}
	})

	t.Run("Client hooks run before request hooks", func(t *testing.T) {
		// Arrange
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		}))
		defer server.Close()

		var order []string
		client := NewClient(server.URL).
			Hook().OnBeforeRequest(func(req *http.Request) error {
			order = append(order, "client")
			return nil
		}).
			Build()

		// Act
		_, err := client.GET("/test").
			Hook().OnBeforeRequest(func(req *http.Request) error {
			order = append(order, "request")
			return nil
		}).
			Send()

		// Assert
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(order) != 2 {
			t.Fatalf("order length got %d, want 2", len(order))
		}
		if order[0] != "client" {
			t.Errorf("order[0] got %q, want %q", order[0], "client")
		}
		if order[1] != "request" {
			t.Errorf("order[1] got %q, want %q", order[1], "request")
		}
	})

	t.Run("Multiple hooks execute in registration order", func(t *testing.T) {
		// Arrange
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		}))
		defer server.Close()

		var order []int
		client := NewClient(server.URL).
			Hook().OnBeforeRequest(func(req *http.Request) error {
			order = append(order, 1)
			return nil
		}).
			Hook().OnBeforeRequest(func(req *http.Request) error {
			order = append(order, 2)
			return nil
		}).
			Hook().OnBeforeRequest(func(req *http.Request) error {
			order = append(order, 3)
			return nil
		}).
			Build()

		// Act
		_, err := client.GET("/test").Send()

		// Assert
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		expected := []int{1, 2, 3}
		if !reflect.DeepEqual(order, expected) {
			t.Errorf("order got %v, want %v", order, expected)
		}
	})

	t.Run("Hooks run on each retry attempt", func(t *testing.T) {
		// Arrange
		attemptCount := 0
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			attemptCount++
			if attemptCount < 3 {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			w.WriteHeader(http.StatusOK)
		}))
		defer server.Close()

		var beforeCount atomic.Int32
		var afterCount atomic.Int32
		client := NewClient(server.URL).
			Hook().OnBeforeRequest(func(req *http.Request) error {
			beforeCount.Add(1)
			return nil
		}).
			Hook().OnAfterResponse(func(req *http.Request, resp *http.Response) {
			afterCount.Add(1)
		}).
			Build()

		// Act
		_, err := client.GET("/test").
			Retry().SetConstantBackoff(10*time.Millisecond, 3).
			Send()

		// Assert
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if got := beforeCount.Load(); got != 3 {
			t.Errorf("before hook count got %d, want 3", got)
		}
		if got := afterCount.Load(); got != 3 {
			t.Errorf("after hook count got %d, want 3", got)
		}
	})

	t.Run("AfterResponse does not run when BeforeRequest aborts", func(t *testing.T) {
		// Arrange
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		}))
		defer server.Close()

		afterCalled := false
		client := NewClient(server.URL).
			Hook().OnBeforeRequest(func(req *http.Request) error {
			return errors.New("abort")
		}).
			Hook().OnAfterResponse(func(req *http.Request, resp *http.Response) {
			afterCalled = true
		}).
			Build()

		// Act
		_, _ = client.GET("/test").Send()

		// Assert
		if afterCalled {
			t.Error("after response hook was called, want not called")
		}
	})

	t.Run("Request-level hook on single request only", func(t *testing.T) {
		// Arrange
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		}))
		defer server.Close()

		var hookCount atomic.Int32
		client := NewClient(server.URL).Build()

		// Act - first request with hook
		_, err := client.GET("/test").
			Hook().OnBeforeRequest(func(req *http.Request) error {
			hookCount.Add(1)
			return nil
		}).
			Send()
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		// Act - second request without hook
		_, err = client.GET("/test").Send()
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		// Assert
		if got := hookCount.Load(); got != 1 {
			t.Errorf("hook count got %d, want 1", got)
		}
	})
}
