package fastshot

import (
	"errors"
	"strings"
	"testing"

	"github.com/opus-domini/fast-shot/constant"
	"github.com/opus-domini/fast-shot/constant/header"
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
		{
			name: "AsXML success",
			setup: func(rb *RequestBodyBuilder) {
				mockBody := new(mock.BodyWrapper)
				mockBody.On("WriteAsXML", tmock.Anything).Return(nil)
				rb.requestConfig.body = mockBody
			},
			method: func(rb *RequestBodyBuilder) *RequestBuilder {
				body := `<example><Key>value</Key></example>`
				return rb.AsXML(&body)
			},
			expectedError: nil,
		},
		{
			name: "AsXML failure",
			setup: func(rb *RequestBodyBuilder) {
				mockBody := new(mock.BodyWrapper)
				mockBody.On("WriteAsXML", tmock.Anything).Return(mockedErr)
				rb.requestConfig.body = mockBody
			},
			method: func(rb *RequestBodyBuilder) *RequestBuilder {
				body := `<example><Key>value</Key></example>`
				return rb.AsXML(&body)
			},
			expectedError: errors.Join(errors.New(constant.ErrMsgMarshalXML), mockedErr),
		},
		{
			name: "AsFormData success",
			setup: func(rb *RequestBodyBuilder) {
				mockBody := new(mock.BodyWrapper)
				mockBody.On("WriteAsFormData", tmock.Anything).Return("multipart/form-data; boundary=test", nil)
				rb.requestConfig.body = mockBody
				rb.requestConfig.httpHeader = newDefaultHttpHeader()
			},
			method: func(rb *RequestBodyBuilder) *RequestBuilder {
				return rb.AsFormData(map[string]string{
					"key1": "value1",
					"key2": "value2",
				})
			},
			expectedError: nil,
		},
		{
			name: "AsFormData success with empty fields",
			setup: func(rb *RequestBodyBuilder) {
				mockBody := new(mock.BodyWrapper)
				mockBody.On("WriteAsFormData", tmock.Anything).Return("multipart/form-data; boundary=test", nil)
				rb.requestConfig.body = mockBody
				rb.requestConfig.httpHeader = newDefaultHttpHeader()
			},
			method: func(rb *RequestBodyBuilder) *RequestBuilder {
				return rb.AsFormData(map[string]string{})
			},
			expectedError: nil,
		},
		{
			name: "AsFormData failure",
			setup: func(rb *RequestBodyBuilder) {
				mockBody := new(mock.BodyWrapper)
				mockBody.On("WriteAsFormData", tmock.Anything).Return("", mockedErr)
				rb.requestConfig.body = mockBody
				rb.requestConfig.httpHeader = newDefaultHttpHeader()
			},
			method: func(rb *RequestBodyBuilder) *RequestBuilder {
				return rb.AsFormData(map[string]string{
					"key1": "value1",
				})
			},
			expectedError: errors.Join(errors.New(constant.ErrMsgSetBody), mockedErr),
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

func TestAsFormData_ContentType(t *testing.T) {
	tests := []struct {
		name           string
		fields         map[string]string
		expectHeader   bool
		headerContains string
	}{
		{
			name: "sets Content-Type header on success",
			fields: map[string]string{
				"field1": "value1",
				"field2": "value2",
			},
			expectHeader:   true,
			headerContains: "multipart/form-data",
		},
		{
			name:           "sets Content-Type header with empty fields",
			fields:         map[string]string{},
			expectHeader:   true,
			headerContains: "multipart/form-data",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			rb := &RequestBodyBuilder{
				parentBuilder: &RequestBuilder{},
				requestConfig: &RequestConfigBase{
					validations: newDefaultValidations(nil),
					body:        newBufferedBody(),
					httpHeader:  newDefaultHttpHeader(),
				},
			}

			// Act
			rb.AsFormData(tt.fields)

			// Assert
			if tt.expectHeader {
				contentType := rb.requestConfig.httpHeader.Get(header.ContentType)
				assert.NotEmpty(t, contentType)
				assert.Contains(t, contentType, tt.headerContains)
			}
			assert.Empty(t, rb.requestConfig.validations.Unwrap())
		})
	}
}

func TestAsFormData_NoContentTypeOnError(t *testing.T) {
	// Arrange
	mockedErr := errors.New("mock error")
	rb := &RequestBodyBuilder{
		parentBuilder: &RequestBuilder{},
		requestConfig: &RequestConfigBase{
			validations: newDefaultValidations(nil),
			httpHeader:  newDefaultHttpHeader(),
		},
	}
	mockBody := new(mock.BodyWrapper)
	mockBody.On("WriteAsFormData", tmock.Anything).Return("", mockedErr)
	rb.requestConfig.body = mockBody

	// Act
	rb.AsFormData(map[string]string{
		"key1": "value1",
	})

	// Assert
	contentType := rb.requestConfig.httpHeader.Get(header.ContentType)
	assert.Empty(t, contentType, "Content-Type should not be set when WriteAsFormData fails")
	err := rb.requestConfig.validations.Get(0)
	assert.Error(t, err)
}
