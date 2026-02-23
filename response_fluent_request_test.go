package fastshot

import (
	"net/http"
	"net/url"
	"reflect"
	"testing"
)

func TestResponseFluentRequest(t *testing.T) {
	tests := []struct {
		name           string
		request        *http.Request
		expectedMethod string
		expectedURL    string
		expectedHeader http.Header
	}{
		{
			name: "GET request",
			request: &http.Request{
				Method: "GET",
				URL:    &url.URL{Scheme: "https", Host: "api.example.com", Path: "/users"},
				Header: http.Header{"Accept": []string{"application/json"}},
			},
			expectedMethod: "GET",
			expectedURL:    "https://api.example.com/users",
			expectedHeader: http.Header{"Accept": []string{"application/json"}},
		},
		{
			name: "POST request with query parameters",
			request: &http.Request{
				Method: "POST",
				URL:    &url.URL{Scheme: "https", Host: "api.example.com", Path: "/users", RawQuery: "page=1&limit=10"},
				Header: http.Header{"Content-Type": []string{"application/json"}},
			},
			expectedMethod: "POST",
			expectedURL:    "https://api.example.com/users?page=1&limit=10",
			expectedHeader: http.Header{"Content-Type": []string{"application/json"}},
		},
		{
			name: "PUT request with multiple headers",
			request: &http.Request{
				Method: "PUT",
				URL:    &url.URL{Scheme: "https", Host: "api.example.com", Path: "/users/123"},
				Header: http.Header{
					"Content-Type":  []string{"application/json"},
					"Authorization": []string{"Bearer token123"},
				},
			},
			expectedMethod: "PUT",
			expectedURL:    "https://api.example.com/users/123",
			expectedHeader: http.Header{
				"Content-Type":  []string{"application/json"},
				"Authorization": []string{"Bearer token123"},
			},
		},
		{
			name: "DELETE request with no headers",
			request: &http.Request{
				Method: "DELETE",
				URL:    &url.URL{Scheme: "https", Host: "api.example.com", Path: "/users/123"},
				Header: http.Header{},
			},
			expectedMethod: "DELETE",
			expectedURL:    "https://api.example.com/users/123",
			expectedHeader: http.Header{},
		},
		{
			name: "PATCH request with fragment",
			request: &http.Request{
				Method: "PATCH",
				URL:    &url.URL{Scheme: "https", Host: "api.example.com", Path: "/users/123", Fragment: "section1"},
				Header: http.Header{"Content-Type": []string{"application/json-patch+json"}},
			},
			expectedMethod: "PATCH",
			expectedURL:    "https://api.example.com/users/123#section1",
			expectedHeader: http.Header{"Content-Type": []string{"application/json-patch+json"}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			response := &Response{
				request: &ResponseFluentRequest{
					request: tt.request,
				},
			}

			// Act
			result := response.Request()

			// Assert
			if result.Raw() != tt.request {
				t.Errorf("Raw() got different request")
			}
			if got := result.Method(); got != tt.expectedMethod {
				t.Errorf("Method() got %q, want %q", got, tt.expectedMethod)
			}
			if got := result.URL(); got != tt.expectedURL {
				t.Errorf("URL() got %q, want %q", got, tt.expectedURL)
			}
			if got := result.Headers(); !reflect.DeepEqual(got, tt.expectedHeader) {
				t.Errorf("Headers() got %v, want %v", got, tt.expectedHeader)
			}
		})
	}
}
