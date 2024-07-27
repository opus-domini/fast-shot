package handler

import (
	"encoding/json"
	"github.com/opus-domini/fast-shot/examples/server/repository"
	"net/http"
)

func GetTuples(w http.ResponseWriter, _ *http.Request) {
	_ = json.NewEncoder(w).Encode(repository.Tuple().GetAll())
}
