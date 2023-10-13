package fastshot

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const (
	ErrMsgValidation       = "invalid request attributes"
	ErrMsgCreateRequest    = "failed to create request"
	ErrMsgMarshalJSON      = "failed to marshal JSON"
	ErrMsgParseQueryString = "failed to parse query string"
	ErrMsgParseURL         = "failed to parse URL"
)

type Request struct {
	ctx           context.Context
	client        *Client
	httpHeader    *http.Header
	httpCookies   []*http.Cookie
	method        string
	path          string
	queryParams   url.Values
	body          io.Reader
	validations   []error
	retries       int
	retryInterval time.Duration
}

func newRequest(client *Client, method, path string) *Request {
	return &Request{
		ctx:         context.Background(),
		client:      client,
		httpHeader:  &http.Header{},
		httpCookies: []*http.Cookie{},
		method:      method,
		path:        path,
		queryParams: url.Values{},
	}
}

func (r *Request) SetContext(ctx context.Context) *Request {
	r.ctx = ctx
	return r
}

func (r *Request) SetHeader(key, value string) *Request {
	r.httpHeader.Set(key, value)
	return r
}

func (r *Request) SetHeaders(headers map[string]string) *Request {
	for key, value := range headers {
		r.SetHeader(key, value)
	}
	return r
}

func (r *Request) AddCookie(cookie *http.Cookie) *Request {
	r.httpCookies = append(r.httpCookies, cookie)
	return r
}

func (r *Request) AddQueryParam(param, value string) *Request {
	r.queryParams.Add(param, value)
	return r
}

func (r *Request) SetQueryParam(param, value string) *Request {
	r.queryParams.Set(param, value)
	return r
}

func (r *Request) SetQueryParams(params map[string]string) *Request {
	for param, value := range params {
		r.SetQueryParam(param, value)
	}
	return r
}

func (r *Request) SetQueryString(query string) *Request {
	// Parse query string
	queryParams, err := url.ParseQuery(strings.TrimSpace(query))
	if err != nil {
		r.validations = append(r.validations, errors.Join(errors.New(ErrMsgParseQueryString), err))
		return r
	}
	// Set query params
	for param, values := range queryParams {
		for _, value := range values {
			r.AddQueryParam(param, value)
		}
	}
	return r
}

func (r *Request) SetRetry(retries int, retryInterval time.Duration) *Request {
	r.retries = retries
	r.retryInterval = retryInterval
	return r
}

func (r *Request) Body(body io.Reader) *Request {
	r.body = body
	return r
}

func (r *Request) BodyJSON(obj interface{}) *Request {
	// Marshal JSON
	bodyBytes, err := json.Marshal(obj)
	if err != nil {
		r.validations = append(r.validations, errors.Join(errors.New(ErrMsgMarshalJSON), err))
		return r
	}
	// Set body
	r.body = bytes.NewBuffer(bodyBytes)
	return r
}

func (r *Request) createFullURL() (*url.URL, error) {
	// Parse base URL and path
	fullURL, err := url.Parse(r.client.baseURL + r.path)
	if err != nil {
		return nil, errors.Join(errors.New(ErrMsgParseURL), err)
	}

	// Add query params
	query := fullURL.Query()
	for param, values := range r.queryParams {
		for _, value := range values {
			query.Add(param, value)
		}
	}
	fullURL.RawQuery = query.Encode()

	return fullURL, nil
}

func (r *Request) createHTTPRequest() (*http.Request, error) {
	// Create full URL
	fullURL, err := r.createFullURL()
	if err != nil {
		return nil, err
	}

	// Create Http Request with context
	request, err := http.NewRequestWithContext(r.ctx, r.method, fullURL.String(), r.body)
	if err != nil {
		return nil, err
	}

	// Add client httpCookies
	for _, cookie := range r.client.httpCookies {
		request.AddCookie(cookie)
	}

	// Add request httpCookies
	for _, cookie := range r.httpCookies {
		request.AddCookie(cookie)
	}

	// Clone and attach client httpHeader
	request.Header = http.Header.Clone(*r.client.httpHeader)

	// Add Request Headers
	for key, values := range *r.httpHeader {
		for _, value := range values {
			request.Header.Add(key, value)
		}
	}

	return request, nil
}

func (r *Request) Send() (Response, error) {
	// Check for validation errors
	if err := errors.Join(r.validations...); err != nil {
		return Response{}, errors.Join(errors.New(ErrMsgValidation), err)
	}

	// Create request
	req, err := r.createHTTPRequest()
	if err != nil {
		return Response{}, errors.Join(errors.New(ErrMsgCreateRequest), err)
	}

	var response *http.Response
	var errAttempts []error

	for i := 0; i <= r.retries; i++ {
		// Execute request
		response, err = r.client.httpClient.Do(req)
		// Check for errors
		resp := Response{Request: r, RawResponse: response}
		if err == nil && !resp.IsError() {
			return resp, nil
		}
		// Append error
		errAttempts = append(errAttempts, fmt.Errorf("attempt %d: %w", i+1, err))
		// Delay before retry (if applicable)
		if i < r.retries {
			time.Sleep(r.retryInterval)
		}
	}

	return Response{Request: r, RawResponse: response}, fmt.Errorf("request failed after %d attempts: %w", r.retries+1, errors.Join(errAttempts...))
}
