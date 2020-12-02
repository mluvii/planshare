package main

import (
	"context"
	"encoding/json"
	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

var topCtx = context.Background()

var rdb = func() *redis.Client {
	addr, configured := os.LookupEnv("REDIS_ADDR")
	if !configured {
		addr = "localhost:6379"
	}

	db, dbEnvErr := strconv.Atoi(os.Getenv("REDIS_DB"))
	if dbEnvErr != nil {
		db = 0
	}

	return redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: os.Getenv("REDIS_PASSWORD"),
		DB:       db,
	})
}()

func handleSave(res http.ResponseWriter, req *http.Request) {
	if req.Method != "POST" {
		res.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		res.WriteHeader(http.StatusBadRequest)
		return
	}

	var js json.RawMessage
	err = json.Unmarshal(body, &js)
	if err != nil {
		log.Println("Cannot unmarshal json", err)
		res.WriteHeader(http.StatusBadRequest)
		return
	}

	ctx := context.WithValue(topCtx, 0, req.Header.Get("X-Request-ID"))

	id := strings.ToLower(strings.Replace(uuid.New().String(), "-", "", -1))

	err = rdb.HSet(ctx, "planshare", id, body).Err()
	if err != nil {
		log.Println("Failed to save", err)
		res.WriteHeader(http.StatusInternalServerError)
		return
	}

	responseBytes, _ := json.Marshal(&struct {
		Id string `json:"id"`
	}{
		Id: id,
	})

	res.WriteHeader(http.StatusCreated)
	res.Header().Set("Content-Type", "application/json")
	res.Write(responseBytes)
}

func handleLoad(res http.ResponseWriter, req *http.Request) {
	if req.Method != "GET" {
		res.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	vars := mux.Vars(req)
	id, found := vars["id"]
	if !found {
		http.NotFound(res, req)
		return
	}

	ctx := context.WithValue(topCtx, 0, req.Header.Get("X-Request-ID"))

	val, err := rdb.HGet(ctx, "planshare", id).Result()
	if err != nil {
		log.Println("Failed to load", err)
		http.NotFound(res, req)
		return
	}

	var js json.RawMessage
	err = json.Unmarshal([]byte(val), &js)
	if err != nil {
		log.Println("Cannot unmarshal json", err)
		http.NotFound(res, req)
		return
	}

	res.Header().Set("Content-Type", "application/json")
	res.Write([]byte(val))
}

func listenProxy() {
	r := mux.NewRouter()
	r.HandleFunc("/{id}", handleLoad).Methods("GET")
	r.HandleFunc("/", handleSave).Methods("POST")

	port := os.Getenv("API_PORT")
	if port == "" {
		port = "8080"
	}

	server := &http.Server{
		Addr:         "127.0.0.1:" + port,
		Handler:      r,
		ReadTimeout:  time.Duration(60) * time.Second,
		WriteTimeout: time.Duration(60) * time.Second,
	}

	err := server.ListenAndServe()
	if err != nil {
		log.Fatalf("Failed to listen and serve: %s\n", err.Error())
	}
}

func main() {
	listenProxy()
}
