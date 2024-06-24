package fastshot

import (
	"errors"
	"net/url"
	"testing"

	"github.com/opus-domini/fast-shot/constant"
	"github.com/stretchr/testify/assert"
)

func TestRequestQueryBuilder(t *testing.T) {
	tests := []struct {
		name           string
		method         func(*RequestBuilder) *RequestBuilder
		expectedQuery  url.Values
		expectedErrors []error
	}{
		{
			name: "Add single parameter",
			method: func(rb *RequestBuilder) *RequestBuilder {
				return rb.Query().AddParam("key", "value")
			},
			expectedQuery: url.Values{
				"key": {"value"},
			},
		},
		{
			name: "Add multiple parameters",
			method: func(rb *RequestBuilder) *RequestBuilder {
				return rb.Query().AddParams(map[string]string{
					"key1": "value1",
					"key2": "value2",
				})
			},
			expectedQuery: url.Values{
				"key1": {"value1"},
				"key2": {"value2"},
			},
		},
		{
			name: "Set single parameter",
			method: func(rb *RequestBuilder) *RequestBuilder {
				return rb.Query().SetParam("key", "value")
			},
			expectedQuery: url.Values{
				"key": {"value"},
			},
		},
		{
			name: "Set multiple parameters",
			method: func(rb *RequestBuilder) *RequestBuilder {
				return rb.Query().SetParams(map[string]string{
					"key1": "value1",
					"key2": "value2",
				})
			},
			expectedQuery: url.Values{
				"key1": {"value1"},
				"key2": {"value2"},
			},
		},
		{
			name: "Set valid raw query string",
			method: func(rb *RequestBuilder) *RequestBuilder {
				return rb.Query().SetRawString("key1=value1&key2=value2")
			},
			expectedQuery: url.Values{
				"key1": {"value1"},
				"key2": {"value2"},
			},
		},
		{
			name: "Set invalid raw query string",
			method: func(rb *RequestBuilder) *RequestBuilder {
				return rb.Query().SetRawString("invalid=%%")
			},
			expectedQuery: url.Values{},
			expectedErrors: []error{
				errors.New(constant.ErrMsgParseQueryString),
			},
		},
		{
			name: "Set empty raw query string",
			method: func(rb *RequestBuilder) *RequestBuilder {
				return rb.Query().SetRawString("")
			},
			expectedQuery: url.Values{},
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
			assert.Equal(t, tt.expectedQuery, rb.request.config.QueryParams())

			if tt.expectedErrors != nil {
				assert.Len(t, rb.request.config.Validations().Unwrap(), len(tt.expectedErrors))
				for i, expectedErr := range tt.expectedErrors {
					assert.ErrorContains(t, rb.request.config.Validations().Get(i), expectedErr.Error())
				}
			} else {
				assert.Empty(t, rb.request.config.Validations().Unwrap())
			}
		})
	}
}
