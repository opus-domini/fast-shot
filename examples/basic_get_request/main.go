package main

import (
	"log/slog"

	"github.com/opus-domini/fast-shot"
	"github.com/opus-domini/fast-shot/examples/server"
	"github.com/opus-domini/fast-shot/examples/server/model"
)

func main() {
	ts := server.Start()
	defer ts.Close()

	// Create a default client with the server URL.
	client := fastshot.DefaultClient(ts.URL)

	// Get all users.
	getUsers(client)

	// Get a user by ID.
	getUser(client, "1")
}

func getUsers(client fastshot.ClientHttpMethods) {
	slog.Info("Get all Users.")

	resp, err := client.GET("/users").Send()
	if err != nil {
		slog.Error("Error getting response.", "error", err)
		return
	}

	slog.Info("Response received.", "status", resp.Status().Text())

	var users []model.User
	if parseErr := resp.Body().AsJSON(&users); parseErr != nil {
		slog.Error("Error parsing response.", "error", parseErr)
	}

	slog.Info("Users received!", "data", users)
}

func getUser(client fastshot.ClientHttpMethods, id string) {
	slog.Info("Get User.", "id", id)

	resp, err := client.GET("/users/" + id).Send()
	if err != nil {
		slog.Error("Error getting response.", "error", err)
		return
	}

	slog.Info("Response received.", "status", resp.Status().Text())

	var user model.User
	if parseErr := resp.Body().AsJSON(&user); parseErr != nil {
		slog.Error("Error parsing response.", "error", parseErr)
	}

	slog.Info("User received!", "data", user)
}
