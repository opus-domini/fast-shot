package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/opus-domini/fast-shot/examples/server/model"
	"github.com/opus-domini/fast-shot/examples/server/repository"
)

func GetUsers(w http.ResponseWriter, _ *http.Request) {
	_ = json.NewEncoder(w).
		Encode(repository.User().GetAll())
}

func GetUser(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil || id < 0 {
		errorResponse(w, ErrorMessage{Status: http.StatusBadRequest, Message: "Invalid User ID"})
		return
	}

	user, found := repository.User().GetById(uint(id))
	if !found {
		errorResponse(w, ErrorMessage{Status: http.StatusNotFound, Message: "User not found"})
		return
	}

	_ = json.NewEncoder(w).Encode(user)
	return
}

func CreateUser(w http.ResponseWriter, r *http.Request) {
	var user *model.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		errorResponse(w, ErrorMessage{Status: http.StatusUnprocessableEntity, Message: "Invalid User request body"})
		return
	}

	newUser := repository.User().Create(user)

	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(newUser)
}
