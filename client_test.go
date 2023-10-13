package fastshot

import (
	"encoding/base64"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestClientBuilder_Auth(t *testing.T) {
	// Arrange
	builder := NewClient("https://example.com")
	// Act
	authBuilder := builder.Auth()
	// Assert
	if authBuilder.parentBuilder != builder {
		t.Errorf("Parent builder not set correctly")
	}
}

func TestClientAuthBuilder_Set(t *testing.T) {
	// Arrange
	builder := NewClient("https://example.com")
	// Act
	authBuilder := builder.Auth()
	authBuilder.Set("value")
	// Assert
	if builder.client.httpHeader.Get("Authorization") != "value" {
		t.Errorf("Authorization header not set correctly")
	}
}

func TestClientAuthBuilder_BasicAuth(t *testing.T) {
	// Arrange
	builder := NewClient("https://example.com")
	// Act
	authBuilder := builder.Auth()
	authBuilder.BasicAuth("username", "password")
	// Assert
	expected := "Basic " + base64.StdEncoding.EncodeToString([]byte("username:password"))
	if builder.client.httpHeader.Get("Authorization") != expected {
		t.Errorf(
			"Header not set correctly, got: %s, want: %s",
			builder.client.httpHeader.Get("Authorization"),
			expected,
		)
	}
}

func TestClientAuthBuilder_BearerToken(t *testing.T) {
	// Arrange
	builder := NewClient("https://example.com")
	// Act
	authBuilder := builder.Auth()
	authBuilder.BearerToken("token")
	// Assert
	if builder.client.httpHeader.Get("Authorization") != "Bearer token" {
		t.Errorf(
			"Header not set correctly, got: %s, want: %s",
			builder.client.httpHeader.Get("Authorization"),
			"Bearer token",
		)
	}
}

func TestClientAuthBuilder_End(t *testing.T) {
	// Arrange
	builder := NewClient("https://example.com")
	// Act
	authBuilder := builder.Auth()
	// Assert
	if authBuilder.End() != builder {
		t.Errorf("Parent builder not returned correctly")
	}
}

func TestClientBuilder_Header(t *testing.T) {
	// Arrange
	builder := NewClient("https://example.com")
	// Act
	headerBuilder := builder.Header()
	// Assert
	if headerBuilder.parentBuilder != builder {
		t.Errorf("Parent builder not set correctly")
	}
}

func TestClientHeaderBuilder_Add(t *testing.T) {
	// Arrange
	builder := NewClient("https://example.com")
	// Act
	headerBuilder := builder.Header()
	headerBuilder.Add("key", "value")
	headerBuilder.Add("key", "value2")
	// Assert
	if !strings.Contains(builder.client.httpHeader.Get("key"), "value") {
		t.Errorf("Header not set correctly")
	}
}

func TestClientHeaderBuilder_Set(t *testing.T) {
	// Arrange
	builder := NewClient("https://example.com")
	// Act
	headerBuilder := builder.Header()
	headerBuilder.Set("key", "value")
	// Assert
	if builder.client.httpHeader.Get("key") != "value" {
		t.Errorf("Header not set correctly")
	}
}

func TestClientHeaderBuilder_AddAccept(t *testing.T) {
	// Arrange
	builder := NewClient("https://example.com")
	valueToFind := "application/xml"
	// Act
	headerBuilder := builder.Header()
	headerBuilder.AddAccept("application/json")
	headerBuilder.AddAccept(valueToFind)
	// Assert
	values := builder.client.httpHeader.Values("Accept")
	valueFound := false
	for _, value := range values {
		if value == valueToFind {
			valueFound = true
			break
		}
	}
	if !valueFound {
		t.Errorf("Header not set correctly")
	}
}

func TestClientHeaderBuilder_AddUserAgent(t *testing.T) {
	// Arrange
	builder := NewClient("https://example.com")
	valueToFind := "chrome"
	// Act
	headerBuilder := builder.Header()
	headerBuilder.AddUserAgent("mobile")
	headerBuilder.AddUserAgent(valueToFind)
	headerBuilder.AddUserAgent("firefox")
	// Assert
	values := builder.client.httpHeader.Values("User-Agent")
	valueFound := false
	for _, value := range values {
		if value == valueToFind {
			valueFound = true
			break
		}
	}
	if !valueFound {
		t.Errorf("Header not set correctly")
	}
}

func TestClientHeaderBuilder_AddContentType(t *testing.T) {
	// Arrange
	builder := NewClient("https://example.com")
	valueToFind := "multipart/form-data; boundary=something"
	// Act
	headerBuilder := builder.Header()
	headerBuilder.AddContentType("text/html; charset=utf-8")
	headerBuilder.AddContentType(valueToFind)
	// Assert
	values := builder.client.httpHeader.Values("Content-Type")
	valueFound := false
	for _, value := range values {
		if value == valueToFind {
			valueFound = true
			break
		}
	}
	if !valueFound {
		t.Errorf("Header not set correctly")
	}
}

func TestClientHeaderBuilder_End(t *testing.T) {
	// Arrange
	builder := NewClient("https://example.com")
	// Act
	headerBuilder := builder.Header()
	// Assert
	if headerBuilder.End() != builder {
		t.Errorf("Parent builder not returned correctly")
	}
}

func TestClientBuilder_Cookie(t *testing.T) {
	// Arrange
	builder := NewClient("https://example.com")
	// Act
	cookieBuilder := builder.Cookie()
	// Assert
	if cookieBuilder.parentBuilder != builder {
		t.Errorf("Parent builder not set correctly")
	}
}

func TestClientCookieBuilder_Add(t *testing.T) {
	// Arrange
	builder := NewClient("https://example.com")
	// Act
	cookieBuilder := builder.Cookie()
	cookieBuilder.Add(&http.Cookie{Name: "name", Value: "value"})
	// Assert
	if len(builder.client.httpCookies) != 1 || builder.client.httpCookies[0].Name != "name" {
		t.Errorf("Cookie not set correctly")
	}
}

func TestClientCookieBuilder_End(t *testing.T) {
	// Arrange
	builder := NewClient("https://example.com")
	// Act
	cookieBuilder := builder.Cookie()
	// Assert
	if cookieBuilder.End() != builder {
		t.Errorf("Parent builder not returned correctly")
	}
}

func TestClientBuilder_Config(t *testing.T) {
	// Arrange
	builder := NewClient("https://example.com")
	// Act
	configBuilder := builder.Config()
	// Assert
	if configBuilder.parentBuilder != builder {
		t.Errorf("Parent builder not set correctly")
	}
}

func TestClientConfigBuilder_SetCustomTransport(t *testing.T) {
	// Arrange
	builder := NewClient("https://example.com")
	// Act
	configBuilder := builder.Config()
	configBuilder.SetCustomTransport(&http.Transport{})
	// Assert
	if builder.client.httpClient.Transport == nil {
		t.Errorf("Transport not set correctly")
	}
}

func TestClientConfigBuilder_SetFollowRedirects(t *testing.T) {
	// Arrange
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/redirect" {
			http.Redirect(w, r, "/target", http.StatusFound) // StatusFound 302
			return
		}
		_, _ = w.Write([]byte("OK"))
	}))
	defer server.Close()

	tests := []struct {
		name            string
		followRedirects bool
		wantFinalURL    string
	}{
		{
			name:            "Follow Redirects",
			followRedirects: true,
			wantFinalURL:    server.URL + "/target",
		},
		{
			name:            "Do Not Follow Redirects",
			followRedirects: false,
			wantFinalURL:    server.URL + "/redirect",
		},
	}

	// Act
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := NewClient(server.URL).
				Config().SetFollowRedirects(tt.followRedirects).End().
				Build()

			resp, err := client.GET("/redirect").Send()
			if err != nil {
				t.Fatalf("Failed: %v", err)
			}

			defer func(Body io.ReadCloser) {
				_ = Body.Close()
			}(resp.RawBody())

			// Assert
			if resp.RawResponse.Request.URL.String() != tt.wantFinalURL {
				t.Errorf("Final URL = %v, want %v", resp.RawResponse.Request.URL, tt.wantFinalURL)
			}
		})
	}
}

func TestClientConfigBuilder_SetTimeout(t *testing.T) {
	// Arrange
	builder := NewClient("https://example.com")
	// Act
	configBuilder := builder.Config()
	configBuilder.SetTimeout(1)
	// Assert
	if builder.client.httpClient.Timeout != 1 {
		t.Errorf("Timeout not set correctly")
	}
}

func TestClientConfigBuilder_End(t *testing.T) {
	// Arrange
	builder := NewClient("https://example.com")
	// Act
	configBuilder := builder.Config()
	// Assert
	if configBuilder.End() != builder {
		t.Errorf("Parent builder not returned correctly")
	}
}

func TestClientBuilder_Build(t *testing.T) {
	// Arrange
	builder := NewClient("https://example.com")
	// Act
	client := builder.Build()
	// Assert
	if client.baseURL != "https://example.com" {
		t.Errorf("BaseURL not set correctly")
	}
}

func TestDefaultClient(t *testing.T) {
	// Arrange
	client := DefaultClient("https://api.example.com")
	// Assert
	if client.baseURL != "https://api.example.com" {
		t.Errorf("BaseURL not set correctly")
	}
}
