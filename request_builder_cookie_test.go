package fastshot

import (
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
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
			assert.Equal(t, rb, result)
			assert.Equal(t, 1, rb.request.config.Cookies().Count())
			assert.Equal(t, tt.cookie, rb.request.config.Cookies().Get(0))
		})
	}
}
