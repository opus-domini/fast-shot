package fastshot

import (
	"io"
	"net/http"
	"strings"
	"testing"
)

func TestResponseFluentStatus(t *testing.T) {
	tests := []struct {
		name                 string
		statusCode           int
		expectedCode         int
		expectedText         string
		expectedInformation  bool
		expectedSuccess      bool
		expectedRedirection  bool
		expectedClientError  bool
		expectedServerError  bool
		expectedOK           bool
		expectedNotFound     bool
		expectedUnauthorized bool
		expectedForbidden    bool
		expectedError        bool
	}{
		{
			name:                 "200 OK",
			statusCode:           http.StatusOK,
			expectedCode:         200,
			expectedText:         "[200] OK",
			expectedInformation:  false,
			expectedSuccess:      true,
			expectedRedirection:  false,
			expectedClientError:  false,
			expectedServerError:  false,
			expectedOK:           true,
			expectedNotFound:     false,
			expectedUnauthorized: false,
			expectedForbidden:    false,
			expectedError:        false,
		},
		{
			name:                 "404 Not Found",
			statusCode:           http.StatusNotFound,
			expectedCode:         404,
			expectedText:         "[404] Not Found",
			expectedInformation:  false,
			expectedSuccess:      false,
			expectedRedirection:  false,
			expectedClientError:  true,
			expectedServerError:  false,
			expectedOK:           false,
			expectedNotFound:     true,
			expectedUnauthorized: false,
			expectedForbidden:    false,
			expectedError:        true,
		},
		{
			name:                 "500 Internal Server Error",
			statusCode:           http.StatusInternalServerError,
			expectedCode:         500,
			expectedText:         "[500] Internal Server Error",
			expectedInformation:  false,
			expectedSuccess:      false,
			expectedRedirection:  false,
			expectedClientError:  false,
			expectedServerError:  true,
			expectedOK:           false,
			expectedNotFound:     false,
			expectedUnauthorized: false,
			expectedForbidden:    false,
			expectedError:        true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			response := &Response{
				rawResponse: &http.Response{
					StatusCode: tt.statusCode,
					Body:       io.NopCloser(strings.NewReader("")),
				},
				status: &ResponseFluentStatus{
					response: &http.Response{
						StatusCode: tt.statusCode,
						Body:       io.NopCloser(strings.NewReader("")),
					},
				},
			}

			// Act
			result := response.Status()

			// Assert
			if got := result.Code(); got != tt.expectedCode {
				t.Errorf("Code() got %d, want %d", got, tt.expectedCode)
			}
			if got := result.Text(); got != tt.expectedText {
				t.Errorf("Text() got %q, want %q", got, tt.expectedText)
			}
			rawResp := response.Raw()
			if rawResp.StatusCode != tt.expectedCode {
				t.Errorf("Raw().StatusCode got %d, want %d", rawResp.StatusCode, tt.expectedCode)
			}
			_ = rawResp.Body.Close()
			_ = response.status.response.Body.Close()
			if got := result.Is1xxInformational(); got != tt.expectedInformation {
				t.Errorf("Is1xxInformational() got %v, want %v", got, tt.expectedInformation)
			}
			if got := result.Is2xxSuccessful(); got != tt.expectedSuccess {
				t.Errorf("Is2xxSuccessful() got %v, want %v", got, tt.expectedSuccess)
			}
			if got := result.Is3xxRedirection(); got != tt.expectedRedirection {
				t.Errorf("Is3xxRedirection() got %v, want %v", got, tt.expectedRedirection)
			}
			if got := result.Is4xxClientError(); got != tt.expectedClientError {
				t.Errorf("Is4xxClientError() got %v, want %v", got, tt.expectedClientError)
			}
			if got := result.Is5xxServerError(); got != tt.expectedServerError {
				t.Errorf("Is5xxServerError() got %v, want %v", got, tt.expectedServerError)
			}
			if got := result.IsOK(); got != tt.expectedOK {
				t.Errorf("IsOK() got %v, want %v", got, tt.expectedOK)
			}
			if got := result.IsNotFound(); got != tt.expectedNotFound {
				t.Errorf("IsNotFound() got %v, want %v", got, tt.expectedNotFound)
			}
			if got := result.IsUnauthorized(); got != tt.expectedUnauthorized {
				t.Errorf("IsUnauthorized() got %v, want %v", got, tt.expectedUnauthorized)
			}
			if got := result.IsForbidden(); got != tt.expectedForbidden {
				t.Errorf("IsForbidden() got %v, want %v", got, tt.expectedForbidden)
			}
			if got := result.IsError(); got != tt.expectedError {
				t.Errorf("IsError() got %v, want %v", got, tt.expectedError)
			}
		})
	}
}
