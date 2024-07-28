package fastshot

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRequestContextBuilder(t *testing.T) {
	type contextKey string

	tests := []struct {
		name string
		ctx  context.Context
	}{
		{
			name: "Set nil context",
			ctx:  nil,
		},
		{
			name: "Set background context",
			ctx:  context.Background(),
		},
		{
			name: "Set context with value",
			ctx:  context.WithValue(context.Background(), contextKey("key"), "value"),
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
			result := rb.Context().Set(tt.ctx)

			// Assert
			assert.Equal(t, rb, result)
			if tt.ctx == nil {
				assert.NotNil(t, rb.request.config.Context().Unwrap())
			} else {
				assert.Equal(t, tt.ctx, rb.request.config.Context().Unwrap())
			}
		})
	}
}
