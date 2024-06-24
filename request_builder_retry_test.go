package fastshot

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestRequestRetryBuilder(t *testing.T) {
	tests := []struct {
		name           string
		method         func(*RequestBuilder) *RequestBuilder
		expectedConfig func(*RetryConfig) bool
	}{
		{
			name: "Set constant backoff",
			method: func(rb *RequestBuilder) *RequestBuilder {
				return rb.Retry().SetConstantBackoff(time.Second, 3)
			},
			expectedConfig: func(rc *RetryConfig) bool {
				return rc.Interval() == time.Second &&
					rc.MaxAttempts() == 3 &&
					rc.BackoffRate() == 1 &&
					rc.JitterStrategy() == JitterStrategyNone
			},
		},
		{
			name: "Set constant backoff with jitter",
			method: func(rb *RequestBuilder) *RequestBuilder {
				return rb.Retry().SetConstantBackoffWithJitter(time.Second, 3)
			},
			expectedConfig: func(rc *RetryConfig) bool {
				return rc.Interval() == time.Second &&
					rc.MaxAttempts() == 3 &&
					rc.BackoffRate() == 1 &&
					rc.JitterStrategy() == JitterStrategyFull
			},
		},
		{
			name: "Set exponential backoff",
			method: func(rb *RequestBuilder) *RequestBuilder {
				return rb.Retry().SetExponentialBackoff(time.Second, 3, 2.0)
			},
			expectedConfig: func(rc *RetryConfig) bool {
				return rc.Interval() == time.Second &&
					rc.MaxAttempts() == 3 &&
					rc.BackoffRate() == 2.0 &&
					rc.JitterStrategy() == JitterStrategyNone
			},
		},
		{
			name: "Set exponential backoff with jitter",
			method: func(rb *RequestBuilder) *RequestBuilder {
				return rb.Retry().SetExponentialBackoffWithJitter(time.Second, 3, 2.0)
			},
			expectedConfig: func(rc *RetryConfig) bool {
				return rc.Interval() == time.Second &&
					rc.MaxAttempts() == 3 &&
					rc.BackoffRate() == 2.0 &&
					rc.JitterStrategy() == JitterStrategyFull
			},
		},
		{
			name: "Set retry condition",
			method: func(rb *RequestBuilder) *RequestBuilder {
				return rb.Retry().WithRetryCondition(func(response *Response) bool {
					return response.Status().Is5xxServerError()
				})
			},
			expectedConfig: func(rc *RetryConfig) bool {
				return rc.ShouldRetry() != nil
			},
		},
		{
			name: "Set max delay",
			method: func(rb *RequestBuilder) *RequestBuilder {
				return rb.Retry().WithMaxDelay(5 * time.Second)
			},
			expectedConfig: func(rc *RetryConfig) bool {
				return rc.MaxDelay() != nil && *rc.MaxDelay() == 5*time.Second
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			rb := &RequestBuilder{
				request: &Request{
					config: newRequestConfigBase("", ""),
				},
			}

			// Act
			result := tt.method(rb)

			// Assert
			assert.Equal(t, rb, result)
			assert.True(t, tt.expectedConfig(rb.request.config.RetryConfig()))
		})
	}
}
