package fastshot

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/opus-domini/fast-shot/constant"
	"io"
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
    buf := new(bytes.Buffer)
    _, err := buf.ReadFrom(body)
    if err != nil {
        b.requestConfig.Validations().Add(errors.Join(errors.New(constant.ErrMsgReadBody), err))
        return b.parentBuilder
    }
    b.requestConfig.SetBody(buf)
    return b.parentBuilder
}

// AsString sets the body as string.
func (b *RequestBodyBuilder) AsString(body string) *RequestBuilder {
	b.requestConfig.SetBody(bytes.NewBufferString(body))
	return b.parentBuilder
}

// AsJSON sets the body as JSON.
func (b *RequestBodyBuilder) AsJSON(obj interface{}) *RequestBuilder {
	// Marshal JSON
	bodyBytes, err := json.Marshal(obj)
	if err != nil {
		b.requestConfig.Validations().Add(errors.Join(errors.New(constant.ErrMsgMarshalJSON), err))
		return b.parentBuilder
	}
	// Set body
	b.requestConfig.SetBody(bytes.NewBuffer(bodyBytes))
	return b.parentBuilder
}
