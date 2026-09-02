package fastshot

import (
	"testing"

	"github.com/opus-domini/fast-shot/constant/header"
	"github.com/opus-domini/fast-shot/constant/mime"
)

func assertHeader(t *testing.T, headers HeaderWrapper, expected map[header.Type]string) {
	t.Helper()
	for key, value := range expected {
		if got := headers.Get(key); got != value {
			t.Errorf("header %s got %q, want %q", key, got, value)
		}
	}
}

func TestHeaderBuilder_Request(t *testing.T) {
	tests := []struct {
		name           string
		method         func(*HeaderBuilder[*RequestBuilder]) *RequestBuilder
		expectedHeader map[header.Type]string
	}{
		{
			name: "Add single header",
			method: func(hb *HeaderBuilder[*RequestBuilder]) *RequestBuilder {
				return hb.Add(header.ContentType, mime.JSON.String())
			},
			expectedHeader: map[header.Type]string{
				header.ContentType: mime.JSON.String(),
			},
		},
		{
			name: "Add multiple headers",
			method: func(hb *HeaderBuilder[*RequestBuilder]) *RequestBuilder {
				return hb.AddAll(map[header.Type]string{
					header.ContentType: mime.JSON.String(),
					header.UserAgent:   "TestAgent",
				})
			},
			expectedHeader: map[header.Type]string{
				header.ContentType: mime.JSON.String(),
				header.UserAgent:   "TestAgent",
			},
		},
		{
			name: "Set single header",
			method: func(hb *HeaderBuilder[*RequestBuilder]) *RequestBuilder {
				return hb.Set(header.ContentType, mime.JSON.String())
			},
			expectedHeader: map[header.Type]string{
				header.ContentType: mime.JSON.String(),
			},
		},
		{
			name: "Set multiple headers",
			method: func(hb *HeaderBuilder[*RequestBuilder]) *RequestBuilder {
				return hb.SetAll(map[header.Type]string{
					header.ContentType: mime.JSON.String(),
					header.UserAgent:   "TestAgent",
				})
			},
			expectedHeader: map[header.Type]string{
				header.ContentType: mime.JSON.String(),
				header.UserAgent:   "TestAgent",
			},
		},
		{
			name: "Add Accept header",
			method: func(hb *HeaderBuilder[*RequestBuilder]) *RequestBuilder {
				return hb.AddAccept(mime.JSON)
			},
			expectedHeader: map[header.Type]string{
				header.Accept: mime.JSON.String(),
			},
		},
		{
			name: "Add Content-Type header",
			method: func(hb *HeaderBuilder[*RequestBuilder]) *RequestBuilder {
				return hb.AddContentType(mime.JSON)
			},
			expectedHeader: map[header.Type]string{
				header.ContentType: mime.JSON.String(),
			},
		},
		{
			name: "Add User-Agent header",
			method: func(hb *HeaderBuilder[*RequestBuilder]) *RequestBuilder {
				return hb.AddUserAgent("TestAgent")
			},
			expectedHeader: map[header.Type]string{
				header.UserAgent: "TestAgent",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			rb := &RequestBuilder{
				request: &Request{
					config: newRequestConfigBase("", "", DefaultJSONCodec()),
				},
			}

			// Act
			result := tt.method(rb.Header())

			// Assert
			if result != rb {
				t.Errorf("got different builder, want same")
			}
			assertHeader(t, rb.request.config.Header(), tt.expectedHeader)
		})
	}
}

func TestHeaderBuilder_Client(t *testing.T) {
	const testAgent = "TestAgent"

	tests := []struct {
		name           string
		method         func(*HeaderBuilder[*ClientBuilder]) *ClientBuilder
		expectedHeader map[header.Type]string
	}{
		{
			name: "Add single header",
			method: func(hb *HeaderBuilder[*ClientBuilder]) *ClientBuilder {
				return hb.Add(header.ContentType, mime.JSON.String())
			},
			expectedHeader: map[header.Type]string{
				header.ContentType: mime.JSON.String(),
			},
		},
		{
			name: "Add multiple headers",
			method: func(hb *HeaderBuilder[*ClientBuilder]) *ClientBuilder {
				return hb.AddAll(map[header.Type]string{
					header.ContentType: mime.JSON.String(),
					header.UserAgent:   testAgent,
				})
			},
			expectedHeader: map[header.Type]string{
				header.ContentType: mime.JSON.String(),
				header.UserAgent:   testAgent,
			},
		},
		{
			name: "Set single header",
			method: func(hb *HeaderBuilder[*ClientBuilder]) *ClientBuilder {
				return hb.Set(header.ContentType, mime.JSON.String())
			},
			expectedHeader: map[header.Type]string{
				header.ContentType: mime.JSON.String(),
			},
		},
		{
			name: "Set multiple headers",
			method: func(hb *HeaderBuilder[*ClientBuilder]) *ClientBuilder {
				return hb.SetAll(map[header.Type]string{
					header.ContentType: mime.JSON.String(),
					header.UserAgent:   testAgent,
				})
			},
			expectedHeader: map[header.Type]string{
				header.ContentType: mime.JSON.String(),
				header.UserAgent:   testAgent,
			},
		},
		{
			name: "Add Accept header",
			method: func(hb *HeaderBuilder[*ClientBuilder]) *ClientBuilder {
				return hb.AddAccept(mime.JSON)
			},
			expectedHeader: map[header.Type]string{
				header.Accept: mime.JSON.String(),
			},
		},
		{
			name: "Add Content-Type header",
			method: func(hb *HeaderBuilder[*ClientBuilder]) *ClientBuilder {
				return hb.AddContentType(mime.JSON)
			},
			expectedHeader: map[header.Type]string{
				header.ContentType: mime.JSON.String(),
			},
		},
		{
			name: "Add User-Agent header",
			method: func(hb *HeaderBuilder[*ClientBuilder]) *ClientBuilder {
				return hb.AddUserAgent(testAgent)
			},
			expectedHeader: map[header.Type]string{
				header.UserAgent: testAgent,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			cb := NewClient("https://example.com")

			// Act
			result := tt.method(cb.Header())

			// Assert
			if result != cb {
				t.Errorf("got different builder, want same")
			}
			assertHeader(t, cb.client.Header(), tt.expectedHeader)
		})
	}
}
