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

// AsJSON decodes the response body as JSON into a value of type T,
// closing the body afterwards.
//
// The type parameter removes the need for a pre-allocated destination:
//
//	type User struct {
//		Name string `json:"name"`
//	}
//	user, err := resp.Body().AsJSON[User]()
func (b *ResponseFluentBody) AsJSON[T any]() (T, error) {
	var value T
	defer b.Close()

	err := b.body.ReadAsJSON(&value)
	return value, err
}

// AsXML decodes the response body as XML into a value of type T,
// closing the body afterwards.
//
// The type parameter removes the need for a pre-allocated destination:
//
//	type User struct {
//		Name string `xml:"name"`
//	}
//	user, err := resp.Body().AsXML[User]()
func (b *ResponseFluentBody) AsXML[T any]() (T, error) {
	var value T
	defer b.Close()

	err := b.body.ReadAsXML(&value)
	return value, err
}
