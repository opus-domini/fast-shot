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

	slog.Info("Get Tuple:", "id", 2)

	// Request a tuple in a busy server!
	resp, err := client.GET("/tuples/2").
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
		slog.Error("Failed to fetch a tuple.", "status", resp.Status().Text())
		return
	}

	// Parse the response body.
	var tuple *model.Tuple

	// Don't need to close the response body here.
	// It's done automatically when using AsBytes, AsString or AsJSON methods.
	if parseErr := resp.Body().AsJSON(&tuple); parseErr != nil {
		slog.Error("Error parsing response.", "error", parseErr)
		return
	}

	// Congratulations! The tuple is here.
	slog.Info("Tuple got!", "tuple", tuple)
}
