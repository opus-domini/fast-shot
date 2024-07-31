package main

import (
	"log/slog"

	fastshot "github.com/opus-domini/fast-shot"
	"github.com/opus-domini/fast-shot/examples/server"
)

func main() {
	// Create a server manager
	serverManager := server.NewManager()

	// Start Test Server #1
	ts1 := serverManager.NewServer()
	defer ts1.Close()

	// Start Test Server #2
	ts2 := serverManager.NewServer()
	defer ts2.Close()

	// Start Test Server #3
	ts3 := serverManager.NewServer()
	defer ts3.Close()

	// Create a custom client with client-side load balancing.
	// The client round-robins the requests between the servers.
	client := fastshot.NewClientLoadBalancer([]string{ts1.URL, ts2.URL, ts3.URL}).Build()

	// Perform health checks on the servers
	for i := 0; i < 5; i++ {
		healthcheck(client)
	}
}

func healthcheck(client fastshot.ClientHttpMethods) {
	// Send a GET request to the root endpoint
	resp, err := client.GET("/").Send()

	// Check if there was an error sending the request.
	if err != nil {
		slog.Error("Error sending the request.", "error", err)
	}

	var healthCheckResponse map[string]interface{}

	// Parse the response body as JSON
	// Note: The response body is automatically closed when using AsBytes, AsString, or AsJSON methods
	if parseErr := resp.Body().AsJSON(&healthCheckResponse); parseErr != nil {
		slog.Error("Error parsing response.", "error", parseErr)
		return
	}

	slog.Info("Health Check received.", "data", healthCheckResponse)
}
