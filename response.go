package fastshot

import (
	"io"
	"net/http"
)

type Response struct {
	Request     *Request
	RawResponse *http.Response
}

func (r *Response) Status() string {
	if r.RawResponse == nil {
		return ""
	}
	return r.RawResponse.Status
}

func (r *Response) StatusCode() int {
	if r.RawResponse == nil {
		return 0
	}
	return r.RawResponse.StatusCode
}

func (r *Response) RawBody() io.ReadCloser {
	if r.RawResponse == nil {
		return nil
	}
	return r.RawResponse.Body
}

func (r *Response) Is1xxInformational() bool {
	if r.RawResponse == nil {
		return false
	}
	return r.RawResponse.StatusCode >= 100 && r.RawResponse.StatusCode < 200
}

func (r *Response) Is2xxSuccessful() bool {
	if r.RawResponse == nil {
		return false
	}
	return r.RawResponse.StatusCode >= 200 && r.RawResponse.StatusCode < 300
}

func (r *Response) Is3xxRedirection() bool {
	if r.RawResponse == nil {
		return false
	}
	return r.RawResponse.StatusCode >= 300 && r.RawResponse.StatusCode < 400
}

func (r *Response) Is4xxClientError() bool {
	if r.RawResponse == nil {
		return false
	}
	return r.RawResponse.StatusCode >= 400 && r.RawResponse.StatusCode < 500
}

func (r *Response) Is5xxServerError() bool {
	if r.RawResponse == nil {
		return false
	}
	return r.RawResponse.StatusCode >= 500 && r.RawResponse.StatusCode < 600
}

func (r *Response) IsError() bool {
	return r.Is4xxClientError() || r.Is5xxServerError()
}
