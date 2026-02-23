package fastshot

import (
	"errors"
	"io"
	"reflect"
	"testing"

	"github.com/opus-domini/fast-shot/mock"
)

func TestResponseFluentBody(t *testing.T) {
	tests := []struct {
		name          string
		setup         func(*mock.BodyWrapper)
		method        func(*ResponseFluentBody) (interface{}, error)
		expected      interface{}
		expectedError error
	}{
		{
			name: "Raw success",
			setup: func(m *mock.BodyWrapper) {
				// No setup needed for Raw method
			},
			method: func(rb *ResponseFluentBody) (interface{}, error) {
				return rb.Raw(), nil
			},
			expected:      &mock.BodyWrapper{},
			expectedError: nil,
		},
		{
			name: "AsBytes success",
			setup: func(m *mock.BodyWrapper) {
				m.ReadFunc = func(p []byte) (int, error) {
					copy(p, "hello")
					return 5, io.EOF
				}
				m.CloseFunc = func() error { return nil }
			},
			method: func(rb *ResponseFluentBody) (interface{}, error) {
				return rb.AsBytes()
			},
			expected:      []byte("hello"),
			expectedError: nil,
		},
		{
			name: "AsBytes read error",
			setup: func(m *mock.BodyWrapper) {
				m.ReadFunc = func(p []byte) (int, error) {
					return 0, errors.New("read error")
				}
				m.CloseFunc = func() error { return nil }
			},
			method: func(rb *ResponseFluentBody) (interface{}, error) {
				return rb.AsBytes()
			},
			expected:      []byte(nil),
			expectedError: errors.New("read error"),
		},
		{
			name: "AsString success",
			setup: func(m *mock.BodyWrapper) {
				m.ReadAsStringFunc = func() (string, error) { return "hello", nil }
				m.CloseFunc = func() error { return nil }
			},
			method: func(rb *ResponseFluentBody) (interface{}, error) {
				return rb.AsString()
			},
			expected:      "hello",
			expectedError: nil,
		},
		{
			name: "AsString error",
			setup: func(m *mock.BodyWrapper) {
				m.ReadAsStringFunc = func() (string, error) { return "", errors.New("string error") }
				m.CloseFunc = func() error { return nil }
			},
			method: func(rb *ResponseFluentBody) (interface{}, error) {
				return rb.AsString()
			},
			expected:      "",
			expectedError: errors.New("string error"),
		},
		{
			name: "AsJSON success",
			setup: func(m *mock.BodyWrapper) {
				m.ReadAsJSONFunc = func(obj interface{}) error {
					arg := obj.(*map[string]string)
					*arg = map[string]string{"key": "value"}
					return nil
				}
				m.CloseFunc = func() error { return nil }
			},
			method: func(rb *ResponseFluentBody) (interface{}, error) {
				var result map[string]string
				err := rb.AsJSON(&result)
				return result, err
			},
			expected:      map[string]string{"key": "value"},
			expectedError: nil,
		},
		{
			name: "AsJSON error",
			setup: func(m *mock.BodyWrapper) {
				m.ReadAsJSONFunc = func(obj interface{}) error { return errors.New("json error") }
				m.CloseFunc = func() error { return nil }
			},
			method: func(rb *ResponseFluentBody) (interface{}, error) {
				var result map[string]string
				err := rb.AsJSON(&result)
				return result, err
			},
			expected:      map[string]string(nil),
			expectedError: errors.New("json error"),
		},
		{
			name: "AsXML success",
			setup: func(m *mock.BodyWrapper) {
				m.ReadAsXMLFunc = func(obj interface{}) error {
					arg := obj.(*struct {
						Key string `xml:"Key"`
					})
					*arg = struct {
						Key string `xml:"Key"`
					}{
						Key: "value",
					}
					return nil
				}
				m.CloseFunc = func() error { return nil }
			},
			method: func(rb *ResponseFluentBody) (interface{}, error) {
				result := struct {
					Key string `xml:"Key"`
				}{}
				err := rb.AsXML(&result)
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
			name: "AsXML error",
			setup: func(m *mock.BodyWrapper) {
				m.ReadAsXMLFunc = func(obj interface{}) error { return errors.New("xml error") }
				m.CloseFunc = func() error { return nil }
			},
			method: func(rb *ResponseFluentBody) (interface{}, error) {
				result := struct {
					Key string `xml:"Key"`
				}{}
				err := rb.AsXML(&result)
				return result, err
			},
			expected: struct {
				Key string `xml:"Key"`
			}{},
			expectedError: errors.New("xml error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			mockBody := &mock.BodyWrapper{}
			tt.setup(mockBody)

			response := &Response{
				body: &ResponseFluentBody{
					body: mockBody,
				},
			}

			// Act
			result, err := tt.method(response.Body())

			// Assert
			if tt.expectedError != nil {
				if err == nil {
					t.Error("expected error, got nil")
				} else if err.Error() != tt.expectedError.Error() {
					t.Errorf("error got %q, want %q", err.Error(), tt.expectedError.Error())
				}
			} else {
				if err != nil {
					t.Fatalf("unexpected error: %v", err)
				}
			}
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("got %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestResponseFluentBodyClose(t *testing.T) {
	t.Run("Close ignores error", func(t *testing.T) {
		called := false
		mockBody := &mock.BodyWrapper{
			CloseFunc: func() error {
				called = true
				return errors.New("close error")
			},
		}

		rb := &ResponseFluentBody{
			body: mockBody,
		}

		rb.Close()

		if !called {
			t.Error("expected Close to be called")
		}
	})
}

func TestResponseFluentBodyCloseErr(t *testing.T) {
	t.Run("CloseErr success", func(t *testing.T) {
		mockBody := &mock.BodyWrapper{
			CloseFunc: func() error { return nil },
		}

		rb := &ResponseFluentBody{
			body: mockBody,
		}

		err := rb.CloseErr()

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})

	t.Run("CloseErr error", func(t *testing.T) {
		mockBody := &mock.BodyWrapper{
			CloseFunc: func() error { return errors.New("close error") },
		}

		rb := &ResponseFluentBody{
			body: mockBody,
		}

		err := rb.CloseErr()

		if err == nil {
			t.Error("expected error, got nil")
		} else if err.Error() != "close error" {
			t.Errorf("error got %q, want %q", err.Error(), "close error")
		}
	})
}
