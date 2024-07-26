package fastshot

import "time"

// BuilderRequestRetry is the interface that wraps the basic methods for setting Request maxAttempts.
var _ BuilderRequestRetry[RequestBuilder] = (*RequestRetryBuilder)(nil)

// RequestRetryBuilder serves as the main entry point for configuring Request maxAttempts.
type RequestRetryBuilder struct {
	parentBuilder *RequestBuilder
	requestConfig *RequestConfigBase
}

// Retry returns a new RequestRetryBuilder for setting custom HTTP CookiesWrapper.
func (b *RequestBuilder) Retry() *RequestRetryBuilder {
	return &RequestRetryBuilder{
		parentBuilder: b,
		requestConfig: b.request.config,
	}
}

// SetConstantBackoff sets the retry interval and maximum attempts.
func (r RequestRetryBuilder) SetConstantBackoff(interval time.Duration, maxAttempts uint) *RequestBuilder {
	r.requestConfig.RetryConfig().SetInterval(interval)
	r.requestConfig.RetryConfig().SetMaxAttempts(maxAttempts)
	r.requestConfig.RetryConfig().SetBackoffRate(1)
	r.requestConfig.RetryConfig().SetJitterStrategy(JitterStrategyNone)
	return r.parentBuilder
}

// SetConstantBackoffWithJitter sets the retry interval and maximum attempts.
func (r RequestRetryBuilder) SetConstantBackoffWithJitter(interval time.Duration, maxAttempts uint) *RequestBuilder {
	r.requestConfig.RetryConfig().SetInterval(interval)
	r.requestConfig.RetryConfig().SetMaxAttempts(maxAttempts)
	r.requestConfig.RetryConfig().SetBackoffRate(1)
	r.requestConfig.RetryConfig().SetJitterStrategy(JitterStrategyFull)
	return r.parentBuilder
}

// SetExponentialBackoff sets the retry interval, maximum attempts, and backoff rate.
func (r RequestRetryBuilder) SetExponentialBackoff(interval time.Duration, maxAttempts uint, backoffRate float64) *RequestBuilder {
	r.requestConfig.RetryConfig().SetInterval(interval)
	r.requestConfig.RetryConfig().SetMaxAttempts(maxAttempts)
	r.requestConfig.RetryConfig().SetBackoffRate(backoffRate)
	r.requestConfig.RetryConfig().SetJitterStrategy(JitterStrategyNone)
	return r.parentBuilder
}

// SetExponentialBackoffWithJitter sets the retry interval, maximum attempts, and backoff rate.
func (r RequestRetryBuilder) SetExponentialBackoffWithJitter(interval time.Duration, maxAttempts uint, backoffRate float64) *RequestBuilder {
	r.requestConfig.RetryConfig().SetInterval(interval)
	r.requestConfig.RetryConfig().SetMaxAttempts(maxAttempts)
	r.requestConfig.RetryConfig().SetBackoffRate(backoffRate)
	r.requestConfig.RetryConfig().SetJitterStrategy(JitterStrategyFull)
	return r.parentBuilder
}

// WithRetryCondition sets the retry condition for the request.
func (r RequestRetryBuilder) WithRetryCondition(shouldRetry func(response *Response) bool) *RequestBuilder {
	r.requestConfig.RetryConfig().SetShouldRetry(shouldRetry)
	return r.parentBuilder
}

// WithMaxDelay sets the maximum delay for the request.
func (r RequestRetryBuilder) WithMaxDelay(duration time.Duration) *RequestBuilder {
	r.requestConfig.RetryConfig().SetMaxDelay(duration)
	return r.parentBuilder
}
