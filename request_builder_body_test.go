package fastshot

import (
	"errors"
	"strings"
	"testing"

	"github.com/opus-domini/fast-shot/constant"
	"github.com/opus-domini/fast-shot/mock"
	"github.com/stretchr/testify/assert"
	tmock "github.com/stretchr/testify/mock"
)

func TestRequestBodyBuilder(t *testing.T) {
	mockedErr := errors.New("mock error")
	tests := []struct {
		name          string
		setup         func(*RequestBodyBuilder)
		method        func(*RequestBodyBuilder) *RequestBuilder
		expectedError error
	}{
		{
			name: "AsReader success",
			setup: func(rb *RequestBodyBuilder) {
				mockBody := new(mock.BodyWrapper)
				mockBody.On("Set", tmock.Anything).Return(nil)
				rb.requestConfig.body = mockBody
			},
			method: func(rb *RequestBodyBuilder) *RequestBuilder {
				return rb.AsReader(strings.NewReader("test"))
			},
			expectedError: nil,
		},
		{
			name: "AsReader failure",
			setup: func(rb *RequestBodyBuilder) {
				mockBody := new(mock.BodyWrapper)
				mockBody.On("Set", tmock.Anything).Return(mockedErr)
				rb.requestConfig.body = mockBody
			},
			method: func(rb *RequestBodyBuilder) *RequestBuilder {
				return rb.AsReader(strings.NewReader("test"))
			},
			expectedError: errors.Join(errors.New(constant.ErrMsgSetBody), mockedErr),
		},
		{
			name: "AsString success",
			setup: func(rb *RequestBodyBuilder) {
				mockBody := new(mock.BodyWrapper)
				mockBody.On("WriteAsString", "test").Return(nil)
				rb.requestConfig.body = mockBody
			},
			method: func(rb *RequestBodyBuilder) *RequestBuilder {
				return rb.AsString("test")
			},
			expectedError: nil,
		},
		{
			name: "AsString failure",
			setup: func(rb *RequestBodyBuilder) {
				mockBody := new(mock.BodyWrapper)
				mockBody.On("WriteAsString", "test").Return(mockedErr)
				rb.requestConfig.body = mockBody
			},
			method: func(rb *RequestBodyBuilder) *RequestBuilder {
				return rb.AsString("test")
			},
			expectedError: errors.Join(errors.New(constant.ErrMsgSetBody), mockedErr),
		},
		{
			name: "AsJSON success",
			setup: func(rb *RequestBodyBuilder) {
				mockBody := new(mock.BodyWrapper)
				mockBody.On("WriteAsJSON", tmock.Anything).Return(nil)
				rb.requestConfig.body = mockBody
			},
			method: func(rb *RequestBodyBuilder) *RequestBuilder {
				return rb.AsJSON(map[string]string{"key": "value"})
			},
			expectedError: nil,
		},
		{
			name: "AsJSON failure",
			setup: func(rb *RequestBodyBuilder) {
				mockBody := new(mock.BodyWrapper)
				mockBody.On("WriteAsJSON", tmock.Anything).Return(mockedErr)
				rb.requestConfig.body = mockBody
			},
			method: func(rb *RequestBodyBuilder) *RequestBuilder {
				return rb.AsJSON(map[string]string{"key": "value"})
			},
			expectedError: errors.Join(errors.New(constant.ErrMsgMarshalJSON), mockedErr),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			rb := &RequestBodyBuilder{
				parentBuilder: &RequestBuilder{},
				requestConfig: &RequestConfigBase{
					validations: newDefaultValidations(nil),
				},
			}
			tt.setup(rb)

			// Act
			result := tt.method(rb)

			// Assert
			assert.Equal(t, rb.parentBuilder, result)
			if tt.expectedError != nil {
				err := rb.requestConfig.validations.Get(0)
				assert.Error(t, err)
				assert.Equal(t, err, tt.expectedError)
			} else {
				assert.Empty(t, rb.requestConfig.validations.Unwrap())
			}
		})
	}
}
