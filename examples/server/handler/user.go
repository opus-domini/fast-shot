package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/opus-domini/fast-shot/examples/server/model"
	"github.com/opus-domini/fast-shot/examples/server/repository"
)

func GetUsers(w http.ResponseWriter, _ *http.Request) {
	_ = json.NewEncoder(w).Encode(repository.User().GetAll())
}

func GetUser(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil || id < 0 {
		http.Error(w, "Invalid User ID", http.StatusBadRequest)
		return
	}

	user, found := repository.User().GetById(uint(id))
	if !found {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	_ = json.NewEncoder(w).Encode(user)
	return
}

func CreateUser(w http.ResponseWriter, r *http.Request) {
	var user *model.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "Invalid User request body", http.StatusBadRequest)
		return
	}

	repository.User().Create(user)
	w.WriteHeader(http.StatusNoContent)
	return
}
