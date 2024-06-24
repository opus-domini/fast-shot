package fastshot

import (
	"net/http"
)

type Response struct {
	rawResponse *http.Response
	// Fluent API
	body    *ResponseFluentBody
	cookie  *ResponseFluentCookie
	header  *ResponseFluentHeader
	request *ResponseFluentRequest
	status  *ResponseFluentStatus
}

func (r *Response) Raw() *http.Response {
	return r.rawResponse
}

func newResponse(response *http.Response) *Response {
	return &Response{
		rawResponse: response,
		// Fluent API
		body: &ResponseFluentBody{
			newUnbufferedBody(response.Body),
		},
		cookie: &ResponseFluentCookie{
			response.Cookies(),
		},
		header: &ResponseFluentHeader{
			response.Header,
		},
		request: &ResponseFluentRequest{
			response.Request,
		},
		status: &ResponseFluentStatus{
			response,
		},
	}
}
