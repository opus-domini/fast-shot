package main

import (
	"log/slog"
	"net/http"

	"github.com/opus-domini/fast-shot"
	"github.com/opus-domini/fast-shot/constant/mime"
	"github.com/opus-domini/fast-shot/examples/server"
	"github.com/opus-domini/fast-shot/examples/server/model"
)

func main() {
	ts := server.Start()
	defer ts.Close()

	// Create a custom client with the server URL.
	client := fastshot.NewClient(ts.URL).
		Header().AddUserAgent("MyAwesomeApp/1.0").
		Header().Add("X-My-Header", "MyValue").
		Cookie().Add(&http.Cookie{Name: "session_id", Value: "xyz123"}).
		Build()

	// Get all resources.
	resp, err := client.GET("/resources").
		Header().AddAccept(mime.JSON).
		Send()
		// Check if there was an error sending the request.
	if err != nil {
		slog.Error("Error sending the request.", "error", err)
	}

	var data []model.Resource

	// Don't need to close the response body here.
	// It's done automatically when using AsBytes, AsString or AsJSON methods.
	if parseErr := resp.Body().AsJSON(&data); parseErr != nil {
		slog.Error("Error parsing response.", "error", parseErr)
		return
	}

	slog.Info("Data received!", "data", data)
}
