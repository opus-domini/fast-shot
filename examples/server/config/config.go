package config

type (
	// Server represents a test server configuration.
	Server struct {
		ID     int
		URL    string
		IsBusy bool
	}
)
