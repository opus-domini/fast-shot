package handler

import (
	"encoding/json"
	"log/slog"
	"math/rand"
	"net/http"
)

func NewMux() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /users", GetUsers)
	mux.HandleFunc("GET /users/{id}", GetUser)
	mux.HandleFunc("POST /users", CreateUser)
	mux.HandleFunc("GET /resources", GetResources)
	mux.HandleFunc("GET /resources/{id}", GetResource)
	mux.HandleFunc("POST /resources", CreateResource)
	mux.HandleFunc("GET /tuples", GetTuples)
	mux.HandleFunc("GET /tuples/{id}", GetTuple)
	mux.HandleFunc("POST /tuples", CreateTuple)
	// Handler default and throw 404
	mux.HandleFunc("/", http.NotFound)
	return mux
}

type ErrorMessage struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
}

func errorResponse(w http.ResponseWriter, errorMessage ErrorMessage) {
	w.WriteHeader(errorMessage.Status)
	_ = json.NewEncoder(w).Encode(errorMessage)
}

func generateServerError() bool {
	// Simulate occasional server errors for retry examples
	busy := rand.Intn(100) < 80
	if busy {
		slog.Warn("Server is busy... try again later.")
	}
	return busy
}
