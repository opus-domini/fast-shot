package fastshot

import (
	"io"
)

type ResponseFluentBody struct {
	body BodyWrapper
}

func (r *Response) Body() *ResponseFluentBody {
	return r.body
}

func (b *ResponseFluentBody) Raw() io.ReadCloser {
	return b.body
}

func (b *ResponseFluentBody) Close() {
	_ = b.body.Close()
}

func (b *ResponseFluentBody) CloseErr() error {
	return b.body.Close()
}

func (b *ResponseFluentBody) AsBytes() ([]byte, error) {
	defer b.Close()

	data, err := io.ReadAll(b.body)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (b *ResponseFluentBody) AsString() (string, error) {
	defer b.Close()

	return b.body.ReadAsString()
}

func (b *ResponseFluentBody) AsJSON(v any) error {
	defer b.Close()

	return b.body.ReadAsJSON(v)
}

// AsJSONOf decodes the response body as JSON into a new value of type T,
// closing the body afterwards.
//
// It is the generic form of AsJSON, returning the decoded value directly
// instead of requiring a pre-allocated destination:
//
//	type User struct {
//		Name string `json:"name"`
//	}
//	user, err := resp.Body().AsJSONOf[User]()
func (b *ResponseFluentBody) AsJSONOf[T any]() (T, error) {
	var value T
	err := b.AsJSON(&value)
	return value, err
}

func (b *ResponseFluentBody) AsXML(v any) error {
	defer b.Close()

	return b.body.ReadAsXML(v)
}

// AsXMLOf decodes the response body as XML into a new value of type T,
// closing the body afterwards.
//
// It is the generic form of AsXML, returning the decoded value directly
// instead of requiring a pre-allocated destination:
//
//	type User struct {
//		Name string `xml:"name"`
//	}
//	user, err := resp.Body().AsXMLOf[User]()
func (b *ResponseFluentBody) AsXMLOf[T any]() (T, error) {
	var value T
	err := b.AsXML(&value)
	return value, err
}
