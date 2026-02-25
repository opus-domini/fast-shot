package fastshot

import (
	"net/url"
	"sync"
	"testing"
)

func TestDefaultBaseURL(t *testing.T) {
	// Arrange
	u, _ := url.Parse("https://example.com")
	base := newDefaultBaseURL(u)

	// Act & Assert
	if got := base.BaseURL(); got != u {
		t.Errorf("got %v, want %v", got, u)
	}
	if got := base.BaseURL(); got != u {
		t.Errorf("got %v, want %v", got, u)
	}
}

func TestBalancedBaseURL_RoundRobin(t *testing.T) {
	// Arrange
	u1, _ := url.Parse("https://a.com")
	u2, _ := url.Parse("https://b.com")
	u3, _ := url.Parse("https://c.com")
	base := newBalancedBaseURL([]*url.URL{u1, u2, u3})

	// Act & Assert
	if got := base.BaseURL(); got != u1 {
		t.Errorf("got %v, want %v", got, u1)
	}
	if got := base.BaseURL(); got != u2 {
		t.Errorf("got %v, want %v", got, u2)
	}
	if got := base.BaseURL(); got != u3 {
		t.Errorf("got %v, want %v", got, u3)
	}
	if got := base.BaseURL(); got != u1 { // wraps around
		t.Errorf("got %v, want %v", got, u1)
	}
	if got := base.BaseURL(); got != u2 {
		t.Errorf("got %v, want %v", got, u2)
	}
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
	for range 100 {
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
		found := false
		for _, valid := range urls {
			if u == valid {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("got %v, want one of %v", u, urls)
		}
	}
}
