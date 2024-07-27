package handler

import "net/http"

func NewMux() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /users", GetUsers)
	mux.HandleFunc("GET /users/{id}", GetUser)
	mux.HandleFunc("POST /users", CreateUser)
	mux.HandleFunc("GET /resource", GetResources)
	mux.HandleFunc("GET /resource/{id}", handleResource)
	mux.HandleFunc("POST /resource", handleResource)
	mux.HandleFunc("GET /tuples", GetTuples)
	mux.HandleFunc("GET /tuples/{id}", handleLoadData)
	mux.HandleFunc("POST /tuples", handleLoadData)
	// Handler default and throw 404
	mux.HandleFunc("/", http.NotFound)
	return mux
}
