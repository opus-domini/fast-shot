package server

import (
	"github.com/opus-domini/fast-shot/examples/server/config"
	"log/slog"
	"net/http/httptest"

	"github.com/opus-domini/fast-shot/examples/server/database"
	"github.com/opus-domini/fast-shot/examples/server/handler"
	"github.com/opus-domini/fast-shot/examples/server/repository"
)

type (
	// Manager manages the shared database and the creation of test servers.
	Manager struct {
		running    []config.Server
		repository *repository.Provider
	}
)

func NewManager() *Manager {
	newState := database.NewState()
	return &Manager{
		repository: repository.NewProvider(newState),
	}
}

func (m *Manager) generateServerID() int {
	return len(m.running) + 1
}

func (m *Manager) newServer(config *config.Server) *httptest.Server {
	ts := httptest.NewServer(handler.NewMux(config, m.repository))
	config.URL = ts.URL
	m.running = append(m.running, *config)
	slog.Info("Test server created!", "config", config)
	return ts
}

func (m *Manager) NewServer() *httptest.Server {
	return m.newServer(&config.Server{
		ID:     m.generateServerID(),
		IsBusy: false,
	})
}

func (m *Manager) NewBusyServer() *httptest.Server {
	return m.newServer(&config.Server{
		ID:     m.generateServerID(),
		IsBusy: true,
	})
}
