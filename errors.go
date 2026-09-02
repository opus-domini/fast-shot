package fastshot

import (
	"errors"
	"fmt"
)

// Sentinel errors returned (wrapped) by Send and the builders.
// Use errors.Is or errors.AsType to inspect them.
var (
	// ErrClientValidation reports invalid client attributes, such as an empty base URL.
	ErrClientValidation = errors.New("invalid client attributes")
	// ErrRequestValidation reports invalid request attributes accumulated while building.
	ErrRequestValidation = errors.New("invalid request attributes")
	// ErrCreateRequest reports a failure to build the underlying *http.Request.
	ErrCreateRequest = errors.New("failed to create request")
	// ErrBeforeRequestHook reports a before-request hook aborting the request.
	ErrBeforeRequestHook = errors.New("before request hook failed")
	// ErrSetBody reports a failure to set the request body.
	ErrSetBody = errors.New("failed to set body")
	// ErrMarshalJSON reports a failure to marshal the body as JSON.
	ErrMarshalJSON = errors.New("failed to marshal JSON")
	// ErrMarshalXML reports a failure to marshal the body as XML.
	ErrMarshalXML = errors.New("failed to marshal XML")
	// ErrParseURL reports a failure to parse a URL.
	ErrParseURL = errors.New("failed to parse URL")
	// ErrParseProxyURL reports a failure to parse a proxy URL.
	ErrParseProxyURL = errors.New("failed to parse proxy URL")
	// ErrParseQueryString reports a failure to parse a raw query string.
	ErrParseQueryString = errors.New("failed to parse query string")
	// ErrEmptyBaseURL reports an empty base URL.
	ErrEmptyBaseURL = errors.New("empty base URL")
)

// RetryError reports that a request exhausted all retry attempts.
//
// Inspect it with errors.AsType (Go 1.26+) or errors.As:
//
//	if retryErr, err := errors.AsType[*fastshot.RetryError](err); err == nil {
//		log.Printf("gave up after %d attempts", retryErr.Attempts)
//	}
type RetryError struct {
	// Attempts is the number of attempts made.
	Attempts uint
	// Err joins one error per failed attempt.
	Err error
}

// Error implements the error interface.
func (e *RetryError) Error() string {
	return fmt.Sprintf("request failed after %d attempts: %v", e.Attempts, e.Err)
}

// Unwrap exposes the joined per-attempt errors.
func (e *RetryError) Unwrap() error {
	return e.Err
}

func newRetryError(attempts uint, errs []error) *RetryError {
	return &RetryError{
		Attempts: attempts,
		Err:      errors.Join(errs...),
	}
}
