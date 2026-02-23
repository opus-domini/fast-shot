package fastshot

import (
	"net/http"
	"testing"
	"time"
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
				if got := cb.client.Cookies().Count(); got != 1 {
					t.Errorf("cookie count got %d, want 1", got)
				}
				cookie := cb.client.Cookies().Get(0)
				if cookie.Name != "session" {
					t.Errorf("Name got %q, want %q", cookie.Name, "session")
				}
				if cookie.Value != "abc123" {
					t.Errorf("Value got %q, want %q", cookie.Value, "abc123")
				}
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
				if got := cb.client.Cookies().Count(); got != 1 {
					t.Errorf("cookie count got %d, want 1", got)
				}
				cookie := cb.client.Cookies().Get(0)
				if cookie.Name != "complex" {
					t.Errorf("Name got %q, want %q", cookie.Name, "complex")
				}
				if cookie.Value != "value" {
					t.Errorf("Value got %q, want %q", cookie.Value, "value")
				}
				if cookie.Path != "/" {
					t.Errorf("Path got %q, want %q", cookie.Path, "/")
				}
				if cookie.Domain != "example.com" {
					t.Errorf("Domain got %q, want %q", cookie.Domain, "example.com")
				}
				if !cookie.Expires.After(time.Now()) {
					t.Errorf("Expires should be in the future")
				}
				if cookie.MaxAge != 86400 {
					t.Errorf("MaxAge got %d, want 86400", cookie.MaxAge)
				}
				if !cookie.Secure {
					t.Error("Secure got false, want true")
				}
				if !cookie.HttpOnly {
					t.Error("HttpOnly got false, want true")
				}
				if cookie.SameSite != http.SameSiteStrictMode {
					t.Errorf("SameSite got %v, want %v", cookie.SameSite, http.SameSiteStrictMode)
				}
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
				if got := cb.client.Cookies().Count(); got != 2 {
					t.Errorf("cookie count got %d, want 2", got)
				}
				if got := cb.client.Cookies().Get(0).Name; got != "first" {
					t.Errorf("first cookie Name got %q, want %q", got, "first")
				}
				if got := cb.client.Cookies().Get(1).Name; got != "second" {
					t.Errorf("second cookie Name got %q, want %q", got, "second")
				}
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
			if result != cb {
				t.Errorf("got different builder, want same")
			}
			tt.assertFunc(t, cb)
		})
	}
}
