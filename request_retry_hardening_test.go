package fastshot

import (
	"context"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"
	"time"
)

func TestRequestBuilder_requestForAttempt(t *testing.T) {
	rb := &RequestBuilder{
		request: &Request{
			config: newRequestConfigBase("", "", DefaultJSONCodec()),
		},
	}

	t.Run("first attempt returns the original request", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "http://example.com", nil)
		got, err := rb.requestForAttempt(req, 0)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if got != req {
			t.Error("got different request, want same")
		}
	})

	t.Run("attempt without replayable body returns the original request", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "http://example.com", nil)
		req.GetBody = nil
		got, err := rb.requestForAttempt(req, 1)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if got != req {
			t.Error("got different request, want same")
		}
	})

	t.Run("body replay error is returned", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "http://example.com", nil)
		req.GetBody = func() (io.ReadCloser, error) { return nil, errors.New("replay boom") }

		_, err := rb.requestForAttempt(req, 1)
		if err == nil || err.Error() != "replay boom" {
			t.Errorf("error got %v, want %q", err, "replay boom")
		}
	})

	t.Run("body replay error aborts the retry loop", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
		}))
		defer server.Close()

		rb := &RequestBuilder{
			request: &Request{
				client: newClientConfigBase(server.URL),
				config: newRequestConfigBase("", "", DefaultJSONCodec()),
			},
		}
		rc := rb.request.config.RetryConfig()
		rc.SetInterval(time.Millisecond)
		rc.SetMaxAttempts(3)
		rc.SetBackoffRate(1)
		rc.SetJitterStrategy(JitterStrategyNone)

		req, reqErr := http.NewRequest(http.MethodGet, server.URL, nil)
		if reqErr != nil {
			t.Fatalf("unexpected error creating request: %v", reqErr)
		}
		replayErr := errors.New("replay boom")
		req.GetBody = func() (io.ReadCloser, error) { return nil, replayErr }

		_, err := rb.executeWithRetry(req)
		retryErrResult, ok := errors.AsType[*RetryError](err)
		if !ok {
			t.Fatalf("error is not *RetryError: %v", err)
		}
		if retryErrResult.Attempts != 3 {
			t.Errorf("Attempts got %d, want 3", retryErrResult.Attempts)
		}
		if !errors.Is(err, replayErr) {
			t.Errorf("error should wrap the replay failure, got %v", err)
		}
	})
}

func TestRequestBuilder_sleepBackoff(t *testing.T) {
	rb := &RequestBuilder{
		request: &Request{
			config: newRequestConfigBase("", "", DefaultJSONCodec()),
		},
	}

	t.Run("non-positive delay returns immediately", func(t *testing.T) {
		if err := rb.sleepBackoff(context.Background(), 0); err != nil {
			t.Errorf("unexpected error: %v", err)
		}
	})

	t.Run("canceled context aborts the backoff", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		cancel()

		err := rb.sleepBackoff(ctx, time.Hour)
		if !errors.Is(err, context.Canceled) {
			t.Errorf("error got %v, want context.Canceled", err)
		}
	})

	t.Run("elapsed delay returns nil", func(t *testing.T) {
		if err := rb.sleepBackoff(context.Background(), time.Millisecond); err != nil {
			t.Errorf("unexpected error: %v", err)
		}
	})
}

func TestRequest_Retry_AbortsOnCanceledContext(t *testing.T) {
	// Arrange
	var attempts atomic.Int32
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		attempts.Add(1)
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer server.Close()

	client := NewClient(server.URL).Build()

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	// Act
	_, err := client.GET("/flaky").
		Context().Set(ctx).
		Retry().SetConstantBackoff(time.Millisecond, 5).
		Send()

	// Assert
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !errors.Is(err, context.Canceled) {
		t.Errorf("error should wrap context.Canceled, got %v", err)
	}
	if _, ok := errors.AsType[*RetryError](err); !ok {
		t.Errorf("error is not *RetryError: %v", err)
	}
	if got := attempts.Load(); got != 0 {
		t.Errorf("server attempts got %d, want 0", got)
	}
}
