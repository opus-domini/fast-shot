package fastshot

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func makeResponseFromServer(statusCode int) *Response {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(statusCode)
	}))
	defer server.Close()

	res, _ := http.Get(server.URL)
	return &Response{RawResponse: res}
}

func makeResponse(statusCode int) *Response {
	return &Response{RawResponse: &http.Response{StatusCode: statusCode}}
}

func TestResponse_Is1xxInformational(t *testing.T) {
	// Test case when RawResponse is nil
	r := &Response{}
	if r.Is1xxInformational() {
		t.Errorf("Is1xxInformational should return false when RawResponse is nil")
	}

	// Test case when status is 100
	r = makeResponse(100)
	if !r.Is1xxInformational() {
		t.Errorf("Is1xxInformational failed for status 100")
	}
}

func TestResponse_Is2xxSuccessful(t *testing.T) {
	// Test case when RawResponse is nil
	r := &Response{}
	if r.Is2xxSuccessful() {
		t.Errorf("Is2xxSuccessful should return false when RawResponse is nil")
	}

	// Test case when status is 200
	r = makeResponse(200)
	if !r.Is2xxSuccessful() {
		t.Errorf("Is2xxSuccessful failed for status 200")
	}
}

func TestResponse_Is3xxRedirection(t *testing.T) {
	// Test case when RawResponse is nil
	r := &Response{}
	if r.Is3xxRedirection() {
		t.Errorf("Is3xxRedirection should return false when RawResponse is nil")
	}

	// Test case when status is 300
	r = makeResponse(300)
	if !r.Is3xxRedirection() {
		t.Errorf("Is3xxRedirection failed for status 300")
	}
}

func TestResponse_Is4xxClientError(t *testing.T) {
	// Test case when RawResponse is nil
	r := &Response{}
	if r.Is4xxClientError() {
		t.Errorf("Is4xxClientError should return false when RawResponse is nil")
	}

	// Test case when status is 400
	r = makeResponse(400)
	if !r.Is4xxClientError() {
		t.Errorf("Is4xxClientError failed for status 400")
	}
}

func TestResponse_Is5xxServerError(t *testing.T) {
	// Test case when RawResponse is nil
	r := &Response{}
	if r.Is5xxServerError() {
		t.Errorf("Is5xxServerError should return false when RawResponse is nil")
	}

	// Test case when status is 500
	r = makeResponse(500)
	if !r.Is5xxServerError() {
		t.Errorf("Is5xxServerError failed for status 500")
	}
}

func TestResponse_IsError(t *testing.T) {
	// Test case when RawResponse is nil
	r := &Response{}
	if r.IsError() {
		t.Errorf("IsError should return false when RawResponse is nil")
	}

	// Test case when status is 400
	r = makeResponse(400)
	if !r.IsError() {
		t.Errorf("IsError failed for status 400")
	}

	// Test case when status is 500
	r = makeResponse(500)
	if !r.IsError() {
		t.Errorf("IsError failed for status 500")
	}
}

func TestResponse_StatusCode(t *testing.T) {
	// Test case when RawResponse is nil
	r := &Response{}
	if r.StatusCode() != 0 {
		t.Errorf("StatusCode should return 0 when RawResponse is nil")
	}

	// Test case when status is 200
	r = makeResponse(200)
	if r.StatusCode() != 200 {
		t.Errorf("StatusCode failed for status 200")
	}
}

func TestResponse_Status(t *testing.T) {
	// Test case when RawResponse is nil
	r := &Response{}
	if r.Status() != "" {
		t.Errorf("Status should return empty string when RawResponse is nil")
	}

	// Test case when status is 200
	r = makeResponseFromServer(200)
	if r.Status() != "200 OK" {
		t.Errorf("Status failed for status 200")
	}
}

func TestResponse_RawBody(t *testing.T) {
	// Test case when RawResponse is nil
	r := &Response{}
	if r.RawBody() != nil {
		t.Errorf("RawBody should return nil when RawResponse is nil")
	}

	// Test case when status is 200
	r = makeResponseFromServer(200)
	if r.RawBody() == nil {
		t.Errorf("RawBody is nil for status 200")
	} else {
		_, err := io.ReadAll(r.RawBody())
		if err != nil {
			return
		}
		_ = r.RawBody().Close()
	}
}
