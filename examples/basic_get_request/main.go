package main

import (
	"log/slog"

	fastshot "github.com/opus-domini/fast-shot"
	"github.com/opus-domini/fast-shot/examples/server"
	"github.com/opus-domini/fast-shot/examples/server/model"
)

func main() {
	// Start the test server
	ts := server.NewManager().NewServer()
	defer ts.Close()

	// Create a default client with the server URL.
	client := fastshot.DefaultClient(ts.URL)

	// Get all users.
	getUsers(client)

	// Get a user by ID.
	getUser(client, "1")

	// Get a user that does not exist.
	getUser(client, "99")
}

func getUsers(client fastshot.ClientHttpMethods) {
	slog.Info("Get all Users.")

	resp, err := client.GET("/users").Send()
	if err != nil {
		slog.Error("Error getting response.", "error", err)
		return
	}

	if resp.Status().IsError() {
		slog.Error("Failed to get data.", "status", resp.Status().Text())
		return
	}

	// Decode straight into a typed value (Go 1.27 generic method).
	// The response body is closed automatically by AsJSON.
	users, err := resp.Body().AsJSON[[]model.User]()
	if err != nil {
		slog.Error("Error parsing response.", "error", err)
		return
	}

	slog.Info("Data received!", "data", users)
}

func getUser(client fastshot.ClientHttpMethods, id string) {
	slog.Info("Get User:", "id", id)

	resp, err := client.GET("/users/" + id).Send()
	if err != nil {
		slog.Error("Error getting response.", "error", err)
		return
	}

	if resp.Status().IsError() {
		slog.Error("Failed to get data.", "status", resp.Status().Text())
		return
	}

	// Decode straight into a typed value (Go 1.27 generic method).
	// The response body is closed automatically by AsJSON.
	user, err := resp.Body().AsJSON[model.User]()
	if err != nil {
		slog.Error("Error parsing response.", "error", err)
		return
	}

	slog.Info("Data received!", "data", user)
}
