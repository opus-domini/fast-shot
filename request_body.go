package fastshot

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/opus-domini/fast-shot/constant"
	"io"
)

// RequestBody is the interface that wraps the basic methods for setting custom HTTP Body's.
var _ RequestBody[RequestBuilder] = (*RequestBodyBuilder)(nil)

// RequestBodyBuilder serves as the main entry point for configuring RequestBody.
type RequestBodyBuilder struct {
	parentBuilder *RequestBuilder
}

// Body returns a new RequestBodyBuilder for setting custom HTTP Body's.
func (b *RequestBuilder) Body() *RequestBodyBuilder {
	return &RequestBodyBuilder{parentBuilder: b}
}

// AsReader sets the body as IO Reader.
func (b *RequestBodyBuilder) AsReader(body io.Reader) *RequestBuilder {
	b.parentBuilder.request.body = body
	return b.parentBuilder
}

// AsString sets the body as string.
func (b *RequestBodyBuilder) AsString(body string) *RequestBuilder {
	b.parentBuilder.request.body = bytes.NewBufferString(body)
	return b.parentBuilder
}

// AsJSON sets the body as JSON.
func (b *RequestBodyBuilder) AsJSON(obj interface{}) *RequestBuilder {
	// Marshal JSON
	bodyBytes, err := json.Marshal(obj)
	if err != nil {
		b.parentBuilder.request.validations = append(b.parentBuilder.request.validations, errors.Join(errors.New(constant.ErrMsgMarshalJSON), err))
		return b.parentBuilder
	}
	// Set body
	b.parentBuilder.request.body = bytes.NewBuffer(bodyBytes)
	return b.parentBuilder
}
