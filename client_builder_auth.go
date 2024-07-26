package fastshot

import (
	"encoding/base64"

	"github.com/opus-domini/fast-shot/constant/header"
)

// BuilderAuth is the interface that wraps the basic methods for setting authentication configurations.
var _ BuilderAuth[ClientBuilder] = (*ClientAuthBuilder)(nil)

// ClientAuthBuilder allows for setting authentication configurations.
type ClientAuthBuilder struct {
	parentBuilder *ClientBuilder
}

// Auth BuilderAuth returns a new ClientAuthBuilder for setting authentication options.
func (b *ClientBuilder) Auth() *ClientAuthBuilder {
	return &ClientAuthBuilder{parentBuilder: b}
}

// Set sets the Authorization header for custom authentication.
func (b *ClientAuthBuilder) Set(value string) *ClientBuilder {
	b.parentBuilder.client.Header().Set(header.Authorization, value)
	return b.parentBuilder
}

// BearerToken sets the Authorization header for Bearer token authentication.
func (b *ClientAuthBuilder) BearerToken(token string) *ClientBuilder {
	b.Set("Bearer " + token)
	return b.parentBuilder
}

// BasicAuth sets the Authorization header for Basic authentication.
func (b *ClientAuthBuilder) BasicAuth(username, password string) *ClientBuilder {
	encoded := base64.StdEncoding.EncodeToString([]byte(username + ":" + password))
	b.Set("Basic " + encoded)
	return b.parentBuilder
}
