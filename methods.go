package fastshot

import "net/http"

func (c *Client) GET(path string) *Request {
	return newRequest(c, http.MethodGet, path)
}

func (c *Client) POST(path string) *Request {
	return newRequest(c, http.MethodPost, path)
}

func (c *Client) PUT(path string) *Request {
	return newRequest(c, http.MethodPut, path)
}

func (c *Client) DELETE(path string) *Request {
	return newRequest(c, http.MethodDelete, path)
}

func (c *Client) PATCH(path string) *Request {
	return newRequest(c, http.MethodPatch, path)
}

func (c *Client) HEAD(path string) *Request {
	return newRequest(c, http.MethodHead, path)
}

func (c *Client) OPTIONS(path string) *Request {
	return newRequest(c, http.MethodOptions, path)
}
