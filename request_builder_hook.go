package fastshot

import (
	"net/http"
)

// BuilderHook is the interface that wraps the basic methods for setting request hooks.
var _ BuilderHook[RequestBuilder] = (*RequestHookBuilder)(nil)

// RequestHookBuilder allows for setting pre-request and post-response hooks at the request level.
type RequestHookBuilder struct {
	parentBuilder *RequestBuilder
	requestConfig *RequestConfigBase
}

// Hook returns a new RequestHookBuilder for setting request hooks.
func (b *RequestBuilder) Hook() *RequestHookBuilder {
	return &RequestHookBuilder{
		parentBuilder: b,
		requestConfig: b.request.config,
	}
}

// OnBeforeRequest adds a pre-request hook to the request.
func (b *RequestHookBuilder) OnBeforeRequest(hook func(*http.Request) error) *RequestBuilder {
	b.requestConfig.AddBeforeRequestHook(hook)
	return b.parentBuilder
}

// OnAfterResponse adds a post-response hook to the request.
func (b *RequestHookBuilder) OnAfterResponse(hook func(*http.Request, *http.Response)) *RequestBuilder {
	b.requestConfig.AddAfterResponseHook(hook)
	return b.parentBuilder
}
