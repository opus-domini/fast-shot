package fastshot

import (
	"encoding/base64"
	"github.com/opus-domini/fast-shot/constant/header"
)

// Auth is the interface that wraps the basic methods for setting authentication configurations.
var _ Auth[RequestBuilder] = (*RequestAuthBuilder)(nil)

// RequestAuthBuilder allows for setting authentication configurations.
type RequestAuthBuilder struct {
	parentBuilder *RequestBuilder
}

// Auth returns a new ClientAuthBuilder for setting authentication options.
func (r *RequestBuilder) Auth() *RequestAuthBuilder {
	return &RequestAuthBuilder{parentBuilder: r}
}

// Set sets the Authorization header for custom authentication.
func (r *RequestAuthBuilder) Set(value string) *RequestBuilder {
	r.parentBuilder.request.httpHeader.Set(header.Authorization, value)
	return r.parentBuilder
}

// BearerToken sets the Authorization header for Bearer token authentication.
func (r *RequestAuthBuilder) BearerToken(token string) *RequestBuilder {
	r.Set("Bearer " + token)
	return r.parentBuilder
}

// BasicAuth sets the Authorization header for Basic authentication.
func (r *RequestAuthBuilder) BasicAuth(username, password string) *RequestBuilder {
	encoded := base64.StdEncoding.EncodeToString([]byte(username + ":" + password))
	r.Set("Basic " + encoded)
	return r.parentBuilder
}
