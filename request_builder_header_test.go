package fastshot

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/opus-domini/fast-shot/constant/header"
	"github.com/opus-domini/fast-shot/constant/mime"
)

func TestRequestHeaderBuilder(t *testing.T) {
	tests := []struct {
		name           string
		method         func(*RequestBuilder) *RequestBuilder
		expectedHeader map[header.Type]string
	}{
		{
			name: "Add single header",
			method: func(rb *RequestBuilder) *RequestBuilder {
				return rb.Header().Add(header.ContentType, mime.JSON.String())
			},
			expectedHeader: map[header.Type]string{
				header.ContentType: mime.JSON.String(),
			},
		},
		{
			name: "Add multiple headers",
			method: func(rb *RequestBuilder) *RequestBuilder {
				return rb.Header().AddAll(map[header.Type]string{
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
			method: func(rb *RequestBuilder) *RequestBuilder {
				return rb.Header().Set(header.ContentType, mime.JSON.String())
			},
			expectedHeader: map[header.Type]string{
				header.ContentType: mime.JSON.String(),
			},
		},
		{
			name: "Set multiple headers",
			method: func(rb *RequestBuilder) *RequestBuilder {
				return rb.Header().SetAll(map[header.Type]string{
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
			method: func(rb *RequestBuilder) *RequestBuilder {
				return rb.Header().AddAccept(mime.JSON)
			},
			expectedHeader: map[header.Type]string{
				header.Accept: mime.JSON.String(),
			},
		},
		{
			name: "Add Content-Type header",
			method: func(rb *RequestBuilder) *RequestBuilder {
				return rb.Header().AddContentType(mime.JSON)
			},
			expectedHeader: map[header.Type]string{
				header.ContentType: mime.JSON.String(),
			},
		},
		{
			name: "Add User-Agent header",
			method: func(rb *RequestBuilder) *RequestBuilder {
				return rb.Header().AddUserAgent("TestAgent")
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
					config: newRequestConfigBase("", ""),
				},
			}

			// Act
			result := tt.method(rb)

			// Assert
			assert.Equal(t, rb, result)
			for key, value := range tt.expectedHeader {
				assert.Equal(t, value, rb.request.config.Header().Get(key))
			}
		})
	}
}
