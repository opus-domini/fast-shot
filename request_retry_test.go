package fastshot

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

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
			resp, err := client.GET("/test").
				Retry().Set(tt.retryConfig.retries, tt.retryConfig.retryInterval).
				Send()

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
