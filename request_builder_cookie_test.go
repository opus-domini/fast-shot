package fastshot

import (
	"net/http"
	"reflect"
	"testing"
	"time"
)

func TestRequestCookieBuilder(t *testing.T) {
	tests := []struct {
		name   string
		cookie *http.Cookie
	}{
		{
			name: "Add simple cookie",
			cookie: &http.Cookie{
				Name:  "session",
				Value: "abc123",
			},
		},
		{
			name: "Add complex cookie",
			cookie: &http.Cookie{
				Name:     "complex",
				Value:    "value",
				Path:     "/",
				Domain:   "example.com",
				Expires:  time.Now().Add(24 * time.Hour),
				MaxAge:   86400,
				Secure:   true,
				HttpOnly: true,
				SameSite: http.SameSiteStrictMode,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			rb := &RequestBuilder{
				request: &Request{
					config: newRequestConfigBase("", ""),
				},
			}

			// Act
			result := rb.Cookie().Add(tt.cookie)

			// Assert
			if result != rb {
				t.Errorf("got different builder, want same")
			}
			if got := rb.request.config.Cookies().Count(); got != 1 {
				t.Errorf("cookie count got %d, want 1", got)
			}
			if got := rb.request.config.Cookies().Get(0); !reflect.DeepEqual(got, tt.cookie) {
				t.Errorf("cookie got %v, want %v", got, tt.cookie)
			}
		})
	}
}
