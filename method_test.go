package fastshot

import (
	"github.com/opus-domini/fast-shot/constant/method"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestClientMethods(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, err := w.Write([]byte("OK"))
		if err != nil {
			return
		}
	}))
	defer server.Close()

	client := DefaultClient(server.URL)

	tests := []struct {
		name       string
		methodFunc func(string) *RequestBuilder
	}{
		{method.CONNECT, client.CONNECT},
		{method.DELETE, client.DELETE},
		{method.GET, client.GET},
		{method.HEAD, client.HEAD},
		{method.OPTIONS, client.OPTIONS},
		{method.PATCH, client.PATCH},
		{method.POST, client.POST},
		{method.PUT, client.PUT},
		{method.TRACE, client.TRACE},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := tt.methodFunc("/")
			resp, _ := req.Send()
			if resp.IsError() {
				t.Errorf("Expected 200, got %d", resp.StatusCode())
			}
		})
	}
}
