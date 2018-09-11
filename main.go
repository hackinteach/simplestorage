package main

import (
	. "./BucketManagement"
	_ "encoding/json"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

// our main function
func main() {
	router := mux.NewRouter()
	router.HandleFunc("/{bucketName}", CreateBucket).Queries("create","{*}").Methods("POST")
	router.HandleFunc("/{bucketName}", DeleteBucket).Queries("delete","{*}").Methods("DELETE")
	router.HandleFunc("/{bucketName}", ListBucket).Queries("list","{*}").Methods("GET")

	log.Fatal(http.ListenAndServe(":8000", router))
}

