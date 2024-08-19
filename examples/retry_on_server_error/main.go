package main

import (
	"log/slog"
	"time"

	fastshot "github.com/opus-domini/fast-shot"
	"github.com/opus-domini/fast-shot/examples/server"
	"github.com/opus-domini/fast-shot/examples/server/model"
)

func main() {
	// Start the test server
	ts := server.NewManager().
		NewServerBuilder().
		EnableBusy().
		Build()
	// Close the server when the function ends.
	defer ts.Close()

	// Create a default client with the server URL.
	client := fastshot.DefaultClient(ts.URL)

	slog.Info("Get Resource:", "id", 2)

	// Request a resource in a busy server!
	resp, err := client.GET("/resources/2").
		Retry().SetConstantBackoff(50*time.Millisecond, 10).
		Send()

	// Check if there was an error sending the request.
	if err != nil {
		slog.Error("Error sending the request.", "error", err)
		return
	}

	// Check if the response is an error.
	if resp.Status().IsError() {
		defer resp.Body().Close()
		slog.Error("Failed to fetch a resource.", "status", resp.Status().Text())
		return
	}

	// Parse the response body.
	var resource *model.Resource

	// Don't need to close the response body here.
	// It's done automatically when using AsBytes, AsString or AsJSON methods.
	if parseErr := resp.Body().AsJSON(&resource); parseErr != nil {
		slog.Error("Error parsing response.", "error", parseErr)
		return
	}

	// Congratulations! The resource is here.
	slog.Info("Resource got!", "data", resource)
}
