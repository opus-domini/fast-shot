package fastshot

import (
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestSentinelErrors(t *testing.T) {
	sentinels := []error{
		ErrClientValidation,
		ErrRequestValidation,
		ErrCreateRequest,
		ErrBeforeRequestHook,
		ErrSetBody,
		ErrMarshalJSON,
		ErrMarshalXML,
		ErrParseURL,
		ErrParseProxyURL,
		ErrParseQueryString,
		ErrEmptyBaseURL,
	}

	for _, sentinel := range sentinels {
		wrapped := fmt.Errorf("context: %w", sentinel)
		if !errors.Is(wrapped, sentinel) {
			t.Errorf("errors.Is failed for %q", sentinel)
		}
	}
}

func TestRetryError(t *testing.T) {
	attemptErr := errors.New("attempt 1: boom")
	retryErr := newRetryError(3, []error{attemptErr})

	t.Run("message includes attempts", func(t *testing.T) {
		want := "request failed after 3 attempts: attempt 1: boom"
		if got := retryErr.Error(); got != want {
			t.Errorf("Error() got %q, want %q", got, want)
		}
	})

	t.Run("unwrap exposes joined errors", func(t *testing.T) {
		if !errors.Is(retryErr, attemptErr) {
			t.Error("errors.Is(retryErr, attemptErr) got false, want true")
		}
	})

	t.Run("inspect with errors.AsType", func(t *testing.T) {
		wrapped := fmt.Errorf("send failed: %w", retryErr)
		inspected, ok := errors.AsType[*RetryError](wrapped)
		if !ok {
			t.Fatalf("errors.AsType did not match *RetryError for %v", wrapped)
		}
		if inspected.Attempts != 3 {
			t.Errorf("Attempts got %d, want 3", inspected.Attempts)
		}
	})

	t.Run("returned by Send after exhausting attempts", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
		}))
		defer server.Close()

		client := NewClient(server.URL).Build()

		_, err := client.GET("/flaky").
			Retry().SetConstantBackoff(time.Millisecond, 2).
			Send()
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		if _, ok := errors.AsType[*RetryError](err); !ok {
			t.Errorf("error is not *RetryError: %v", err)
		}
	})
}
