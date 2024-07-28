package fastshot

import (
	"net/url"
	"time"

	"github.com/opus-domini/fast-shot/constant/method"
)

type (
	// RequestConfigBase encapsulates all configurations for a request.
	RequestConfigBase struct {
		ctx         ContextWrapper
		httpHeader  HeaderWrapper
		httpCookies CookiesWrapper
		method      method.Type
		path        string
		queryParams url.Values
		body        BodyWrapper
		validations ValidationsWrapper
		retryConfig *RetryConfig
	}

	// JitterStrategy represents the strategy for jitter.
	JitterStrategy string

	// RetryConfig represents the configuration for the retry mechanism.
	RetryConfig struct {
		shouldRetry    func(response *Response) bool
		interval       time.Duration
		maxAttempts    uint
		backoffRate    float64
		maxDelay       *time.Duration
		jitterStrategy JitterStrategy
	}
)

const (
	// JitterStrategyNone NONE is the default jitter strategy.
	JitterStrategyNone JitterStrategy = "NONE"
	// JitterStrategyFull FULL is the full jitter strategy.
	JitterStrategyFull JitterStrategy = "FULL"
)

// Context returns the context for the request.
func (c *RequestConfigBase) Context() ContextWrapper {
	return c.ctx
}

// Header returns the header for the request.
func (c *RequestConfigBase) Header() HeaderWrapper {
	return c.httpHeader
}

// Cookies returns the cookies for the request.
func (c *RequestConfigBase) Cookies() CookiesWrapper {
	return c.httpCookies
}

// Method returns the method for the request.
func (c *RequestConfigBase) Method() method.Type {
	return c.method
}

// Path returns the path for the request.
func (c *RequestConfigBase) Path() string {
	return c.path
}

// QueryParams returns the query parameters for the request.
func (c *RequestConfigBase) QueryParams() url.Values {
	return c.queryParams
}

// Body returns the body for the request.
func (c *RequestConfigBase) Body() BodyWrapper {
	return c.body
}

// Validations returns the validations for the request.
func (c *RequestConfigBase) Validations() ValidationsWrapper {
	return c.validations
}

// RetryConfig returns the retry configuration for the request.
func (c *RequestConfigBase) RetryConfig() *RetryConfig {
	return c.retryConfig
}

// ShouldRetry returns the retry condition for the request.
func (c *RetryConfig) ShouldRetry() func(response *Response) bool {
	return c.shouldRetry
}

// SetShouldRetry sets the retry condition for the request.
func (c *RetryConfig) SetShouldRetry(shouldRetry func(response *Response) bool) {
	c.shouldRetry = shouldRetry
}

// Interval returns the retry interval for the request.
func (c *RetryConfig) Interval() time.Duration {
	return c.interval
}

// SetInterval sets the retry interval for the request.
func (c *RetryConfig) SetInterval(duration time.Duration) {
	c.interval = duration
}

// MaxAttempts returns the retry configuration for the request.
func (c *RetryConfig) MaxAttempts() uint {
	return c.maxAttempts
}

// SetMaxAttempts sets the retry configuration for the request.
func (c *RetryConfig) SetMaxAttempts(attempts uint) {
	c.maxAttempts = attempts
}

// BackoffRate returns the retry backoff rate for the request.
func (c *RetryConfig) BackoffRate() float64 {
	return c.backoffRate
}

// SetBackoffRate sets the retry backoff rate for the request.
func (c *RetryConfig) SetBackoffRate(rate float64) {
	c.backoffRate = rate
}

// MaxDelay returns the retry maximum delay for the request.
func (c *RetryConfig) MaxDelay() *time.Duration {
	return c.maxDelay
}

// SetMaxDelay sets the retry maximum delay for the request.
func (c *RetryConfig) SetMaxDelay(duration time.Duration) {
	c.maxDelay = &duration
}

// JitterStrategy returns the retry jitter strategy for the request.
func (c *RetryConfig) JitterStrategy() JitterStrategy {
	return c.jitterStrategy
}

// SetJitterStrategy sets the retry jitter strategy for the request.
func (c *RetryConfig) SetJitterStrategy(strategy JitterStrategy) {
	c.jitterStrategy = strategy
}

// NewRequestConfigBase creates a new request configuration.
func newRequestConfigBase(method method.Type, path string) *RequestConfigBase {
	return &RequestConfigBase{
		ctx:         newDefaultContext(),
		httpHeader:  newDefaultHttpHeader(),
		httpCookies: newDefaultHttpCookies(),
		method:      method,
		path:        path,
		queryParams: url.Values{},
		body:        newBufferedBody(),
		validations: newDefaultValidations(nil),
		retryConfig: &RetryConfig{
			shouldRetry:    func(response *Response) bool { return response.Status().IsError() },
			interval:       1 * time.Second,
			backoffRate:    2.0,
			jitterStrategy: JitterStrategyNone,
		},
	}
}
