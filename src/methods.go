package main

import (
	"net/http"
	"io"
	"time"
	"crypto/sha256"
	"encoding/json"
	"encoding/base64"

	"github.com/gorilla/mux"
	"github.com/go-redis/redis/v8"
)

type Link struct {
	ID string `json:"id"`
	Location string `json:"location"`
	Clicks int `json:"clicks"`
	Password string `json:"password"`
	Created string `json:"created"`
}

func getLink(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	val, redis_err := rdb.HGet(ctx, params["id"], "L").Result()
	
	if redis_err == redis.Nil {
		w.Header().Set("Content-Type", "text/plain")
		io.WriteString(w, "Link doesn't exist")

		return
	} else if redis_err != nil {
		panic(redis_err)

		w.Header().Set("Content-Type", "text/plain")
		io.WriteString(w, "Something went wrong")

		return
	} 

	http.Redirect(w, r, val, 301)

	incr_err := rdb.HIncrBy(ctx, params["id"], "C", 1).Err()

	if incr_err != nil {
		panic(incr_err)
	}
}

func createLink(w http.ResponseWriter, r *http.Request) {
	var link Link
	json.NewDecoder(r.Body).Decode(&link)

	if exists, _ := rdb.Exists(ctx, link.ID).Result(); exists == 1 {
		w.Header().Set("Content-Type", "text/plain")
		io.WriteString(w, "Link " + link.ID + " already exists!")

		return
	}
	
	hashedpass := sha256.Sum256([]byte(link.Password))
	link.Password = base64.StdEncoding.EncodeToString(hashedpass[:])

	link.Created = time.Now().UTC().String()

	err := rdb.HSet(ctx, link.ID, "L", link.Location, "C", 0, "P", link.Password, "T", link.Created).Err()
	if err != nil {
		panic(err)

		w.Header().Set("Content-Type", "text/plain")	
		io.WriteString(w, "Something went wrong!")

		return
	}

	w.Header().Set("Content-Type", "application/json")	
	json.NewEncoder(w).Encode(link)
}

func deleteLink(w http.ResponseWriter, r *http.Request) {
	var link Link
	json.NewDecoder(r.Body).Decode(&link)

	hashedpass := sha256.Sum256([]byte(link.Password))
	link.Password = base64.StdEncoding.EncodeToString(hashedpass[:])

	val, err := rdb.HGet(ctx, link.ID, "P").Result()

	if err != nil {
		w.Header().Set("Content-Type", "text/plain")
		io.WriteString(w, "Failed to delete link")
		
		return
	}

	if val == link.Password {
		rdb.Del(ctx, link.ID)
		w.Header().Set("Content-Type", "text/plain")
		io.WriteString(w, link.ID + " Deleted")
	} else {
		w.Header().Set("Content-Type", "text/plain")
		io.WriteString(w, "Incorrect password")
	}

}

func getClicks(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain")

	params := mux.Vars(r)
	val, err := rdb.HGet(ctx, params["id"], "C").Result()

	if err == redis.Nil {
		io.WriteString(w, "That link doesn't exist!")
	} else if err != nil {
		panic(err)
		io.WriteString(w, "Failed to get clicks")
	} else {
		io.WriteString(w, val)
	}

}

func docs(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain")

	io.WriteString(w, `
	Welcome to my link shortener
	Available routes:
	GET /{id} -> redirect to the original url
	POST /create -> create a link
		{
			"ID": "this is my link id",
			"Location": "https://www.youtube.com/watch?v=dQw4w9WgXcQ",
			"Password": "this is used to delete your link"
		}
	DELETE /delete -> delete a link
		{
			"ID": "id of the link to be deleted",
			"Password": "so only you can delete it!"
		}
	GET /clicks/{id} -> number of times a link was clicked
	`)
}