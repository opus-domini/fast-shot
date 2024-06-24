package fastshot

import "net/http"

type ResponseFluentRequest struct {
	request *http.Request
}

func (r *Response) Request() *ResponseFluentRequest {
	return r.request
}

func (r *ResponseFluentRequest) Raw() *http.Request {
	return r.request
}

func (r *ResponseFluentRequest) Method() string {
	return r.request.Method
}

func (r *ResponseFluentRequest) URL() string {
	return r.request.URL.String()
}

func (r *ResponseFluentRequest) Headers() http.Header {
	return r.request.Header
}
