package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/opus-domini/fast-shot/examples/server/repository"
)

func GetResources(w http.ResponseWriter, _ *http.Request) {
	_ = json.NewEncoder(w).Encode(repository.Resource().GetAll())
}

func handleResource(w http.ResponseWriter, _ *http.Request) {
	// Simulate occasional server errors for retry examples
	if time.Now().UnixNano()%2 == 0 {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	// Simulate successful response
	_ = json.NewEncoder(w).
		Encode(
			map[string]string{
				"message": "Resource data",
			},
		)
}

func handleLoadData(w http.ResponseWriter, _ *http.Request) {
	// Simulate responses from different "servers" for load balancing
	serverID := (time.Now().UnixNano() / 1e6) % 2
	// Simulate data from different servers
	_ = json.NewEncoder(w).
		Encode(
			map[string]string{
				"message": fmt.Sprintf("Data from Server %d", serverID),
			},
		)
}
