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
)

type Request struct {
	client      *Client
	ctx         context.Context
	method      string
	path        string
	headers     http.Header
	queryParams url.Values
	body        io.Reader
	validations []error
}

func newRequest(client *Client, method, path string) *Request {
	return &Request{
		client:      client,
		ctx:         context.Background(),
		method:      method,
		path:        path,
		headers:     http.Header{},
		queryParams: url.Values{},
	}
}

func (r *Request) SetContext(ctx context.Context) *Request {
	r.ctx = ctx
	return r
}

func (r *Request) SetHeader(key, value string) *Request {
	r.headers.Set(key, value)
	return r
}

func (r *Request) SetHeaders(headers map[string]string) *Request {
	for key, value := range headers {
		r.SetHeader(key, value)
	}
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
		r.validations = append(r.validations, fmt.Errorf("failed to parse query string: %w", err))
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

func (r *Request) Body(body io.Reader) *Request {
	r.body = body
	return r
}

func (r *Request) BodyJSON(obj interface{}) *Request {
	// Marshal JSON
	bodyBytes, err := json.Marshal(obj)
	if err != nil {
		r.validations = append(r.validations, fmt.Errorf("failed to marshal JSON: %w", err))
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
		return nil, fmt.Errorf("failed to parse URL: %w", err)
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

func (r *Request) validate() error {
	if len(r.validations) > 0 {
		return errors.Join(r.validations...)
	}
	return nil
}

func (r *Request) Send() (Response, error) {
	// Check for validation errors
	if err := r.validate(); err != nil {
		return Response{}, err
	}

	// Create full URL
	fullURL, err := r.createFullURL()
	if err != nil {
		return Response{}, err
	}

	// Create request
	req, err := http.NewRequestWithContext(r.ctx, r.method, fullURL.String(), r.body)
	if err != nil {
		return Response{}, err
	}

	// Add cookies
	for _, cookie := range r.client.httpCookies {
		req.AddCookie(cookie)
	}

	// Clone and attach client headers
	req.Header = http.Header.Clone(*r.client.httpHeader)

	// Add Request Headers
	for key, values := range r.headers {
		for _, value := range values {
			req.Header.Add(key, value)
		}
	}

	// Execute request
	response, err := r.client.httpClient.Do(req)

	return Response{Request: r, RawResponse: response}, nil
}
