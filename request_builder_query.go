package fastshot

import (
	"errors"
	"net/url"
	"strings"

	"github.com/opus-domini/fast-shot/constant"
)

// BuilderRequestQuery is the interface that wraps the basic methods for setting query parameters.
var _ BuilderRequestQuery[RequestBuilder] = (*RequestQueryBuilder)(nil)

// RequestQueryBuilder serves as the main entry point for configuring BuilderRequestQuery.
type RequestQueryBuilder struct {
	parentBuilder *RequestBuilder
	requestConfig *RequestConfigBase
}

// Query returns a new RequestQueryBuilder for setting query parameters.
func (b *RequestBuilder) Query() *RequestQueryBuilder {
	return &RequestQueryBuilder{
		parentBuilder: b,
		requestConfig: b.request.config,
	}
}

// AddParam adds a query parameter to the HTTP request. If parameter already exists, it will be appended.
func (b *RequestQueryBuilder) AddParam(param, value string) *RequestBuilder {
	b.requestConfig.queryParams.Add(param, value)
	return b.parentBuilder
}

// AddParams adds multiple query parameters to the HTTP request. If parameter already exists, it will be appended.
func (b *RequestQueryBuilder) AddParams(params map[string]string) *RequestBuilder {
	for param, value := range params {
		b.AddParam(param, value)
	}
	return b.parentBuilder
}

// SetParam sets a query parameter to the HTTP request. If parameter already exists, it will be overwritten.
func (b *RequestQueryBuilder) SetParam(param, value string) *RequestBuilder {
	b.requestConfig.queryParams.Set(param, value)
	return b.parentBuilder
}

// SetParams sets multiple query parameters to the HTTP request. If parameter already exists, it will be overwritten.
func (b *RequestQueryBuilder) SetParams(params map[string]string) *RequestBuilder {
	for param, value := range params {
		b.SetParam(param, value)
	}
	return b.parentBuilder
}

// SetRawString sets query parameters from a raw query string.
func (b *RequestQueryBuilder) SetRawString(query string) *RequestBuilder {
	// Parse query string
	queryParams, err := url.ParseQuery(strings.TrimSpace(query))
	if err != nil {
		b.requestConfig.Validations().Add(errors.Join(errors.New(constant.ErrMsgParseQueryString), err))
		return b.parentBuilder
	}
	// Set query params
	for param, values := range queryParams {
		for _, value := range values {
			b.SetParam(param, value)
		}
	}
	return b.parentBuilder
}
