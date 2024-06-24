package fastshot

import (
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestClientCookieBuilder(t *testing.T) {
	tests := []struct {
		name       string
		cookie     *http.Cookie
		assertFunc func(*testing.T, *ClientBuilder)
	}{
		{
			name: "Add simple cookie",
			cookie: &http.Cookie{
				Name:  "session",
				Value: "abc123",
			},
			assertFunc: func(t *testing.T, cb *ClientBuilder) {
				assert.Equal(t, 1, cb.client.Cookies().Count())
				cookie := cb.client.Cookies().Get(0)
				assert.Equal(t, "session", cookie.Name)
				assert.Equal(t, "abc123", cookie.Value)
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
			assertFunc: func(t *testing.T, cb *ClientBuilder) {
				assert.Equal(t, 1, cb.client.Cookies().Count())
				cookie := cb.client.Cookies().Get(0)
				assert.Equal(t, "complex", cookie.Name)
				assert.Equal(t, "value", cookie.Value)
				assert.Equal(t, "/", cookie.Path)
				assert.Equal(t, "example.com", cookie.Domain)
				assert.True(t, cookie.Expires.After(time.Now()))
				assert.Equal(t, 86400, cookie.MaxAge)
				assert.True(t, cookie.Secure)
				assert.True(t, cookie.HttpOnly)
				assert.Equal(t, http.SameSiteStrictMode, cookie.SameSite)
			},
		},
		{
			name: "Add multiple cookies",
			cookie: &http.Cookie{
				Name:  "first",
				Value: "value1",
			},
			assertFunc: func(t *testing.T, cb *ClientBuilder) {
				cb.Cookie().Add(&http.Cookie{Name: "second", Value: "value2"})
				assert.Equal(t, 2, cb.client.Cookies().Count())
				assert.Equal(t, "first", cb.client.Cookies().Get(0).Name)
				assert.Equal(t, "second", cb.client.Cookies().Get(1).Name)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			cb := NewClient("https://example.com")

			// Act
			result := cb.Cookie().Add(tt.cookie)

			// Assert
			assert.Equal(t, cb, result)
			tt.assertFunc(t, cb)
		})
	}
}
