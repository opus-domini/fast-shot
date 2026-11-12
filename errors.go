package fastshot

import (
	"errors"
	"fmt"

	"github.com/opus-domini/fast-shot/constant"
)

// Sentinel errors returned (wrapped) by Send and the builders.
// Use errors.Is or errors.AsType to inspect them.
var (
	// ErrClientValidation reports invalid client attributes, such as an empty base URL.
	ErrClientValidation = errors.New(constant.ErrMsgClientValidation)
	// ErrRequestValidation reports invalid request attributes accumulated while building.
	ErrRequestValidation = errors.New(constant.ErrMsgRequestValidation)
	// ErrCreateRequest reports a failure to build the underlying *http.Request.
	ErrCreateRequest = errors.New(constant.ErrMsgCreateRequest)
	// ErrBeforeRequestHook reports a before-request hook aborting the request.
	ErrBeforeRequestHook = errors.New(constant.ErrMsgBeforeRequestHook)
	// ErrSetBody reports a failure to set the request body.
	ErrSetBody = errors.New(constant.ErrMsgSetBody)
	// ErrMarshalJSON reports a failure to marshal the body as JSON.
	ErrMarshalJSON = errors.New(constant.ErrMsgMarshalJSON)
	// ErrMarshalXML reports a failure to marshal the body as XML.
	ErrMarshalXML = errors.New(constant.ErrMsgMarshalXML)
	// ErrParseURL reports a failure to parse a URL.
	ErrParseURL = errors.New(constant.ErrMsgParseURL)
	// ErrParseProxyURL reports a failure to parse a proxy URL.
	ErrParseProxyURL = errors.New(constant.ErrMsgParseProxyURL)
	// ErrParseQueryString reports a failure to parse a raw query string.
	ErrParseQueryString = errors.New(constant.ErrMsgParseQueryString)
	// ErrEmptyBaseURL reports an empty base URL.
	ErrEmptyBaseURL = errors.New(constant.ErrMsgEmptyBaseURL)
)

// RetryError reports that a request exhausted all retry attempts.
//
// Inspect it with errors.AsType (Go 1.26+) or errors.As:
//
//	if retryErr, ok := errors.AsType[*fastshot.RetryError](err); ok {
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
