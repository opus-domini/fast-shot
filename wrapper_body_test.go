package fastshot

import (
	"bytes"
	"errors"
	"io"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBufferedBody(t *testing.T) {
	tests := []struct {
		name          string
		setup         func(*BufferedBody)
		method        func(*BufferedBody) (interface{}, error)
		expected      interface{}
		expectedError error
	}{
		{
			name: "Read success",
			setup: func(b *BufferedBody) {
				b.buffer.WriteString("hello world")
			},
			method: func(b *BufferedBody) (interface{}, error) {
				buf := make([]byte, 5)
				n, err := b.Read(buf)
				return string(buf[:n]), err
			},
			expected:      "hello",
			expectedError: nil,
		},
		{
			name: "Read more than buffer size",
			setup: func(b *BufferedBody) {
				b.buffer.WriteString("hello")
			},
			method: func(b *BufferedBody) (interface{}, error) {
				buf := make([]byte, 10)
				n, err := b.Read(buf)
				return string(buf[:n]), err
			},
			expected:      "hello",
			expectedError: nil,
		},
		{
			name:  "Read empty buffer",
			setup: func(b *BufferedBody) {},
			method: func(b *BufferedBody) (interface{}, error) {
				buf := make([]byte, 5)
				n, err := b.Read(buf)
				return string(buf[:n]), err
			},
			expected:      "",
			expectedError: io.EOF,
		},
		{
			name: "Close success",
			setup: func(b *BufferedBody) {
				// No setup needed for Close()
			},
			method: func(b *BufferedBody) (interface{}, error) {
				return nil, b.Close()
			},
			expected:      nil,
			expectedError: nil,
		},
		{
			name: "ReadAsJSON success",
			setup: func(b *BufferedBody) {
				b.buffer.WriteString(`{"key": "value"}`)
			},
			method: func(b *BufferedBody) (interface{}, error) {
				var result map[string]string
				err := b.ReadAsJSON(&result)
				return result, err
			},
			expected:      map[string]string{"key": "value"},
			expectedError: nil,
		},
		{
			name: "ReadAsJSON error",
			setup: func(b *BufferedBody) {
				b.buffer.WriteString(`invalid json`)
			},
			method: func(b *BufferedBody) (interface{}, error) {
				var result map[string]string
				err := b.ReadAsJSON(&result)
				return nil, err
			},
			expected:      nil,
			expectedError: errors.New("invalid character 'i' looking for beginning of value"),
		},
		{
			name: "WriteAsJSON success",
			setup: func(b *BufferedBody) {
				// No setup needed
			},
			method: func(b *BufferedBody) (interface{}, error) {
				return nil, b.WriteAsJSON(map[string]string{"key": "value"})
			},
			expected:      nil,
			expectedError: nil,
		},
		{
			name: "ReadAsString success",
			setup: func(b *BufferedBody) {
				b.buffer.WriteString("hello world")
			},
			method: func(b *BufferedBody) (interface{}, error) {
				return b.ReadAsString()
			},
			expected:      "hello world",
			expectedError: nil,
		},
		{
			name: "WriteAsString success",
			setup: func(b *BufferedBody) {
				// No setup needed
			},
			method: func(b *BufferedBody) (interface{}, error) {
				return nil, b.WriteAsString("hello world")
			},
			expected:      nil,
			expectedError: nil,
		},
		{
			name: "Set success",
			setup: func(b *BufferedBody) {
				// No setup needed
			},
			method: func(b *BufferedBody) (interface{}, error) {
				return nil, b.Set(strings.NewReader("hello world"))
			},
			expected:      nil,
			expectedError: nil,
		},
		{
			name: "Unwrap success",
			setup: func(b *BufferedBody) {
				b.buffer.WriteString("hello world")
			},
			method: func(b *BufferedBody) (interface{}, error) {
				buf := new(bytes.Buffer)
				_, err := buf.ReadFrom(b.Unwrap())
				return buf.String(), err
			},
			expected:      "hello world",
			expectedError: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			body := newBufferedBody()
			tt.setup(body)

			// Act
			result, err := tt.method(body)

			// Assert
			assert.Equal(t, tt.expected, result)
			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError.Error(), err.Error())
			}
		})
	}
}

func TestUnbufferedBody(t *testing.T) {
	tests := []struct {
		name          string
		reader        io.ReadCloser
		method        func(*UnbufferedBody) (interface{}, error)
		expected      interface{}
		expectedError error
	}{
		{
			name:   "Read success",
			reader: io.NopCloser(strings.NewReader("hello world")),
			method: func(b *UnbufferedBody) (interface{}, error) {
				buf := make([]byte, 5)
				n, err := b.Read(buf)
				return string(buf[:n]), err
			},
			expected:      "hello",
			expectedError: nil,
		},
		{
			name:   "Read more than buffer size",
			reader: io.NopCloser(strings.NewReader("hello")),
			method: func(b *UnbufferedBody) (interface{}, error) {
				buf := make([]byte, 10)
				n, err := b.Read(buf)
				return string(buf[:n]), err
			},
			expected:      "hello",
			expectedError: nil,
		},
		{
			name:   "Close success",
			reader: io.NopCloser(strings.NewReader("hello world")),
			method: func(b *UnbufferedBody) (interface{}, error) {
				return nil, b.Close()
			},
			expected:      nil,
			expectedError: nil,
		},
		{
			name:   "ReadAsJSON success",
			reader: io.NopCloser(strings.NewReader(`{"key": "value"}`)),
			method: func(b *UnbufferedBody) (interface{}, error) {
				var result map[string]string
				err := b.ReadAsJSON(&result)
				return result, err
			},
			expected:      map[string]string{"key": "value"},
			expectedError: nil,
		},
		{
			name:   "ReadAsJSON error",
			reader: io.NopCloser(strings.NewReader(`invalid json`)),
			method: func(b *UnbufferedBody) (interface{}, error) {
				var result map[string]string
				err := b.ReadAsJSON(&result)
				return nil, err
			},
			expected:      nil,
			expectedError: errors.New("invalid character 'i' looking for beginning of value"),
		},
		{
			name:   "WriteAsJSON success",
			reader: io.NopCloser(strings.NewReader("")),
			method: func(b *UnbufferedBody) (interface{}, error) {
				return nil, b.WriteAsJSON(map[string]string{"key": "value"})
			},
			expected:      nil,
			expectedError: nil,
		},
		{
			name:   "ReadAsString success",
			reader: io.NopCloser(strings.NewReader("hello world")),
			method: func(b *UnbufferedBody) (interface{}, error) {
				return b.ReadAsString()
			},
			expected:      "hello world",
			expectedError: nil,
		},
		{
			name:   "WriteAsString success",
			reader: io.NopCloser(strings.NewReader("")),
			method: func(b *UnbufferedBody) (interface{}, error) {
				return nil, b.WriteAsString("hello world")
			},
			expected:      nil,
			expectedError: nil,
		},
		{
			name:   "Set success",
			reader: io.NopCloser(strings.NewReader("")),
			method: func(b *UnbufferedBody) (interface{}, error) {
				return nil, b.Set(strings.NewReader("hello world"))
			},
			expected:      nil,
			expectedError: nil,
		},
		{
			name:   "Unwrap success",
			reader: io.NopCloser(strings.NewReader("hello world")),
			method: func(b *UnbufferedBody) (interface{}, error) {
				buf := new(bytes.Buffer)
				_, err := buf.ReadFrom(b.Unwrap())
				return buf.String(), err
			},
			expected:      "hello world",
			expectedError: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			body := newUnbufferedBody(tt.reader)

			// Act
			result, err := tt.method(body)

			// Assert
			assert.Equal(t, tt.expected, result)
			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError.Error(), err.Error())
			}
		})
	}
}
