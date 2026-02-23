package fastshot

import (
	"net/http"
	"reflect"
	"sort"
	"testing"
)

func TestResponseFluentHeader(t *testing.T) {
	tests := []struct {
		name           string
		headers        http.Header
		getKey         string
		expectedGet    string
		getAllKey      string
		expectedGetAll []string
		expectedKeys   []string
	}{
		{
			name:           "Empty headers",
			headers:        http.Header{},
			getKey:         "Content-Type",
			expectedGet:    "",
			getAllKey:      "Content-Type",
			expectedGetAll: nil,
			expectedKeys:   []string{},
		},
		{
			name: "Single header",
			headers: http.Header{
				"Content-Type": []string{"application/json"},
			},
			getKey:         "Content-Type",
			expectedGet:    "application/json",
			getAllKey:      "Content-Type",
			expectedGetAll: []string{"application/json"},
			expectedKeys:   []string{"Content-Type"},
		},
		{
			name: "Multiple headers",
			headers: http.Header{
				"Content-Type":    []string{"application/json"},
				"Accept-Encoding": []string{"gzip", "deflate"},
			},
			getKey:         "Accept-Encoding",
			expectedGet:    "gzip",
			getAllKey:      "Accept-Encoding",
			expectedGetAll: []string{"gzip", "deflate"},
			expectedKeys:   []string{"Accept-Encoding", "Content-Type"},
		},
		{
			name: "Case-insensitive header keys",
			headers: http.Header{
				"Content-Type": []string{"application/json"},
			},
			getKey:         "CONTENT-TYPE",
			expectedGet:    "application/json",
			getAllKey:      "Content-Type",
			expectedGetAll: []string{"application/json"},
			expectedKeys:   []string{"Content-Type"},
		},
		{
			name: "Non-existent header",
			headers: http.Header{
				"Content-Type": []string{"application/json"},
			},
			getKey:         "X-Custom-Header",
			expectedGet:    "",
			getAllKey:      "X-Custom-Header",
			expectedGetAll: nil,
			expectedKeys:   []string{"Content-Type"},
		},
		{
			name: "Multiple values for single header",
			headers: http.Header{
				"Set-Cookie": []string{"session=abc123", "user=john"},
			},
			getKey:         "Set-Cookie",
			expectedGet:    "session=abc123",
			getAllKey:      "Set-Cookie",
			expectedGetAll: []string{"session=abc123", "user=john"},
			expectedKeys:   []string{"Set-Cookie"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			response := &Response{
				header: &ResponseFluentHeader{
					header: tt.headers,
				},
			}

			// Act
			result := response.Header()

			// Assert
			if got := result.Get(tt.getKey); got != tt.expectedGet {
				t.Errorf("Get(%q) got %q, want %q", tt.getKey, got, tt.expectedGet)
			}
			if got := result.GetAll(tt.getAllKey); !reflect.DeepEqual(got, tt.expectedGetAll) {
				t.Errorf("GetAll(%q) got %v, want %v", tt.getAllKey, got, tt.expectedGetAll)
			}
			got := result.Keys()
			want := tt.expectedKeys
			sort.Strings(got)
			sort.Strings(want)
			if !reflect.DeepEqual(got, want) {
				t.Errorf("Keys() got %v, want %v", got, want)
			}
		})
	}
}
