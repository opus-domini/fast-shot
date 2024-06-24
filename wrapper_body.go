package fastshot

import (
	"bytes"
	"encoding/json"
	"io"
)

// DefaultBody implements BodyWrapper interface and provides a default HTTP context.
var _ BodyWrapper = (*DefaultBody)(nil)

// DefaultBody implements ContextWrapper interface and provides a default HTTP context.
type DefaultBody struct {
	body io.Reader
}

// Unwrap will return the underlying body
func (c *DefaultBody) Unwrap() io.Reader {
	return c.body
}

// Set will set the body
func (c *DefaultBody) Set(body io.Reader) {
	c.body = body
}

// SetAsJSON will set the body as JSON
func (c *DefaultBody) SetAsJSON(obj interface{}) error {
	bodyBytes, err := json.Marshal(obj)
	if err != nil {
		return err
	}
	c.body = bytes.NewBuffer(bodyBytes)
	return nil
}

// newDefaultBody initializes a new DefaultBody with a given body.
func newDefaultBody() *DefaultBody {
	return &DefaultBody{}
}
