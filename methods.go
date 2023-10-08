package fastshot

import "net/http"

// GET is a shortcut for NewRequest(http.MethodGet, path).
func (c *Client) GET(path string) *Request {
	return newRequest(c, http.MethodGet, path)
}

// POST is a shortcut for NewRequest(http.MethodPost, path).
func (c *Client) POST(path string) *Request {
	return newRequest(c, http.MethodPost, path)
}

// PUT is a shortcut for NewRequest(http.MethodPut, path).
func (c *Client) PUT(path string) *Request {
	return newRequest(c, http.MethodPut, path)
}

// DELETE is a shortcut for NewRequest(http.MethodDelete, path).
func (c *Client) DELETE(path string) *Request {
	return newRequest(c, http.MethodDelete, path)
}

// PATCH is a shortcut for NewRequest(http.MethodPatch, path).
func (c *Client) PATCH(path string) *Request {
	return newRequest(c, http.MethodPatch, path)
}

// HEAD is a shortcut for NewRequest(http.MethodHead, path).
func (c *Client) HEAD(path string) *Request {
	return newRequest(c, http.MethodHead, path)
}

// OPTIONS is a shortcut for NewRequest(http.MethodOptions, path).
func (c *Client) OPTIONS(path string) *Request {
	return newRequest(c, http.MethodOptions, path)
}
