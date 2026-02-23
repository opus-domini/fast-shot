package fastshot

import (
	"net/http"
)

// BuilderHook is the interface that wraps the basic methods for setting request hooks.
var _ BuilderHook[ClientBuilder] = (*ClientHookBuilder)(nil)

// ClientHookBuilder allows for setting pre-request and post-response hooks at the client level.
type ClientHookBuilder struct {
	parentBuilder *ClientBuilder
}

// Hook returns a new ClientHookBuilder for setting request hooks.
func (b *ClientBuilder) Hook() *ClientHookBuilder {
	return &ClientHookBuilder{parentBuilder: b}
}

// OnBeforeRequest adds a pre-request hook to the client.
func (b *ClientHookBuilder) OnBeforeRequest(hook func(*http.Request) error) *ClientBuilder {
	b.parentBuilder.client.AddBeforeRequestHook(hook)
	return b.parentBuilder
}

// OnAfterResponse adds a post-response hook to the client.
func (b *ClientHookBuilder) OnAfterResponse(hook func(*http.Request, *http.Response)) *ClientBuilder {
	b.parentBuilder.client.AddAfterResponseHook(hook)
	return b.parentBuilder
}
