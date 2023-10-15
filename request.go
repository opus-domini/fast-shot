package fastshot

import (
	"context"
	"errors"
	"fmt"
	"github.com/opus-domini/fast-shot/constant"
	"io"
	"net/http"
	"net/url"
	"time"
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

type RequestBuilder struct {
	request *Request
}

func newRequest(client *Client, method, path string) *RequestBuilder {
	return &RequestBuilder{
		request: &Request{
			ctx:         context.Background(),
			client:      client,
			httpHeader:  &http.Header{},
			httpCookies: []*http.Cookie{},
			method:      method,
			path:        path,
			queryParams: url.Values{},
		},
	}
}

func (b *RequestBuilder) createFullURL() (*url.URL, error) {
	// Parse base URL and path
	fullURL, err := url.Parse(b.request.client.baseURL + b.request.path)
	if err != nil {
		return nil, errors.Join(errors.New(constant.ErrMsgParseURL), err)
	}

	// Add query params
	query := fullURL.Query()
	for param, values := range b.request.queryParams {
		for _, value := range values {
			query.Add(param, value)
		}
	}
	fullURL.RawQuery = query.Encode()

	return fullURL, nil
}

func (b *RequestBuilder) createHTTPRequest() (*http.Request, error) {
	// Create full URL
	fullURL, err := b.createFullURL()
	if err != nil {
		return nil, err
	}

	// Create Http Request with context
	request, err := http.NewRequestWithContext(b.request.ctx, b.request.method, fullURL.String(), b.request.body)
	if err != nil {
		return nil, err
	}

	// Add client httpCookies
	for _, cookie := range b.request.client.httpCookies {
		request.AddCookie(cookie)
	}

	// Add request httpCookies
	for _, cookie := range b.request.httpCookies {
		request.AddCookie(cookie)
	}

	// Clone and attach client httpHeader
	request.Header = http.Header.Clone(*b.request.client.httpHeader)

	// Add Request Headers
	for key, values := range *b.request.httpHeader {
		for _, value := range values {
			request.Header.Add(key, value)
		}
	}

	return request, nil
}

func (b *RequestBuilder) Send() (Response, error) {
	// Check for validation errors
	if err := errors.Join(b.request.validations...); err != nil {
		return Response{}, errors.Join(errors.New(constant.ErrMsgValidation), err)
	}

	// Create request
	req, err := b.createHTTPRequest()
	if err != nil {
		return Response{}, errors.Join(errors.New(constant.ErrMsgCreateRequest), err)
	}

	var response *http.Response
	var errAttempts []error

	for i := 0; i <= b.request.retries; i++ {
		// Execute request
		response, err = b.request.client.httpClient.Do(req)
		// Check for errors
		resp := Response{Request: b.request, RawResponse: response}
		if err == nil && !resp.IsError() {
			return resp, nil
		}
		// Append error
		errAttempts = append(errAttempts, fmt.Errorf("attempt %d: %w", i+1, err))
		// Delay before retry (if applicable)
		if i < b.request.retries {
			time.Sleep(b.request.retryInterval)
		}
	}

	return Response{Request: b.request, RawResponse: response}, fmt.Errorf("request failed after %d attempts: %w", b.request.retries+1, errors.Join(errAttempts...))
}
