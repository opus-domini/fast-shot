package fastshot

import "github.com/opus-domini/fast-shot/constant/method"

// GET is a shortcut for NewRequest(c, method.GET, path).
func (c *Client) GET(path string) *RequestBuilder {
	return newRequest(c, method.GET, path)
}

// POST is a shortcut for NewRequest(c, method.POST, path).
func (c *Client) POST(path string) *RequestBuilder {
	return newRequest(c, method.POST, path)
}

// PUT is a shortcut for NewRequest(c, method.PUT, path).
func (c *Client) PUT(path string) *RequestBuilder {
	return newRequest(c, method.PUT, path)
}

// DELETE is a shortcut for NewRequest(c, method.DELETE, path).
func (c *Client) DELETE(path string) *RequestBuilder {
	return newRequest(c, method.DELETE, path)
}

// PATCH is a shortcut for NewRequest(c, method.PATCH, path).
func (c *Client) PATCH(path string) *RequestBuilder {
	return newRequest(c, method.PATCH, path)
}

// HEAD is a shortcut for NewRequest(c, method.HEAD, path).
func (c *Client) HEAD(path string) *RequestBuilder {
	return newRequest(c, method.HEAD, path)
}

// CONNECT is a shortcut for NewRequest(c, method.CONNECT, path).
func (c *Client) CONNECT(path string) *RequestBuilder {
	return newRequest(c, method.CONNECT, path)
}

// OPTIONS is a shortcut for NewRequest(c, method.OPTIONS, path).
func (c *Client) OPTIONS(path string) *RequestBuilder {
	return newRequest(c, method.OPTIONS, path)
}

// TRACE is a shortcut for NewRequest(c, method.TRACE, path).
func (c *Client) TRACE(path string) *RequestBuilder {
	return newRequest(c, method.TRACE, path)
}
