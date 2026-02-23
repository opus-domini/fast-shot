package fastshot

import (
	"net/http"
	"testing"
	"time"
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
			if len(result) != len(tt.expected) {
				t.Fatalf("cookie count got %d, want %d", len(result), len(tt.expected))
			}
			for i, expectedCookie := range tt.expected {
				assertCookiesEqual(t, expectedCookie, result[i])
			}
		})
	}
}

func assertCookiesEqual(t *testing.T, expected, actual *http.Cookie) {
	t.Helper()
	if actual.Name != expected.Name {
		t.Errorf("Name got %q, want %q", actual.Name, expected.Name)
	}
	if actual.Value != expected.Value {
		t.Errorf("Value got %q, want %q", actual.Value, expected.Value)
	}
	if actual.Path != expected.Path {
		t.Errorf("Path got %q, want %q", actual.Path, expected.Path)
	}
	if actual.Domain != expected.Domain {
		t.Errorf("Domain got %q, want %q", actual.Domain, expected.Domain)
	}
	if actual.MaxAge != expected.MaxAge {
		t.Errorf("MaxAge got %d, want %d", actual.MaxAge, expected.MaxAge)
	}
	if actual.Secure != expected.Secure {
		t.Errorf("Secure got %v, want %v", actual.Secure, expected.Secure)
	}
	if actual.HttpOnly != expected.HttpOnly {
		t.Errorf("HttpOnly got %v, want %v", actual.HttpOnly, expected.HttpOnly)
	}
	if actual.SameSite != expected.SameSite {
		t.Errorf("SameSite got %v, want %v", actual.SameSite, expected.SameSite)
	}

	if !expected.Expires.IsZero() {
		if diff := expected.Expires.Sub(actual.Expires).Abs(); diff > time.Second {
			t.Errorf("Expires diff %v exceeds %v", diff, time.Second)
		}
	}
}
