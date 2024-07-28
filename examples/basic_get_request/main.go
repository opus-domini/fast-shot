package main

import (
	"log/slog"

	fastshot "github.com/opus-domini/fast-shot"
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
	handleResponse(resp, &[]model.User{})
}

func getUser(client fastshot.ClientHttpMethods, id string) {
	slog.Info("Get User:", "id", id)

	resp, err := client.GET("/users/" + id).Send()
	if err != nil {
		slog.Error("Error getting response.", "error", err)
		return
	}
	handleResponse(resp, &model.User{})
}

func handleResponse(resp *fastshot.Response, data interface{}) {
	slog.Info("Response:", "status", resp.Status().Text())

	if resp.Status().IsError() {
		slog.Error("Failed to get data.")
		return
	}

	// Don't need to close the response body here.
	// It's done automatically when using AsBytes, AsString or AsJSON methods.
	if parseErr := resp.Body().AsJSON(data); parseErr != nil {
		slog.Error("Error parsing response.", "error", parseErr)
		return
	}

	slog.Info("Data received!", "data", data)
}
