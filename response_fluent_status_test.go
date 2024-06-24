package fastshot

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
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
				status: &ResponseFluentStatus{
					response: &http.Response{StatusCode: tt.statusCode},
				},
			}

			// Act
			result := response.Status()

			// Assert
			assert.Equal(t, tt.expectedCode, result.Code())
			assert.Equal(t, tt.expectedText, result.Text())
			assert.Equal(t, tt.expectedInformation, result.Is1xxInformational())
			assert.Equal(t, tt.expectedSuccess, result.Is2xxSuccessful())
			assert.Equal(t, tt.expectedRedirection, result.Is3xxRedirection())
			assert.Equal(t, tt.expectedClientError, result.Is4xxClientError())
			assert.Equal(t, tt.expectedServerError, result.Is5xxServerError())
			assert.Equal(t, tt.expectedOK, result.IsOK())
			assert.Equal(t, tt.expectedNotFound, result.IsNotFound())
			assert.Equal(t, tt.expectedUnauthorized, result.IsUnauthorized())
			assert.Equal(t, tt.expectedForbidden, result.IsForbidden())
			assert.Equal(t, tt.expectedError, result.IsError())
		})
	}
}
