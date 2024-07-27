package handler

import (
	"encoding/json"
	"net/http"

	"github.com/opus-domini/fast-shot/examples/server/repository"
)

func GetTuples(w http.ResponseWriter, _ *http.Request) {
	_ = json.NewEncoder(w).Encode(repository.Tuple().GetAll())
}
