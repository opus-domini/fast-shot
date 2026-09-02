package fastshot

import (
	"encoding/base64"
	"net/http"

	"github.com/opus-domini/fast-shot/constant/header"
)

// AuthBuilder provides fluent HTTP authentication configuration.
//
// The generic type parameter P is the parent builder returned by every
// method to keep the fluent chain going, so the same builder serves both
// ClientBuilder and RequestBuilder:
//
//	fastshot.NewClient("https://api.example.com").
//		Auth().BearerToken("my-access-token").
//		Build()
//
//	client.GET("/protected").
//		Auth().BasicAuth("username", "password").
//		Send()
type AuthBuilder[P any] struct {
	parent P
	header http.Header
}

// Auth returns an AuthBuilder for setting authentication on the request.
func (b *RequestBuilder) Auth() *AuthBuilder[*RequestBuilder] {
	return &AuthBuilder[*RequestBuilder]{
		parent: b,
		header: b.request.config.httpHeader,
	}
}

// Auth returns an AuthBuilder for setting authentication on the client.
func (b *ClientBuilder) Auth() *AuthBuilder[*ClientBuilder] {
	return &AuthBuilder[*ClientBuilder]{
		parent: b,
		header: b.client.Header(),
	}
}

// Set sets the Authorization header for custom authentication.
func (b *AuthBuilder[P]) Set(value string) P {
	b.header.Set(header.Authorization.String(), value)
	return b.parent
}

// BearerToken sets the Authorization header for Bearer token authentication.
func (b *AuthBuilder[P]) BearerToken(token string) P {
	return b.Set("Bearer " + token)
}

// BasicAuth sets the Authorization header for Basic authentication.
func (b *AuthBuilder[P]) BasicAuth(username, password string) P {
	encoded := base64.StdEncoding.EncodeToString([]byte(username + ":" + password))
	return b.Set("Basic " + encoded)
}
