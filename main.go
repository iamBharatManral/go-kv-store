package main

import (
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/v1/{key}", PutHandler).Methods("PUT")
	r.HandleFunc("/v1/key/{key}", GetHandler).Methods("GET")
	r.HandleFunc("/v1/key/{key}", DeleteHandler).Methods("DELETE")
	log.Fatal(http.ListenAndServe(":8080", r))
}
