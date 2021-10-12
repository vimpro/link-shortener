package main

import (
	"net/http"
	"log"

	"github.com/gorilla/mux"
	"context"
    "github.com/go-redis/redis/v8"
)

var ctx = context.Background()
var rdb *redis.Client

func main() {
	rdb = redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
		Password: "",
		DB: 0,
	})

	router := mux.NewRouter()
	router.HandleFunc("/", docs).Methods("GET")
	router.HandleFunc("/create", createLink).Methods("POST")
	router.HandleFunc("/delete", deleteLink).Methods("DELETE")
	router.HandleFunc("/clicks/{id}", getClicks).Methods("GET")
	router.HandleFunc("/{id}", getLink).Methods("GET")

	log.Fatal(http.ListenAndServe(":8000", router))
}