package fastshot

import (
	"errors"
	"io"

	"github.com/opus-domini/fast-shot/constant"
)

// BuilderRequestBody is the interface that wraps the basic methods for setting custom HTTP Body's.
var _ BuilderRequestBody[RequestBuilder] = (*RequestBodyBuilder)(nil)

// RequestBodyBuilder serves as the main entry point for configuring BuilderRequestBody.
type RequestBodyBuilder struct {
	parentBuilder *RequestBuilder
	requestConfig *RequestConfigBase
}

// Body returns a new RequestBodyBuilder for setting custom HTTP Body's.
func (b *RequestBuilder) Body() *RequestBodyBuilder {
	return &RequestBodyBuilder{
		parentBuilder: b,
		requestConfig: b.request.config,
	}
}

// AsReader sets the body as IO Reader.
func (b *RequestBodyBuilder) AsReader(body io.Reader) *RequestBuilder {
	err := b.requestConfig.Body().Set(body)
	if err != nil {
		b.requestConfig.Validations().Add(errors.Join(errors.New(constant.ErrMsgSetBody), err))
	}
	return b.parentBuilder
}

// AsString sets the body as string.
func (b *RequestBodyBuilder) AsString(body string) *RequestBuilder {
	err := b.requestConfig.Body().WriteAsString(body)
	if err != nil {
		b.requestConfig.Validations().Add(errors.Join(errors.New(constant.ErrMsgSetBody), err))
	}
	return b.parentBuilder
}

// AsJSON sets the body as JSON.
func (b *RequestBodyBuilder) AsJSON(obj interface{}) *RequestBuilder {
	err := b.requestConfig.Body().WriteAsJSON(obj)
	if err != nil {
		b.requestConfig.Validations().Add(errors.Join(errors.New(constant.ErrMsgMarshalJSON), err))
	}
	return b.parentBuilder
}
