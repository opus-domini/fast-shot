package fastshot

import (
	"errors"
	"io"
	"testing"

	"github.com/opus-domini/fast-shot/mock"
	"github.com/stretchr/testify/assert"
	tmock "github.com/stretchr/testify/mock"
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
				m.On("Read", tmock.Anything).Run(func(args tmock.Arguments) {
					b := args.Get(0).([]byte)
					copy(b, "hello")
				}).Return(5, io.EOF).Once()
				m.On("Close").Return(nil).Once()
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
				m.On("Read", tmock.Anything).Return(0, errors.New("read error")).Once()
				m.On("Close").Return(nil).Once()
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
				m.On("ReadAsString").Return("hello", nil).Once()
				m.On("Close").Return(nil).Once()
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
				m.On("ReadAsString").Return("", errors.New("string error")).Once()
				m.On("Close").Return(nil).Once()
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
				m.On("ReadAsJSON", tmock.Anything).Run(func(args tmock.Arguments) {
					arg := args.Get(0).(*map[string]string)
					*arg = map[string]string{"key": "value"}
				}).Return(nil).Once()
				m.On("Close").Return(nil).Once()
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
				m.On("ReadAsJSON", tmock.Anything).Return(errors.New("json error")).Once()
				m.On("Close").Return(nil).Once()
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
				m.On("ReadAsXML", tmock.Anything).Run(func(args tmock.Arguments) {
					arg := args.Get(0).(*struct {
						Key string `xml:"Key"`
					})
					*arg = struct {
						Key string `xml:"Key"`
					}{
						Key: "value",
					}
				}).Return(nil).Once()
				m.On("Close").Return(nil).Once()
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
				m.On("ReadAsXML", tmock.Anything).Return(errors.New("xml error")).Once()
				m.On("Close").Return(nil).Once()
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
			mockBody := new(mock.BodyWrapper)
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
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError.Error(), err.Error())
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, tt.expected, result)

			mockBody.AssertExpectations(t)
		})
	}
}

func TestResponseFluentBodyClose(t *testing.T) {
	t.Run("Close ignores error", func(t *testing.T) {
		mockBody := new(mock.BodyWrapper)
		mockBody.On("Close").Return(errors.New("close error")).Once()

		rb := &ResponseFluentBody{
			body: mockBody,
		}

		rb.Close()

		mockBody.AssertExpectations(t)
	})
}

func TestResponseFluentBodyCloseErr(t *testing.T) {
	t.Run("CloseErr success", func(t *testing.T) {
		mockBody := new(mock.BodyWrapper)
		mockBody.On("Close").Return(nil).Once()

		rb := &ResponseFluentBody{
			body: mockBody,
		}

		err := rb.CloseErr()

		assert.NoError(t, err)
		mockBody.AssertExpectations(t)
	})

	t.Run("CloseErr error", func(t *testing.T) {
		mockBody := new(mock.BodyWrapper)
		mockBody.On("Close").Return(errors.New("close error")).Once()

		rb := &ResponseFluentBody{
			body: mockBody,
		}

		err := rb.CloseErr()

		assert.Error(t, err)
		assert.Equal(t, "close error", err.Error())
		mockBody.AssertExpectations(t)
	})
}
