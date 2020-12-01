package main

import (
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

func handleSave(res http.ResponseWriter, req *http.Request) {
	if req.Method != "POST" {
		res.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
}

func handleLoad(res http.ResponseWriter, req *http.Request) {
	if req.Method != "GET" {
		res.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
}

func listenProxy() {
	r := mux.NewRouter()
	r.HandleFunc("/save", handleSave)
	r.HandleFunc("/{:id}", handleLoad)

	server := &http.Server{
		Addr:         "127.0.0.1:8080",
		Handler:      r,
		ReadTimeout:  time.Duration(60) * time.Second,
		WriteTimeout: time.Duration(60) * time.Second,
	}

	error := server.ListenAndServe()
	if error != nil {
		log.Fatalf("Failed to listen and serve: %s\n", error.Error())
	}
}

func main() {
	listenProxy()
}
