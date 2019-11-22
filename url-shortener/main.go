package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/teris-io/shortid"
)

// ShortenedURL data type
type ShortenedURL struct {
	Shortened string
	Original  string
}

func generateUniqueID() (id string, err error) {
	sid, err := shortid.New(1, shortid.DefaultABC, 2342)
	if err != nil {
		panic(err)
	}

	uniqueID, err := sid.Generate()
	return uniqueID, err
}

func handlePostRequests(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	type RequestBody struct {
		Original string
	}
	var body RequestBody
	err := decoder.Decode(&body)
	if err != nil {
		panic(err)
	}
	fmt.Println(body.Original)

	id, err := generateUniqueID()
	if err != nil {
		panic(err)
	}

	new := ShortenedURL{id, body.Original}
	fmt.Printf("Created new item %s", new.Shortened)
}

func setupRoutes() {
	fmt.Println("Listening on port 8080")

	// setup routes
	http.HandleFunc("/short", handlePostRequests)

	// serve routes
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func main() {
	setupRoutes()
}
