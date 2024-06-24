package fastshot

import (
	"net/url"
	"sync/atomic"
)

type (
	// DefaultBaseURL implements ConfigBaseURL interface and provides a single base URL.
	DefaultBaseURL struct {
		baseURL *url.URL
	}

	// BalancedBaseURL implements ConfigBaseURL interface and provides load balancing.
	BalancedBaseURL struct {
		baseURLs       []*url.URL
		currentBaseURL uint32
	}
)

// BaseURL for DefaultBaseURL returns the base URL.
func (c *DefaultBaseURL) BaseURL() *url.URL {
	return c.baseURL
}

// BaseURL for BalancedBaseURL returns the next base URL in the list.
func (c *BalancedBaseURL) BaseURL() *url.URL {
	currentIndex := atomic.LoadUint32(&c.currentBaseURL)
	atomic.AddUint32(&c.currentBaseURL, 1)
	c.currentBaseURL = c.currentBaseURL % uint32(len(c.baseURLs))
	return c.baseURLs[currentIndex]
}
