package main

import (
	"net/http"
	"log"
	"context"

	"github.com/gorilla/mux"
    "github.com/go-redis/redis/v8"
)

var ctx = context.Background()
var rdb *redis.Client

// configuration
var address  = "localhost:6379"
var password = ""
var database = 0

func main() {
	rdb = redis.NewClient(&redis.Options{
		Addr: address,
		Password: password,
		DB: database,
	})

	// ping redis
	// make sure the client is connected to the server
	ping := rdb.Ping(ctx)

	if ping.Err() != nil {
		_, err := ping.Result()
		log.Fatalf("Could not connect to Redis!\n\nAddress: %s\nPassword: %s\nDatabase: %v\nError: %s", address, password, database, err.Error())
	} else {
		router := mux.NewRouter()

		router.HandleFunc("/", docs).Methods("GET")
		router.HandleFunc("/create", createLink).Methods("POST")
		router.HandleFunc("/delete", deleteLink).Methods("DELETE")
		router.HandleFunc("/clicks/{id}", getClicks).Methods("GET")
		router.HandleFunc("/{id}", getLink).Methods("GET")

		log.Fatal(http.ListenAndServe(":8000", router))
	}
}