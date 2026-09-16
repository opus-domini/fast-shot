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

	"github.com/opus-domini/fast-shot/constant/method"
)

type Request struct {
	client ClientConfig
	config *RequestConfigBase
}

type RequestBuilder struct {
	request *Request
}

func newRequest(client ClientConfig, method method.Type, path string) *RequestBuilder {
	return &RequestBuilder{
		request: &Request{
			client: client,
			config: newRequestConfigBase(method, path, client.JSONCodec()),
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
		b.request.config.Context().Unwrap(),
		b.request.config.Method().String(),
		fullURL.String(),
		b.request.config.Body().Unwrap(),
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

	// Add Client Headers
	for key, values := range *b.request.client.Header().Unwrap() {
		for _, value := range values {
			request.Header.Add(key, value)
		}
	}

	// Add Request Headers
	for key, values := range *b.request.config.Header().Unwrap() {
		for _, value := range values {
			request.Header.Add(key, value)
		}
	}

	return request, nil
}

func (b *RequestBuilder) runBeforeRequestHooks(req *http.Request) error {
	for _, hook := range b.request.client.BeforeRequestHooks() {
		if err := hook(req); err != nil {
			return err
		}
	}
	for _, hook := range b.request.config.BeforeRequestHooks() {
		if err := hook(req); err != nil {
			return err
		}
	}
	return nil
}

func (b *RequestBuilder) runAfterResponseHooks(req *http.Request, resp *http.Response) {
	//nolint:bodyclose // False positive: iterating hook functions, not handling a response body.
	for _, hook := range b.request.client.AfterResponseHooks() {
		hook(req, resp)
	}
	//nolint:bodyclose // False positive: iterating hook functions, not handling a response body.
	for _, hook := range b.request.config.AfterResponseHooks() {
		hook(req, resp)
	}
}

func (b *RequestBuilder) execute(request *http.Request) (*Response, error) {
	// Run before-request hooks
	if err := b.runBeforeRequestHooks(request); err != nil {
		return nil, fmt.Errorf("%w: %w", ErrBeforeRequestHook, err)
	}

	// Execute request
	//nolint:bodyclose // Response body is closed by the caller via Response APIs.
	response, err := b.request.client.HttpClient().Do(request)
	if err != nil {
		return nil, err
	}

	// Run after-response hooks
	b.runAfterResponseHooks(request, response)

	return newResponse(response, b.request.client.JSONCodec()), nil
}

// requestForAttempt returns the request to send on the given attempt.
// Retries get a fresh clone with a reset body (when replayable) so consumed
// request bodies are not silently sent empty on subsequent attempts.
func (b *RequestBuilder) requestForAttempt(req *http.Request, attempt uint) (*http.Request, error) {
	if attempt == 0 || req.GetBody == nil {
		return req, nil
	}
	body, err := req.GetBody()
	if err != nil {
		return nil, err
	}
	attemptReq := req.Clone(req.Context())
	attemptReq.Body = body
	return attemptReq, nil
}

// sleepBackoff waits for the retry delay, aborting early if the request
// context is canceled.
func (b *RequestBuilder) sleepBackoff(ctx context.Context, delay time.Duration) error {
	if delay <= 0 {
		return nil
	}
	timer := time.NewTimer(delay)
	defer timer.Stop()
	select {
	case <-timer.C:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

func (b *RequestBuilder) executeWithRetry(req *http.Request) (*Response, error) {
	config := b.request.config.RetryConfig()
	ctx := req.Context()
	var errAttempts []error

	for attempt := range config.MaxAttempts() {
		attemptReq, err := b.requestForAttempt(req, attempt)
		if err != nil {
			return nil, newRetryError(config.MaxAttempts(), append(errAttempts, err))
		}

		// Execute request
		response, errExecution := b.execute(attemptReq)
		// Check for errors
		if errExecution == nil {
			if !config.ShouldRetry()(response) {
				return response, nil
			}
			errExecution = errors.New(response.Status().Text())
			// Release the connection of the discarded attempt. Since Go 1.27
			// the body is auto-drained on close, keeping the connection reusable.
			_ = response.Raw().Body.Close()
		}
		// Append error
		errAttempts = append(errAttempts, fmt.Errorf("attempt %d: %w", attempt+1, errExecution))
		// Delay before retry, except after the final attempt
		if attempt == config.MaxAttempts()-1 {
			break
		}
		if err := b.sleepBackoff(ctx, b.calculateRetryDelay(attempt)); err != nil {
			return nil, newRetryError(config.MaxAttempts(), append(errAttempts, err))
		}
	}

	return nil, newRetryError(config.MaxAttempts(), errAttempts)
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

func (b *RequestBuilder) Send() (*Response, error) {
	// Check for client validation errors
	if err := errors.Join(b.request.client.Validations().Unwrap()...); err != nil {
		return nil, fmt.Errorf("%w: %w", ErrClientValidation, err)
	}

	// Check for request validation errors
	if err := errors.Join(b.request.config.Validations().Unwrap()...); err != nil {
		return nil, fmt.Errorf("%w: %w", ErrRequestValidation, err)
	}

	// Create request
	req, err := b.createHTTPRequest()
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrCreateRequest, err)
	}

	// Check if maxAttempts are enabled
	if b.request.config.RetryConfig() != nil && b.request.config.RetryConfig().MaxAttempts() > 1 {
		return b.executeWithRetry(req)
	}

	// Execute the request
	return b.execute(req)
}
