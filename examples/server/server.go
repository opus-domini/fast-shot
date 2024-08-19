package server

import (
	"log/slog"
	"net/http/httptest"
	"sync"

	"github.com/opus-domini/fast-shot/examples/server/config"
	"github.com/opus-domini/fast-shot/examples/server/database"
	"github.com/opus-domini/fast-shot/examples/server/handler"
	"github.com/opus-domini/fast-shot/examples/server/repository"
)

type (
	// Manager manages the shared database and the creation of test servers.
	Manager struct {
		running    []config.Server
		repository *repository.Provider
		mutex      sync.Mutex
	}

	// ServerBuilder helps in building a server with custom configurations.
	ServerBuilder struct {
		manager *Manager
		config  *config.Server
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
	m.mutex.Lock()
	defer m.mutex.Unlock()

	ts := httptest.NewServer(handler.NewMux(config, m.repository))
	config.ID = m.generateServerID()
	config.URL = ts.URL
	m.running = append(m.running, *config)
	slog.Info("Test server created!", "config", config)
	return ts
}

func (m *Manager) NewServer() *httptest.Server {
	return m.newServer(&config.Server{})
}

// NewServerBuilder creates a new instance of ServerBuilder.
func (m *Manager) NewServerBuilder() *ServerBuilder {
	return &ServerBuilder{
		manager: m,
		config:  &config.Server{},
	}
}

// EnableBusy sets the EnableBusy flag.
func (b *ServerBuilder) EnableBusy() *ServerBuilder {
	b.config.EnableBusy = true
	return b
}

// EnableHeaderDebug sets the EnableHeaderDebug flag.
func (b *ServerBuilder) EnableHeaderDebug() *ServerBuilder {
	b.config.EnableHeaderDebug = true
	return b
}

// Build creates a new server with the specified configurations.
func (b *ServerBuilder) Build() *httptest.Server {
	return b.manager.newServer(b.config)
}
