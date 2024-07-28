package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/opus-domini/fast-shot/examples/server/model"
	"github.com/opus-domini/fast-shot/examples/server/repository"
)

func GetResources(w http.ResponseWriter, _ *http.Request) {
	_ = json.NewEncoder(w).Encode(repository.Resource().GetAll())
}

func GetResource(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil || id < 0 {
		errorResponse(w, ErrorMessage{Status: http.StatusBadRequest, Message: "Invalid Resource ID"})
		return
	}

	resource, found := repository.Resource().GetById(uint(id))
	if !found {
		errorResponse(w, ErrorMessage{Status: http.StatusNotFound, Message: "Resource not found"})
		return
	}

	_ = json.NewEncoder(w).Encode(resource)
	return
}

func CreateResource(w http.ResponseWriter, r *http.Request) {
	var resource *model.Resource
	if err := json.NewDecoder(r.Body).Decode(&resource); err != nil {
		errorResponse(w, ErrorMessage{Status: http.StatusUnprocessableEntity, Message: "Invalid Resource request body"})
		return
	}

	newResource := repository.Resource().Create(resource)

	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(newResource)
}
