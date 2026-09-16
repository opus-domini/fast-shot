package fastshot

import (
	"net/http"
	"net/http/httptest"
	"runtime"
	"runtime/pprof"
	"strings"
	"testing"
	"time"
)

// TestNoGoroutineLeaks exercises the client (including the retry path, which
// discards intermediate responses) and asserts that the Go 1.27 goroutineleak
// profile stays empty afterwards.
//
// The profile reports goroutines blocked on unreachable concurrency primitives
// (channels, mutexes), so it guards against channel/mutex leaks in hooks,
// retries, and body handling. Connection hygiene for discarded responses is
// enforced by closing their bodies in executeWithRetry.
func TestNoGoroutineLeaks(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))

	client := DefaultClient(server.URL)

	// Retry path: discarded 5xx responses must have their bodies closed so
	// connections are released instead of leaking.
	_, _ = client.GET("/flaky").
		Retry().SetConstantBackoff(time.Millisecond, 3).
		Send()

	// Happy path with body consumed.
	resp, err := client.GET("/ok").Send()
	if err == nil {
		_, _ = resp.Body().AsString()
	}

	// Close the server explicitly so idle keep-alive connections are torn
	// down and their goroutines exit before the leak profile is taken.
	server.Close()

	// The leak profile is computed from reachability, so give the GC a chance
	// to collect the discarded responses and their connections.
	for range 5 {
		runtime.GC()
		time.Sleep(10 * time.Millisecond)
	}

	profile := pprof.Lookup("goroutineleak")
	if profile == nil {
		t.Skip("goroutineleak profile not available in this build")
	}

	var buf strings.Builder
	if err := profile.WriteTo(&buf, 1); err != nil {
		t.Fatalf("unexpected error writing goroutineleak profile: %v", err)
	}
	// The text profile always writes a header such as
	// "goroutineleak profile: total 0"; count the reported goroutines.
	if count := profile.Count(); count > 0 {
		t.Errorf("goroutine leaks detected (%d):\n%s", count, buf.String())
	}
}
