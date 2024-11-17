package fastshot

import (
	"bytes"
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

func (b *ResponseFluentBody) AsBytes() ([]byte, error) {
	defer b.Close()

	buf := new(bytes.Buffer)
	_, err := buf.ReadFrom(b.body)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func (b *ResponseFluentBody) AsString() (string, error) {
	defer b.Close()

	return b.body.ReadAsString()
}

func (b *ResponseFluentBody) AsJSON(v interface{}) error {
	defer b.Close()

	return b.body.ReadAsJSON(v)
}

func (b *ResponseFluentBody) AsXML(v interface{}) error {
	defer b.Close()

	return b.body.ReadAsXML(v)
}
