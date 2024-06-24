package fastshot

import (
	"fmt"
	"net/http"
)

type ResponseFluentStatus struct {
	response *http.Response
}

func (r *Response) Status() *ResponseFluentStatus {
	return r.status
}

func (r *ResponseFluentStatus) Code() int {
	return r.response.StatusCode
}

func (r *ResponseFluentStatus) Text() string {
	return fmt.Sprintf(
		"[%d] %s",
		r.response.StatusCode,
		http.StatusText(r.response.StatusCode),
	)
}

func (r *ResponseFluentStatus) Is1xxInformational() bool {
	return r.response.StatusCode >= 100 && r.response.StatusCode < 200
}

func (r *ResponseFluentStatus) Is2xxSuccessful() bool {
	return r.response.StatusCode >= 200 && r.response.StatusCode < 300
}

func (r *ResponseFluentStatus) Is3xxRedirection() bool {
	return r.response.StatusCode >= 300 && r.response.StatusCode < 400
}

func (r *ResponseFluentStatus) Is4xxClientError() bool {
	return r.response.StatusCode >= 400 && r.response.StatusCode < 500
}

func (r *ResponseFluentStatus) Is5xxServerError() bool {
	return r.response.StatusCode >= 500 && r.response.StatusCode < 600
}

func (r *ResponseFluentStatus) IsOK() bool {
	return r.Code() == http.StatusOK
}

func (r *ResponseFluentStatus) IsNotFound() bool {
	return r.Code() == http.StatusNotFound
}

func (r *ResponseFluentStatus) IsUnauthorized() bool {
	return r.Code() == http.StatusUnauthorized
}

func (r *ResponseFluentStatus) IsForbidden() bool {
	return r.Code() == http.StatusForbidden
}

func (r *ResponseFluentStatus) IsError() bool {
	return r.Is4xxClientError() || r.Is5xxServerError()
}
