package fastshot

import (
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
		methodFunc func(string) *Request
	}{
		{http.MethodGet, client.GET},
		{http.MethodPost, client.POST},
		{http.MethodPut, client.PUT},
		{http.MethodDelete, client.DELETE},
		{http.MethodPatch, client.PATCH},
		{http.MethodHead, client.HEAD},
		{http.MethodOptions, client.OPTIONS},
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
