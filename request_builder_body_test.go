package fastshot

import (
	"errors"
	"io"
	"strings"
	"testing"

	"github.com/opus-domini/fast-shot/constant"
	"github.com/opus-domini/fast-shot/constant/header"
	"github.com/opus-domini/fast-shot/mock"
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
				rb.requestConfig.body = &mock.BodyWrapper{
					SetFunc: func(body io.Reader) error { return nil },
				}
			},
			method: func(rb *RequestBodyBuilder) *RequestBuilder {
				return rb.AsReader(strings.NewReader("test"))
			},
			expectedError: nil,
		},
		{
			name: "AsReader failure",
			setup: func(rb *RequestBodyBuilder) {
				rb.requestConfig.body = &mock.BodyWrapper{
					SetFunc: func(body io.Reader) error { return mockedErr },
				}
			},
			method: func(rb *RequestBodyBuilder) *RequestBuilder {
				return rb.AsReader(strings.NewReader("test"))
			},
			expectedError: errors.Join(errors.New(constant.ErrMsgSetBody), mockedErr),
		},
		{
			name: "AsString success",
			setup: func(rb *RequestBodyBuilder) {
				rb.requestConfig.body = &mock.BodyWrapper{
					WriteAsStringFunc: func(body string) error { return nil },
				}
			},
			method: func(rb *RequestBodyBuilder) *RequestBuilder {
				return rb.AsString("test")
			},
			expectedError: nil,
		},
		{
			name: "AsString failure",
			setup: func(rb *RequestBodyBuilder) {
				rb.requestConfig.body = &mock.BodyWrapper{
					WriteAsStringFunc: func(body string) error { return mockedErr },
				}
			},
			method: func(rb *RequestBodyBuilder) *RequestBuilder {
				return rb.AsString("test")
			},
			expectedError: errors.Join(errors.New(constant.ErrMsgSetBody), mockedErr),
		},
		{
			name: "AsJSON success",
			setup: func(rb *RequestBodyBuilder) {
				rb.requestConfig.body = &mock.BodyWrapper{
					WriteAsJSONFunc: func(obj interface{}) error { return nil },
				}
			},
			method: func(rb *RequestBodyBuilder) *RequestBuilder {
				return rb.AsJSON(map[string]string{"key": "value"})
			},
			expectedError: nil,
		},
		{
			name: "AsJSON failure",
			setup: func(rb *RequestBodyBuilder) {
				rb.requestConfig.body = &mock.BodyWrapper{
					WriteAsJSONFunc: func(obj interface{}) error { return mockedErr },
				}
			},
			method: func(rb *RequestBodyBuilder) *RequestBuilder {
				return rb.AsJSON(map[string]string{"key": "value"})
			},
			expectedError: errors.Join(errors.New(constant.ErrMsgMarshalJSON), mockedErr),
		},
		{
			name: "AsXML success",
			setup: func(rb *RequestBodyBuilder) {
				rb.requestConfig.body = &mock.BodyWrapper{
					WriteAsXMLFunc: func(obj interface{}) error { return nil },
				}
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
				rb.requestConfig.body = &mock.BodyWrapper{
					WriteAsXMLFunc: func(obj interface{}) error { return mockedErr },
				}
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
				rb.requestConfig.body = &mock.BodyWrapper{
					WriteAsFormDataFunc: func(fields map[string]string) (string, error) {
						return "multipart/form-data; boundary=test", nil
					},
				}
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
				rb.requestConfig.body = &mock.BodyWrapper{
					WriteAsFormDataFunc: func(fields map[string]string) (string, error) {
						return "multipart/form-data; boundary=test", nil
					},
				}
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
				rb.requestConfig.body = &mock.BodyWrapper{
					WriteAsFormDataFunc: func(fields map[string]string) (string, error) {
						return "", mockedErr
					},
				}
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
			if result != rb.parentBuilder {
				t.Errorf("got different builder, want same")
			}
			if tt.expectedError != nil {
				err := rb.requestConfig.validations.Get(0)
				if err == nil {
					t.Error("expected error, got nil")
				} else if err.Error() != tt.expectedError.Error() {
					t.Errorf("error got %q, want %q", err.Error(), tt.expectedError.Error())
				}
			} else {
				if got := rb.requestConfig.validations.Unwrap(); len(got) != 0 {
					t.Errorf("validations got %v, want empty", got)
				}
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
				if contentType == "" {
					t.Error("Content-Type is empty, want non-empty")
				}
				if !strings.Contains(contentType, tt.headerContains) {
					t.Errorf("Content-Type %q does not contain %q", contentType, tt.headerContains)
				}
			}
			if got := rb.requestConfig.validations.Unwrap(); len(got) != 0 {
				t.Errorf("validations got %v, want empty", got)
			}
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
	rb.requestConfig.body = &mock.BodyWrapper{
		WriteAsFormDataFunc: func(fields map[string]string) (string, error) {
			return "", mockedErr
		},
	}

	// Act
	rb.AsFormData(map[string]string{
		"key1": "value1",
	})

	// Assert
	contentType := rb.requestConfig.httpHeader.Get(header.ContentType)
	if contentType != "" {
		t.Errorf("Content-Type should not be set when WriteAsFormData fails, got %q", contentType)
	}
	err := rb.requestConfig.validations.Get(0)
	if err == nil {
		t.Error("expected error, got nil")
	}
}
