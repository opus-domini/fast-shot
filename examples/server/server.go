package server

import (
	"log/slog"
	"net/http/httptest"

	"github.com/opus-domini/fast-shot/examples/server/database"
	"github.com/opus-domini/fast-shot/examples/server/handler"
	"github.com/opus-domini/fast-shot/examples/server/repository"
)

// Manager manages the shared database and the creation of test servers.
type Manager struct {
	serverCount int
	repository  *repository.Provider
}

func NewManager() *Manager {
	newState := database.NewState()
	return &Manager{
		repository: repository.NewProvider(newState),
	}
}

func (m *Manager) NewServer() *httptest.Server {
	// Increment the server count and create a new server.
	m.serverCount++

	// Get the server ID.
	serverID := m.serverCount

	// Create a new test server.
	ts := httptest.NewServer(handler.NewMux(serverID, m.repository))

	slog.Info("Test server for examples created!", "url", ts.URL, "serverID", serverID)

	return ts
}
