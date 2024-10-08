package handler

import (
	"encoding/json"
	"log/slog"
	"math/rand"
	"net/http"
	"strconv"
	"time"

	"github.com/opus-domini/fast-shot/examples/server/config"
	"github.com/opus-domini/fast-shot/examples/server/model"
	"github.com/opus-domini/fast-shot/examples/server/repository"
)

type (
	// Middleware is a function that wraps an HTTP handler.
	Middleware func(http.Handler) http.Handler

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

func NewMux(config *config.Server, repository *repository.Provider) http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/", HealthCheck(config.ID))
	mux.HandleFunc("GET /users", GetAll(repository.User))
	mux.HandleFunc("GET /users/{id}", GetByID(repository.User))
	mux.HandleFunc("POST /users", Create[*model.User](repository.User))
	mux.HandleFunc("GET /resources", GetAll(repository.Resource))
	mux.HandleFunc("GET /resources/{id}", GetByID(repository.Resource))
	mux.HandleFunc("POST /resources", Create[*model.Resource](repository.Resource))

	middlewares := []Middleware{
		HeaderDebugMiddleware(config),
		BusyMiddleware(config),
	}

	return handleWithMiddlewares(middlewares, mux)
}

func HeaderDebugMiddleware(config *config.Server) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if config.EnableHeaderDebug {
				slog.Info("Request Headers:", "headers", r.Header)
			}
			next.ServeHTTP(w, r)
		})
	}
}

func BusyMiddleware(config *config.Server) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if config.EnableBusy && shouldSimulateServerError() {
				slog.Error("Server is busy!", "serverID", config.ID)
				writeErrorResponse(w, ErrorMessage{Status: http.StatusInternalServerError, Message: "Internal Server Error"})
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

func HealthCheck(serverID int) func(w http.ResponseWriter, _ *http.Request) {
	return func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(HealthCheckMessage{
			ServerID:  strconv.Itoa(serverID),
			Status:    "UP",
			Timestamp: time.Now().Format(time.RFC3339),
		})
	}
}

func GetAll(repository repository.Repository) func(w http.ResponseWriter, _ *http.Request) {
	return func(w http.ResponseWriter, _ *http.Request) {
		_ = json.NewEncoder(w).
			Encode(repository.GetAll())
	}
}

func GetByID(repository repository.Repository) func(w http.ResponseWriter, _ *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
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

func Create[T model.Model](repository repository.Repository) func(w http.ResponseWriter, _ *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var raw T
		if err := json.NewDecoder(r.Body).Decode(&raw); err != nil {
			writeErrorResponse(w, ErrorMessage{Status: http.StatusUnprocessableEntity, Message: "Invalid request body"})
			return
		}

		newResource := repository.Create(raw)

		w.WriteHeader(http.StatusCreated)
		_ = json.NewEncoder(w).Encode(newResource)
	}
}

// handleWithMiddlewares applies a list of middlewares to an HTTP handler in reverse order.
// Middlewares are applied in reverse order because each middleware wraps the next one.
// The last middleware in the list is the first to be executed, and it calls the next middleware in the chain.
// This continues until the final handler is reached.
func handleWithMiddlewares(middlewares []Middleware, mux http.Handler) http.Handler {
	for i := len(middlewares) - 1; i >= 0; i-- {
		mux = middlewares[i](mux)
	}
	return mux
}

func writeErrorResponse(w http.ResponseWriter, errorMessage ErrorMessage) {
	w.WriteHeader(errorMessage.Status)
	_ = json.NewEncoder(w).Encode(errorMessage)
}

func shouldSimulateServerError() bool {
	// Simulate occasional server errors for retry examples
	return rand.Intn(100) < 80
}
