package fastshot

import (
	"net/http"
)

// HookBuilder provides fluent configuration of pre-request and post-response hooks.
//
// Pre-request hooks can inspect or modify the *http.Request and optionally abort
// by returning an error. Post-response hooks are observational and receive both
// the request and the response.
//
// The generic type parameter P is the parent builder returned by every
// method to keep the fluent chain going, so the same builder serves both
// ClientBuilder and RequestBuilder:
//
//	fastshot.NewClient("https://api.example.com").
//		Hook().OnBeforeRequest(func(req *http.Request) error {
//			req.Header.Set("X-Request-ID", uuid.New().String()) // stdlib uuid, Go 1.27+
//			return nil
//		}).
//		Build()
type HookBuilder[P any] struct {
	parent    P
	addBefore func(func(*http.Request) error)
	addAfter  func(func(*http.Request, *http.Response))
}

// Hook returns a HookBuilder for setting request hooks on the request.
func (b *RequestBuilder) Hook() *HookBuilder[*RequestBuilder] {
	return &HookBuilder[*RequestBuilder]{
		parent:    b,
		addBefore: b.request.config.AddBeforeRequestHook,
		addAfter:  b.request.config.AddAfterResponseHook,
	}
}

// Hook returns a HookBuilder for setting request hooks on the client.
func (b *ClientBuilder) Hook() *HookBuilder[*ClientBuilder] {
	return &HookBuilder[*ClientBuilder]{
		parent:    b,
		addBefore: b.client.AddBeforeRequestHook,
		addAfter:  b.client.AddAfterResponseHook,
	}
}

// OnBeforeRequest adds a pre-request hook.
func (b *HookBuilder[P]) OnBeforeRequest(hook func(*http.Request) error) P {
	b.addBefore(hook)
	return b.parent
}

// OnAfterResponse adds a post-response hook.
func (b *HookBuilder[P]) OnAfterResponse(hook func(*http.Request, *http.Response)) P {
	b.addAfter(hook)
	return b.parent
}
