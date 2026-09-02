package fastshot

import (
	"fmt"
	"io"

	"github.com/opus-domini/fast-shot/constant/header"
)

// RequestBodyBuilder serves as the main entry point for configuring the request body.
type RequestBodyBuilder struct {
	parentBuilder *RequestBuilder
	requestConfig *RequestConfigBase
}

// Body returns a new RequestBodyBuilder for setting the request body.
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
		b.requestConfig.addValidation(fmt.Errorf("%w: %w", ErrSetBody, err))
	}
	return b.parentBuilder
}

// AsString sets the body as string.
func (b *RequestBodyBuilder) AsString(body string) *RequestBuilder {
	err := b.requestConfig.Body().WriteAsString(body)
	if err != nil {
		b.requestConfig.addValidation(fmt.Errorf("%w: %w", ErrSetBody, err))
	}
	return b.parentBuilder
}

// AsJSON sets the body as JSON.
func (b *RequestBodyBuilder) AsJSON(obj any) *RequestBuilder {
	err := b.requestConfig.Body().WriteAsJSON(obj)
	if err != nil {
		b.requestConfig.addValidation(fmt.Errorf("%w: %w", ErrMarshalJSON, err))
	}
	return b.parentBuilder
}

// AsXML sets the body as XML.
func (b *RequestBodyBuilder) AsXML(obj any) *RequestBuilder {
	err := b.requestConfig.Body().WriteAsXML(obj)
	if err != nil {
		b.requestConfig.addValidation(fmt.Errorf("%w: %w", ErrMarshalXML, err))
	}
	return b.parentBuilder
}

// AsFormData sets the body as multipart/form-data.
func (b *RequestBodyBuilder) AsFormData(fields map[string]string) *RequestBuilder {
	contentType, err := b.requestConfig.Body().WriteAsFormData(fields)
	if err != nil {
		b.requestConfig.addValidation(fmt.Errorf("%w: %w", ErrSetBody, err))
	} else {
		b.requestConfig.httpHeader.Set(header.ContentType.String(), contentType)
	}
	return b.parentBuilder
}
