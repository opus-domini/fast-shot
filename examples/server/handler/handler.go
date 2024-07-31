package handler

import (
	"encoding/json"
	"log/slog"
	"math/rand"
	"net/http"
	"strconv"
	"time"

	"github.com/opus-domini/fast-shot/examples/server/model"
	"github.com/opus-domini/fast-shot/examples/server/repository"
)

type (
	// HealthCheckMessage is a simple struct for health check messages.
	HealthCheckMessage struct {
		ServerID  string `json:"server_id"`
		Status    string `json:"status"`
		Timestamp string `json:"timestamp"`
	}

	// ErrorMessage is a simple struct for error messages.
	ErrorMessage struct {
		Status  int    `json:"status"`
		Message string `json:"message"`
	}
)

func NewMux(serverID int, repository *repository.Provider) http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(HealthCheckMessage{
			ServerID:  strconv.Itoa(serverID),
			Status:    "UP",
			Timestamp: time.Now().Format(time.RFC3339),
		})
	})
	mux.HandleFunc("GET /users", GetAll(repository.User, false))
	mux.HandleFunc("GET /users/{id}", GetByID(repository.User, false))
	mux.HandleFunc("POST /users", Create(repository.User, false))
	mux.HandleFunc("GET /resources", GetAll(repository.Resource, false))
	mux.HandleFunc("GET /resources/{id}", GetByID(repository.Resource, false))
	mux.HandleFunc("POST /resources", Create(repository.Resource, false))
	mux.HandleFunc("GET /tuples", GetAll(repository.Tuple, true))
	mux.HandleFunc("GET /tuples/{id}", GetByID(repository.Tuple, true))
	mux.HandleFunc("POST /tuples", Create(repository.Tuple, true))
	return mux
}

func GetAll(repository repository.Repository, withBusyServer bool) func(w http.ResponseWriter, _ *http.Request) {
	return func(w http.ResponseWriter, _ *http.Request) {
		if withBusyServer && serverIsBusy() {
			writeErrorResponse(w, ErrorMessage{Status: http.StatusInternalServerError, Message: "Internal Server Error"})
			return
		}
		_ = json.NewEncoder(w).
			Encode(repository.GetAll())
	}
}

func GetByID(repository repository.Repository, withBusyServer bool) func(w http.ResponseWriter, _ *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if withBusyServer && serverIsBusy() {
			writeErrorResponse(w, ErrorMessage{Status: http.StatusInternalServerError, Message: "Internal Server Error"})
			return
		}

		id, err := strconv.Atoi(r.PathValue("id"))
		if err != nil || id < 0 {
			writeErrorResponse(w, ErrorMessage{Status: http.StatusBadRequest, Message: "Invalid ID"})
			return
		}

		resource, found := repository.GetById(uint(id))
		if !found {
			writeErrorResponse(w, ErrorMessage{Status: http.StatusNotFound, Message: "Not found"})
			return
		}

		_ = json.NewEncoder(w).Encode(resource)
		return
	}
}

func Create(repository repository.Repository, withBusyServer bool) func(w http.ResponseWriter, _ *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if withBusyServer && serverIsBusy() {
			writeErrorResponse(w, ErrorMessage{Status: http.StatusInternalServerError, Message: "Internal Server Error"})
			return
		}

		var raw model.Model
		if err := json.NewDecoder(r.Body).Decode(&raw); err != nil {
			writeErrorResponse(w, ErrorMessage{Status: http.StatusUnprocessableEntity, Message: "Invalid request body"})
			return
		}

		newResource := repository.Create(raw)

		w.WriteHeader(http.StatusCreated)
		_ = json.NewEncoder(w).Encode(newResource)
	}
}

func writeErrorResponse(w http.ResponseWriter, errorMessage ErrorMessage) {
	w.WriteHeader(errorMessage.Status)
	_ = json.NewEncoder(w).Encode(errorMessage)
}

func serverIsBusy() bool {
	// Simulate occasional server errors for retry examples
	if rand.Intn(100) < 80 {
		slog.Warn("Server is busy... try again later.")
		return true
	}
	return false
}
