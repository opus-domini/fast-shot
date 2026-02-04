package fastshot

import "context"

// Compile-time check that DefaultContext implements ContextWrapper.
var _ ContextWrapper = (*DefaultContext)(nil)

// DefaultContext implements ContextWrapper interface and provides a default HTTP context.
type DefaultContext struct {
	ctx context.Context
}

// Unwrap will return the underlying context
func (c *DefaultContext) Unwrap() context.Context {
	return c.ctx
}

// Set will set the context
func (c *DefaultContext) Set(ctx context.Context) {
	if ctx != nil {
		c.ctx = ctx
	}
}

// newDefaultContext initializes a new DefaultContext.
func newDefaultContext() *DefaultContext {
	return &DefaultContext{
		ctx: context.Background(),
	}
}
