package fastshot

import "net/http"

type ResponseFluentCookie struct {
	cookies []*http.Cookie
}

func (r *Response) Cookie() *ResponseFluentCookie {
	return r.cookie
}

func (c *ResponseFluentCookie) GetAll() []*http.Cookie {
	return c.cookies
}
