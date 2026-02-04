package fastshot

import (
	"net/url"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDefaultBaseURL(t *testing.T) {
	// Arrange
	u, _ := url.Parse("https://example.com")
	base := newDefaultBaseURL(u)

	// Act & Assert
	assert.Equal(t, u, base.BaseURL())
	assert.Equal(t, u, base.BaseURL())
}

func TestBalancedBaseURL_RoundRobin(t *testing.T) {
	// Arrange
	u1, _ := url.Parse("https://a.com")
	u2, _ := url.Parse("https://b.com")
	u3, _ := url.Parse("https://c.com")
	base := newBalancedBaseURL([]*url.URL{u1, u2, u3})

	// Act & Assert
	assert.Equal(t, u1, base.BaseURL())
	assert.Equal(t, u2, base.BaseURL())
	assert.Equal(t, u3, base.BaseURL())
	assert.Equal(t, u1, base.BaseURL()) // wraps around
	assert.Equal(t, u2, base.BaseURL())
}

func TestBalancedBaseURL_Concurrent(t *testing.T) {
	// Arrange
	u1, _ := url.Parse("https://a.com")
	u2, _ := url.Parse("https://b.com")
	u3, _ := url.Parse("https://c.com")
	urls := []*url.URL{u1, u2, u3}
	base := newBalancedBaseURL(urls)

	// Act
	var wg sync.WaitGroup
	results := make(chan *url.URL, 100)
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			results <- base.BaseURL()
		}()
	}
	wg.Wait()
	close(results)

	// Assert - all returned URLs must be one of the valid base URLs
	for u := range results {
		assert.Contains(t, urls, u)
	}
}
