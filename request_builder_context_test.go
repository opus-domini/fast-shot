package fastshot

import (
	"context"
	"reflect"
	"testing"
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
			if result != rb {
				t.Errorf("got different builder, want same")
			}
			if tt.ctx == nil {
				if rb.request.config.Context().Unwrap() == nil {
					t.Error("context got nil, want non-nil")
				}
			} else {
				if got := rb.request.config.Context().Unwrap(); !reflect.DeepEqual(got, tt.ctx) {
					t.Errorf("context got %v, want %v", got, tt.ctx)
				}
			}
		})
	}
}
