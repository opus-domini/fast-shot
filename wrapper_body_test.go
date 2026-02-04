package fastshot

import (
	"bytes"
	"errors"
	"io"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWrapperBody_Buffered(t *testing.T) {
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
			name: "ReadAsXML success",
			setup: func(b *BufferedBody) {
				b.buffer.WriteString(`<example><Key>value</Key></example>`)
			},
			method: func(b *BufferedBody) (interface{}, error) {
				result := struct {
					Key string `xml:"Key"`
				}{}
				err := b.ReadAsXML(&result)
				return result, err
			},
			expected: struct {
				Key string `xml:"Key"`
			}{
				Key: "value",
			},
			expectedError: nil,
		},
		{
			name: "ReadAsXML error",
			setup: func(b *BufferedBody) {
				b.buffer.WriteString(`<>invalid xml`)
			},
			method: func(b *BufferedBody) (interface{}, error) {
				result := struct {
					Key string `xml:"Key"`
				}{}
				err := b.ReadAsXML(&result)
				return nil, err
			},
			expected:      nil,
			expectedError: errors.New("XML syntax error on line 1: expected element name after <"),
		},
		{
			name: "WriteAsXML success",
			setup: func(b *BufferedBody) {
				// No setup needed
			},
			method: func(b *BufferedBody) (interface{}, error) {
				result := struct {
					Key string `xml:"Key"`
				}{
					Key: "value",
				}
				return nil, b.WriteAsXML(&result)
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
			name: "WriteAsString twice to test buffer reset",
			setup: func(b *BufferedBody) {
				_ = b.WriteAsString("first")
				_ = b.WriteAsString("second")
			},
			method: func(b *BufferedBody) (interface{}, error) {
				buf := make([]byte, 6)
				n, err := b.Read(buf)
				return string(buf[:n]), err
			},
			expected:      "second",
			expectedError: nil,
		},
		{
			name: "ReadAsString twice to test buffer read",
			setup: func(b *BufferedBody) {
				_ = b.WriteAsString("body content")
			},
			method: func(b *BufferedBody) (interface{}, error) {
				_, _ = b.ReadAsString()
				str, err := b.ReadAsString()
				return str, err
			},
			expected:      "body content",
			expectedError: nil,
		},
		{
			name: "WriteAsFormData success",
			setup: func(b *BufferedBody) {
				// No setup needed
			},
			method: func(b *BufferedBody) (interface{}, error) {
				contentType, err := b.WriteAsFormData(map[string]string{"field": "value"})
				if err != nil {
					return nil, err
				}
				return contentType != "", nil
			},
			expected:      true,
			expectedError: nil,
		},
		{
			name: "WriteAsFormData empty fields",
			setup: func(b *BufferedBody) {
				// No setup needed
			},
			method: func(b *BufferedBody) (interface{}, error) {
				contentType, err := b.WriteAsFormData(map[string]string{})
				if err != nil {
					return nil, err
				}
				return contentType != "", nil
			},
			expected:      true,
			expectedError: nil,
		},
		{
			name: "WriteAsFormData error",
			setup: func(b *BufferedBody) {
				writeFormDataFn = func(_ io.Writer, _ map[string]string) (string, error) {
					return "", errors.New("form data error")
				}
			},
			method: func(b *BufferedBody) (interface{}, error) {
				defer func() { writeFormDataFn = writeFormData }()
				return b.WriteAsFormData(map[string]string{"k": "v"})
			},
			expected:      "",
			expectedError: errors.New("form data error"),
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

func TestWrapperBody_Unbuffered(t *testing.T) {
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
			name:   "WriteAsJSON error",
			reader: io.NopCloser(strings.NewReader("")),
			method: func(b *UnbufferedBody) (interface{}, error) {
				return nil, b.WriteAsJSON(make(chan int))
			},
			expected:      nil,
			expectedError: errors.New("json: unsupported type: chan int"),
		},
		{
			name:   "ReadAsXML success",
			reader: io.NopCloser(strings.NewReader(`<example><Key>value</Key></example>`)),
			method: func(b *UnbufferedBody) (interface{}, error) {
				result := struct {
					Key string `xml:"Key"`
				}{}
				err := b.ReadAsXML(&result)
				return result, err
			},
			expected: struct {
				Key string `xml:"Key"`
			}{
				Key: "value",
			},
			expectedError: nil,
		},
		{
			name:   "ReadAsXML error",
			reader: io.NopCloser(strings.NewReader(`<>invalid xml`)),
			method: func(b *UnbufferedBody) (interface{}, error) {
				result := struct {
					Key string `xml:"Key"`
				}{}
				err := b.ReadAsXML(&result)
				return nil, err
			},
			expected:      nil,
			expectedError: errors.New("XML syntax error on line 1: expected element name after <"),
		},
		{
			name:   "WriteAsXML success",
			reader: io.NopCloser(strings.NewReader("")),
			method: func(b *UnbufferedBody) (interface{}, error) {
				return nil, b.WriteAsXML(map[string]string{"key": "value"})
			},
			expected:      nil,
			expectedError: nil,
		},
		{
			name:   "WriteAsXML error",
			reader: io.NopCloser(strings.NewReader("")),
			method: func(b *UnbufferedBody) (interface{}, error) {
				return nil, b.WriteAsXML(make(chan int))
			},
			expected:      nil,
			expectedError: errors.New("xml: unsupported type: chan int"),
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
			name:   "ReadAsString twice to test empty reader",
			reader: io.NopCloser(strings.NewReader("hello world")),
			method: func(b *UnbufferedBody) (interface{}, error) {
				_, _ = b.ReadAsString()
				return b.ReadAsString()
			},
			expected:      "",
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
		{
			name:   "WriteAsJSON error",
			reader: io.NopCloser(strings.NewReader("")),
			method: func(b *UnbufferedBody) (interface{}, error) {
				return nil, b.WriteAsJSON(make(chan int))
			},
			expected:      nil,
			expectedError: errors.New("json: unsupported type: chan int"),
		},
		{
			name:   "ReadAsString error", // Covering "if err != nil"
			reader: io.NopCloser(&errorReader{}),
			method: func(b *UnbufferedBody) (interface{}, error) {
				return b.ReadAsString()
			},
			expected:      "",
			expectedError: errors.New("mock error"), // Expect the reader's error
		},
		{
			name:   "WriteAsFormData success",
			reader: io.NopCloser(strings.NewReader("")),
			method: func(b *UnbufferedBody) (interface{}, error) {
				contentType, err := b.WriteAsFormData(map[string]string{"field": "value"})
				if err != nil {
					return nil, err
				}
				return contentType != "", nil
			},
			expected:      true,
			expectedError: nil,
		},
		{
			name:   "WriteAsFormData empty fields",
			reader: io.NopCloser(strings.NewReader("")),
			method: func(b *UnbufferedBody) (interface{}, error) {
				contentType, err := b.WriteAsFormData(map[string]string{})
				if err != nil {
					return nil, err
				}
				return contentType != "", nil
			},
			expected:      true,
			expectedError: nil,
		},
		{
			name:   "WriteAsFormData error",
			reader: io.NopCloser(strings.NewReader("")),
			method: func(b *UnbufferedBody) (interface{}, error) {
				writeFormDataFn = func(_ io.Writer, _ map[string]string) (string, error) {
					return "", errors.New("form data error")
				}
				defer func() { writeFormDataFn = writeFormData }()
				return b.WriteAsFormData(map[string]string{"k": "v"})
			},
			expected:      "",
			expectedError: errors.New("form data error"),
		},
		{
			name:   "Set with io.ReadCloser", // Covering "if closer, ok..."
			reader: io.NopCloser(strings.NewReader("")),
			method: func(b *UnbufferedBody) (interface{}, error) {
				return nil, b.Set(io.NopCloser(strings.NewReader("test")))
			},
			expected:      nil,
			expectedError: nil,
		},
		{
			name:   "Set with io.Reader", // Covering "else" block
			reader: io.NopCloser(strings.NewReader("")),
			method: func(b *UnbufferedBody) (interface{}, error) {
				return nil, b.Set(strings.NewReader("test"))
			},
			expected:      nil,
			expectedError: nil,
		},
		{
			name:   "WriteAsJSON assigns reader contents",
			reader: io.NopCloser(strings.NewReader("")),
			method: func(b *UnbufferedBody) (interface{}, error) {
				// Use struct for deterministic JSON ordering
				type payload struct {
					Key string `json:"key"`
				}
				if err := b.WriteAsJSON(payload{Key: "value"}); err != nil {
					return nil, err
				}
				return b.ReadAsString()
			},
			expected:      "{\"key\":\"value\"}\n",
			expectedError: nil,
		},
		{
			name:   "WriteAsXML assigns reader contents",
			reader: io.NopCloser(strings.NewReader("")),
			method: func(b *UnbufferedBody) (interface{}, error) {
				type payload struct { // deterministic root element name "payload"
					Key string `xml:"Key"`
				}
				if err := b.WriteAsXML(payload{Key: "value"}); err != nil {
					return nil, err
				}
				return b.ReadAsString()
			},
			expected:      "<payload><Key>value</Key></payload>",
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

func TestWriteFormData(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		var buf bytes.Buffer
		contentType, err := writeFormData(&buf, map[string]string{"key": "value"})
		assert.NoError(t, err)
		assert.Contains(t, contentType, "multipart/form-data")
		assert.True(t, buf.Len() > 0)
	})

	t.Run("writer error on WriteField", func(t *testing.T) {
		w := &errorWriter{failAfter: 0}
		contentType, err := writeFormData(w, map[string]string{"key": "value"})
		assert.Error(t, err)
		assert.Empty(t, contentType)
	})

	t.Run("writer error on Close", func(t *testing.T) {
		// Allow enough bytes for the field write (~110 bytes) but fail
		// when Close writes the closing boundary (~68 more bytes).
		w := &errorWriter{failAfter: 150}
		contentType, err := writeFormData(w, map[string]string{"k": "v"})
		assert.Error(t, err)
		assert.Empty(t, contentType)
	})
}

// errorWriter fails after writing failAfter bytes.
type errorWriter struct {
	written   int
	failAfter int
}

func (w *errorWriter) Write(p []byte) (int, error) {
	if w.written+len(p) > w.failAfter {
		return 0, errors.New("write error")
	}
	w.written += len(p)
	return len(p), nil
}

// Helper type
type errorReader struct{}

func (r *errorReader) Read(_ []byte) (n int, err error) {
	return 0, errors.New("mock error")
}

func (r *errorReader) Close() error {
	return nil
}
