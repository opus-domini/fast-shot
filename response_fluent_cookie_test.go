package fastshot

import (
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestResponseFluentCookie(t *testing.T) {
	tests := []struct {
		name     string
		cookies  []*http.Cookie
		expected []*http.Cookie
	}{
		{
			name:     "No cookies",
			cookies:  []*http.Cookie{},
			expected: []*http.Cookie{},
		},
		{
			name: "Single cookie",
			cookies: []*http.Cookie{
				{Name: "session", Value: "abc123"},
			},
			expected: []*http.Cookie{
				{Name: "session", Value: "abc123"},
			},
		},
		{
			name: "Multiple cookies",
			cookies: []*http.Cookie{
				{Name: "session", Value: "abc123"},
				{Name: "user", Value: "john_doe"},
			},
			expected: []*http.Cookie{
				{Name: "session", Value: "abc123"},
				{Name: "user", Value: "john_doe"},
			},
		},
		{
			name: "Cookies with various attributes",
			cookies: []*http.Cookie{
				{Name: "session", Value: "abc123", Path: "/", Domain: "example.com", Expires: time.Now().Add(24 * time.Hour), MaxAge: 86400, Secure: true, HttpOnly: true},
				{Name: "preference", Value: "dark_mode", SameSite: http.SameSiteStrictMode},
			},
			expected: []*http.Cookie{
				{Name: "session", Value: "abc123", Path: "/", Domain: "example.com", MaxAge: 86400, Secure: true, HttpOnly: true},
				{Name: "preference", Value: "dark_mode", SameSite: http.SameSiteStrictMode},
			},
		},
		{
			name: "Cookies with empty values",
			cookies: []*http.Cookie{
				{Name: "empty_cookie", Value: ""},
				{Name: "normal_cookie", Value: "has_value"},
			},
			expected: []*http.Cookie{
				{Name: "empty_cookie", Value: ""},
				{Name: "normal_cookie", Value: "has_value"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			response := &Response{
				cookie: &ResponseFluentCookie{
					cookies: tt.cookies,
				},
			}

			// Act
			result := response.Cookie().GetAll()

			// Assert
			assert.Equal(t, len(tt.expected), len(result))
			for i, expectedCookie := range tt.expected {
				assertCookiesEqual(t, expectedCookie, result[i])
			}
		})
	}
}

func assertCookiesEqual(t *testing.T, expected, actual *http.Cookie) {
	assert.Equal(t, expected.Name, actual.Name)
	assert.Equal(t, expected.Value, actual.Value)
	assert.Equal(t, expected.Path, actual.Path)
	assert.Equal(t, expected.Domain, actual.Domain)
	assert.Equal(t, expected.MaxAge, actual.MaxAge)
	assert.Equal(t, expected.Secure, actual.Secure)
	assert.Equal(t, expected.HttpOnly, actual.HttpOnly)
	assert.Equal(t, expected.SameSite, actual.SameSite)

	if !expected.Expires.IsZero() {
		assert.WithinDuration(t, expected.Expires, actual.Expires, time.Second)
	}
}
