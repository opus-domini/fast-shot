package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/opus-domini/fast-shot/examples/server/model"
	"github.com/opus-domini/fast-shot/examples/server/repository"
)

func GetTuples(w http.ResponseWriter, _ *http.Request) {
	// Simulate occasional server errors for retry examples
	if generateServerError() {
		errorResponse(w, ErrorMessage{Status: http.StatusInternalServerError, Message: "Internal Server Error"})
		return
	}

	_ = json.NewEncoder(w).Encode(repository.Tuple().GetAll())
}

func GetTuple(w http.ResponseWriter, r *http.Request) {
	// Simulate occasional server errors for retry examples
	if generateServerError() {
		errorResponse(w, ErrorMessage{Status: http.StatusInternalServerError, Message: "Internal Server Error"})
		return
	}

	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil || id < 0 {
		errorResponse(w, ErrorMessage{Status: http.StatusBadRequest, Message: "Invalid Tuple ID"})
		return
	}

	tuple, found := repository.Tuple().GetById(uint(id))
	if !found {
		errorResponse(w, ErrorMessage{Status: http.StatusNotFound, Message: "Tuple not found"})
		return
	}

	_ = json.NewEncoder(w).Encode(tuple)
	return
}

func CreateTuple(w http.ResponseWriter, r *http.Request) {
	// Simulate occasional server errors for retry examples
	if generateServerError() {
		errorResponse(w, ErrorMessage{Status: http.StatusInternalServerError, Message: "Internal Server Error"})
		return
	}

	var tuple *model.Tuple
	if err := json.NewDecoder(r.Body).Decode(&tuple); err != nil {
		errorResponse(w, ErrorMessage{Status: http.StatusUnprocessableEntity, Message: "Invalid Tuple request body"})
		return
	}

	newTuple := repository.Tuple().Create(tuple)

	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(newTuple)
}
