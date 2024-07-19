package fastshot

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func logServerResponse(responseStatus string) {
	fmt.Printf("timestamp: %s, response: %s\n", time.Now().Format(time.StampMilli), responseStatus)
}

func TestRequest_Send_Retry(t *testing.T) {
	type retryConfig struct {
		retryBuilder func(requestBuilder *RequestBuilder) *RequestBuilder
		maxAttempts  uint
		interval     time.Duration
	}

	tests := []struct {
		name        string
		retryConfig retryConfig
		serverFunc  func() (http.HandlerFunc, *int)
		expectError bool
		expectCount int
	}{
		{
			name: "Retry Constant Backoff Successful",
			retryConfig: retryConfig{
				retryBuilder: func(requestBuilder *RequestBuilder) *RequestBuilder {
					return requestBuilder.
						Retry().SetConstantBackoff(time.Millisecond, 3)
				},
				interval:    time.Millisecond,
				maxAttempts: 3,
			},
			serverFunc: func() (http.HandlerFunc, *int) {
				retryCount := 0
				return func(w http.ResponseWriter, r *http.Request) {
					retryCount++
					if retryCount < 3 {
						logServerResponse("500 SERVER ERROR")
						w.WriteHeader(http.StatusInternalServerError)
						return
					}
					logServerResponse("200 OK")
					w.WriteHeader(http.StatusOK)
					_ = json.NewEncoder(w).Encode(map[string]string{"message": "Success!"})
				}, &retryCount
			},
			expectError: false,
			expectCount: 3,
		},
		{
			name: "Retry Constant Backoff Unsuccessful",
			retryConfig: retryConfig{
				retryBuilder: func(requestBuilder *RequestBuilder) *RequestBuilder {
					return requestBuilder.
						Retry().SetConstantBackoff(time.Millisecond, 3)
				},
				maxAttempts: 3,
				interval:    time.Millisecond,
			},
			serverFunc: func() (http.HandlerFunc, *int) {
				retryCount := 0
				return func(w http.ResponseWriter, r *http.Request) {
					retryCount++
					logServerResponse("500 SERVER ERROR")
					w.WriteHeader(http.StatusInternalServerError)
				}, &retryCount
			},
			expectError: true,
			expectCount: 3,
		},
		{
			name: "Retry Constant Backoff With Jitter Successful",
			retryConfig: retryConfig{
				retryBuilder: func(requestBuilder *RequestBuilder) *RequestBuilder {
					return requestBuilder.
						Retry().SetConstantBackoffWithJitter(time.Millisecond, 3)
				},
				interval:    time.Millisecond,
				maxAttempts: 3,
			},
			serverFunc: func() (http.HandlerFunc, *int) {
				retryCount := 0
				return func(w http.ResponseWriter, r *http.Request) {
					retryCount++
					if retryCount < 3 {
						logServerResponse("500 SERVER ERROR")
						w.WriteHeader(http.StatusInternalServerError)
						return
					}
					logServerResponse("200 OK")
					w.WriteHeader(http.StatusOK)
					_ = json.NewEncoder(w).Encode(map[string]string{"message": "Success!"})
				}, &retryCount
			},
			expectError: false,
			expectCount: 3,
		},
		{
			name: "Retry Exponential Backoff Successful",
			retryConfig: retryConfig{
				retryBuilder: func(requestBuilder *RequestBuilder) *RequestBuilder {
					return requestBuilder.
						Retry().SetExponentialBackoff(5*time.Millisecond, 5, 2)
				},
			},
			serverFunc: func() (http.HandlerFunc, *int) {
				retryCount := 0
				return func(w http.ResponseWriter, r *http.Request) {
					retryCount++
					if retryCount < 5 {
						logServerResponse("500 SERVER ERROR")
						w.WriteHeader(http.StatusInternalServerError)
						return
					}
					logServerResponse("200 OK")
					w.WriteHeader(http.StatusOK)
					_ = json.NewEncoder(w).Encode(map[string]string{"message": "Success!"})
				}, &retryCount
			},
			expectError: false,
			expectCount: 5,
		},
		{
			name: "Retry Exponential Backoff Unsuccessful",
			retryConfig: retryConfig{
				retryBuilder: func(requestBuilder *RequestBuilder) *RequestBuilder {
					return requestBuilder.
						Retry().SetExponentialBackoff(5*time.Millisecond, 5, 2)
				},
			},
			serverFunc: func() (http.HandlerFunc, *int) {
				retryCount := 0
				return func(w http.ResponseWriter, r *http.Request) {
					retryCount++
					logServerResponse("500 SERVER ERROR")
					w.WriteHeader(http.StatusInternalServerError)
				}, &retryCount
			},
			expectError: true,
			expectCount: 5,
		},
		{
			name: "Retry Exponential Backoff With Jitter Successful",
			retryConfig: retryConfig{
				retryBuilder: func(requestBuilder *RequestBuilder) *RequestBuilder {
					return requestBuilder.
						Retry().SetExponentialBackoffWithJitter(5*time.Millisecond, 5, 2)
				},
			},
			serverFunc: func() (http.HandlerFunc, *int) {
				retryCount := 0
				return func(w http.ResponseWriter, r *http.Request) {
					retryCount++
					if retryCount < 5 {
						logServerResponse("500 SERVER ERROR")
						w.WriteHeader(http.StatusInternalServerError)
						return
					}
					logServerResponse("200 OK")
					w.WriteHeader(http.StatusOK)
					_ = json.NewEncoder(w).Encode(map[string]string{"message": "Success!"})
				}, &retryCount
			},
			expectError: false,
			expectCount: 5,
		},
		{
			name: "Retry Exponential Backoff With Max Delay Successful",
			retryConfig: retryConfig{
				retryBuilder: func(requestBuilder *RequestBuilder) *RequestBuilder {
					return requestBuilder.
						Retry().SetExponentialBackoff(1*time.Minute, 5, 2).
						Retry().WithMaxDelay(10 * time.Millisecond)
				},
			},
			serverFunc: func() (http.HandlerFunc, *int) {
				retryCount := 0
				return func(w http.ResponseWriter, r *http.Request) {
					retryCount++
					if retryCount < 5 {
						logServerResponse("500 SERVER ERROR")
						w.WriteHeader(http.StatusInternalServerError)
						return
					}
					logServerResponse("200 OK")
					w.WriteHeader(http.StatusOK)
					_ = json.NewEncoder(w).Encode(map[string]string{"message": "Success!"})
				}, &retryCount
			},
			expectError: false,
			expectCount: 5,
		},
		{
			name: "Retry Exponential Backoff With Custom Retry Condition Successful",
			retryConfig: retryConfig{
				retryBuilder: func(requestBuilder *RequestBuilder) *RequestBuilder {
					return requestBuilder.
						Retry().SetExponentialBackoff(1*time.Millisecond, 5, 2).
						Retry().
						WithRetryCondition(
							func(response Response) bool {
								return response.IsError() || response.RawResponse.StatusCode == http.StatusNoContent
							},
						)
				},
			},
			serverFunc: func() (http.HandlerFunc, *int) {
				retryCount := 0
				return func(w http.ResponseWriter, r *http.Request) {
					retryCount++
					if retryCount < 5 {
						logServerResponse("204 NO CONTENT")
						w.WriteHeader(http.StatusNoContent)
						return
					}
					logServerResponse("200 OK")
					w.WriteHeader(http.StatusOK)
					_ = json.NewEncoder(w).Encode(map[string]string{"message": "Success!"})
				}, &retryCount
			},
			expectError: false,
			expectCount: 5,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			serverFunc, retryCount := tt.serverFunc()
			server := httptest.NewServer(serverFunc)
			defer server.Close()

			requestBuilder := DefaultClient(server.URL).GET("/test")

			// Act
			resp, err := tt.retryConfig.retryBuilder(requestBuilder).Send()

			// Assert
			if *retryCount != tt.expectCount {
				t.Errorf("Unexpected retry count: got %v, want %v", *retryCount, tt.expectCount)
			}

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
