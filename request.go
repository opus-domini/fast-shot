package fastshot

import (
	"context"
	"errors"
	"fmt"
	"math"
	"math/rand/v2"
	"net/http"
	"net/url"
	"time"

	"github.com/opus-domini/fast-shot/constant"
)

type Request struct {
	client ClientConfig
	config *RequestConfigBase
}

type RequestBuilder struct {
	request *Request
}

func newRequest(client ClientConfig, method, path string) *RequestBuilder {
	return &RequestBuilder{
		request: &Request{
			client: client,
			config: &RequestConfigBase{
				ctx: context.Background(),
				httpHeader: &DefaultHttpHeader{
					header: &http.Header{},
				},
				httpCookies: &DefaultHttpCookies{
					cookies: []*http.Cookie{},
				},
				method:      method,
				path:        path,
				queryParams: url.Values{},
				validations: &DefaultValidations{
					validations: []error{},
				},
				retryConfig: &RetryConfig{
					shouldRetry:    func(response Response) bool { return response.IsError() },
					interval:       1 * time.Second,
					backoffRate:    2.0,
					jitterStrategy: JitterStrategyNone,
				},
			},
		},
	}
}

func (b *RequestBuilder) createFullURL() *url.URL {
	// Parse base URL and path
	fullURL := b.request.client.BaseURL().JoinPath(b.request.config.Path())

	// Add query params
	query := fullURL.Query()
	for param, values := range b.request.config.QueryParams() {
		for _, value := range values {
			query.Add(param, value)
		}
	}
	fullURL.RawQuery = query.Encode()

	return fullURL
}

func (b *RequestBuilder) createHTTPRequest() (*http.Request, error) {
	// Create full URL
	fullURL := b.createFullURL()

	// Create Http Request with context
	request, err := http.NewRequestWithContext(
		b.request.config.Context(),
		b.request.config.Method(),
		fullURL.String(),
		b.request.config.Body(),
	)
	if err != nil {
		return nil, err
	}

	// Add client httpCookies
	for _, cookie := range b.request.client.Cookies().Unwrap() {
		request.AddCookie(cookie)
	}

	// Add request httpCookies
	for _, cookie := range b.request.config.Cookies().Unwrap() {
		request.AddCookie(cookie)
	}

	// Clone and attach client httpHeader
	request.Header = http.Header.Clone(*b.request.client.Header().Unwrap())

	// Add Request Headers
	for key, values := range *b.request.config.Header().Unwrap() {
		for _, value := range values {
			request.Header.Add(key, value)
		}
	}

	return request, nil
}

func (b *RequestBuilder) execute(req *http.Request) (Response, error) {
	// Execute request
	response, err := b.request.client.HttpClient().Do(req)

	return Response{Request: b.request, RawResponse: response}, err
}

func (b *RequestBuilder) executeWithRetry() (Response, error) {
	var config = b.request.config.RetryConfig()
	var errExecution error
	var errAttempts []error
	var response Response

	for attempt := uint(0); attempt < config.MaxAttempts(); attempt++ {
		// For each unique attempt, create a new request instance to avoid side effects from previous attempts (e.g. request body)
		req, err := b.createHTTPRequest()
		if err != nil {
			return Response{}, errors.Join(errors.New(constant.ErrMsgCreateRequest), err)
		}
		// Execute request
		response, errExecution = b.execute(req)
		// Check for errors
		if errExecution == nil {
			if !config.ShouldRetry()(response) {
				return response, nil
			}
			errExecution = errors.New(response.StatusText())
		}
		// Append error
		errAttempts = append(errAttempts, fmt.Errorf("attempt %d: %w", attempt+1, errExecution))
		// Delay before retry
		delay := b.calculateRetryDelay(attempt)
		time.Sleep(delay)
	}

	return response,
		fmt.Errorf(
			"request failed after %d attempts: %w",
			config.MaxAttempts(),
			errors.Join(errAttempts...),
		)
}

func (b *RequestBuilder) calculateRetryDelay(attempt uint) time.Duration {
	config := b.request.config.RetryConfig()
	delay := float64(config.Interval()) * math.Pow(config.BackoffRate(), float64(attempt))

	if config.MaxDelay() != nil {
		delay = math.Min(delay, float64(*config.MaxDelay()))
	}

	if config.JitterStrategy() == JitterStrategyFull {
		delay = rand.Float64() * delay
	}

	return time.Duration(delay)
}

func (b *RequestBuilder) Send() (Response, error) {
	// Check for client validation errors
	if err := errors.Join(b.request.client.Validations().Unwrap()...); err != nil {
		return Response{}, errors.Join(errors.New(constant.ErrMsgClientValidation), err)
	}

	// Check for request validation errors
	if err := errors.Join(b.request.config.Validations().Unwrap()...); err != nil {
		return Response{}, errors.Join(errors.New(constant.ErrMsgRequestValidation), err)
	}

	// Check if maxAttempts are enabled
	if b.request.config.RetryConfig() != nil && b.request.config.RetryConfig().MaxAttempts() > 1 {
		return b.executeWithRetry()
	}

	// Create request
	req, err := b.createHTTPRequest()
	if err != nil {
		return Response{}, errors.Join(errors.New(constant.ErrMsgCreateRequest), err)
	}

	// Execute the request
	return b.execute(req)
}
