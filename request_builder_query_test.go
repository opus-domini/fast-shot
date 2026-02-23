package fastshot

import (
	"errors"
	"net/url"
	"reflect"
	"testing"

	"github.com/opus-domini/fast-shot/constant"
)

func TestRequestQueryBuilder(t *testing.T) {
	tests := []struct {
		name          string
		method        func(*RequestBuilder) *RequestBuilder
		expectedQuery url.Values
		expectedError error
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
			expectedError: errors.Join(errors.New(constant.ErrMsgParseQueryString), url.EscapeError("%%")),
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
			if result != rb {
				t.Errorf("got different builder, want same")
			}
			if got := rb.request.config.QueryParams(); !reflect.DeepEqual(got, tt.expectedQuery) {
				t.Errorf("QueryParams() got %v, want %v", got, tt.expectedQuery)
			}

			if tt.expectedError != nil {
				if got := len(rb.request.config.Validations().Unwrap()); got != 1 {
					t.Errorf("validations count got %d, want 1", got)
				}
				if got := rb.request.config.Validations().Get(0); got.Error() != tt.expectedError.Error() {
					t.Errorf("validation got %q, want %q", got.Error(), tt.expectedError.Error())
				}
			} else {
				if got := rb.request.config.Validations().Unwrap(); len(got) != 0 {
					t.Errorf("validations got %v, want empty", got)
				}
			}
		})
	}
}
