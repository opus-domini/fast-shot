package main

import (
	"log/slog"
	"time"

	fastshot "github.com/opus-domini/fast-shot"
	"github.com/opus-domini/fast-shot/examples/server"
	"github.com/opus-domini/fast-shot/examples/server/model"
)

func main() {
	ts := server.Start()
	defer ts.Close()

	// Create a default client with the server URL.
	client := fastshot.DefaultClient(ts.URL)

	// Create a new user.
	newUser := &model.User{
		Name:      "John",
		Birthdate: time.Date(1990, 1, 1, 0, 0, 0, 0, time.UTC),
	}

	// Send the new user to the server.
	resp, err := client.POST("/users").
		Body().AsJSON(newUser).
		Send()

	// Check if there was an error sending the request.
	if err != nil {
		slog.Error("Error sending the request.", "error", err)
		return
	}

	// Check if the response is an error.
	if resp.Status().Is5xxServerError() {
		defer resp.Body().Close()
		slog.Error("Failed to create user, server error.", "status", resp.Status().Text())
		return
	}

	// Check if the response is a client error.
	if resp.Status().Is4xxClientError() {
		defer resp.Body().Close()
		slog.Error("Failed to create user, some client error.", "status", resp.Status().Text())
		return
	}

	// Congratulations! The user was created.
	slog.Info("User created!", "status", resp.Status().Text())

	// Parse the response body.
	var createdUser *model.User

	// Don't need to close the response body here.
	// It's done automatically when using AsBytes, AsString or AsJSON methods.
	if parseErr := resp.Body().AsJSON(&createdUser); parseErr != nil {
		slog.Error("Error parsing response.", "error", parseErr)
		return
	}

	// Print the created user.
	slog.Info("User data:", "data", createdUser)
}
