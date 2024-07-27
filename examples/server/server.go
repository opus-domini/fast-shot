package server

import (
	"log/slog"
	"net/http/httptest"

	"github.com/opus-domini/fast-shot/examples/server/handler"
)

// Start a new test server using the handler defined in the examples package.
// This server is intended for testing purposes only.
func Start() *httptest.Server {
	ts := httptest.NewServer(handler.NewMux())
	slog.Info("Test server for examples start!", "url", ts.URL)
	return ts
}
