package fastshot

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"

	"github.com/opus-domini/fast-shot/constant"
)

// ClientBuilder serves as the main entry point for configuring HTTP clients.
type ClientBuilder struct {
	client Client
}

// NewClient initializes a new ClientBuilder with a given baseURL.
func NewClient(baseURL string) *ClientBuilder {
	var validations []error

	if baseURL == "" {
		validations = append(validations, errors.New(constant.ErrMsgEmptyBaseURL))
	}

	parsedURL, err := url.Parse(baseURL)
	if err != nil {
		validations = append(validations, errors.Join(errors.New(constant.ErrMsgParseURL), err))
	}

	return &ClientBuilder{
		client: &ClientConfigBase{
			httpClient: &DefaultHttpClient{
				client: &http.Client{},
			},
			httpHeader:  &http.Header{},
			httpCookies: []*http.Cookie{},
			validations: validations,
			ConfigBaseURL: &DefaultBaseURL{
				baseURL: parsedURL,
			},
		},
	}
}

// NewClientLoadBalancer initializes a new ClientBuilder with a given baseURLs.
func NewClientLoadBalancer(baseURLs []string) *ClientBuilder {
	var validations []error

	var parsedURLs []*url.URL
	for index, baseURL := range baseURLs {
		if baseURL == "" {
			validations = append(validations, fmt.Errorf("base URL %d: %s", index, constant.ErrMsgEmptyBaseURL))
			continue
		}

		parsedURL, err := url.Parse(baseURL)
		if err != nil {
			validations = append(validations, errors.Join(errors.New(constant.ErrMsgParseURL), err))
		}
		parsedURLs = append(parsedURLs, parsedURL)
	}

	if len(parsedURLs) == 0 {
		validations = append(validations, errors.New(constant.ErrMsgEmptyBaseURL))
	}

	return &ClientBuilder{
		client: &ClientConfigBase{
			httpClient: &DefaultHttpClient{
				client: &http.Client{},
			},
			httpHeader:  &http.Header{},
			httpCookies: []*http.Cookie{},
			validations: validations,
			ConfigBaseURL: &BalancedBaseURL{
				baseURLs:       parsedURLs,
				currentBaseURL: 0,
			},
		},
	}
}

// DefaultClient initializes a new default ClientConfig with a given baseURL.
func DefaultClient(baseURL string) ClientHttpMethods {
	return NewClient(baseURL).Build()
}

// DefaultClientLoadBalancer initializes a new default ClientConfig with a given baseURLs.
func DefaultClientLoadBalancer(baseURLs []string) ClientHttpMethods {
	return NewClientLoadBalancer(baseURLs).Build()
}

// Build finalizes the ClientBuilder configurations and returns a new ClientConfig.
func (b *ClientBuilder) Build() ClientHttpMethods {
	return b.client
}
