package main

import (
	. "./BucketManagement"
	. "./ObjectManagement"
	_ "encoding/json"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

// our main function
func main() {
	router := mux.NewRouter()
	/* Bucket Management */
	router.HandleFunc("/{bucketName}", CreateBucket).Queries("create","{create}").Methods("POST")
	router.HandleFunc("/{bucketName}", DeleteBucket).Queries("delete","{delete}").Methods("DELETE")
	router.HandleFunc("/{bucketName}", ListBucket).Queries("list","{list}").Methods("GET")

	/* Object Management */
	router.HandleFunc("/{bucketName}/{objectName}", CreateTicket).Queries("create","{create}").Methods("POST")
	router.HandleFunc("/{bucketName}/{objectName}", UploadAll).Queries("partNumber","{partNumber}").Methods("PUT")

	log.Fatal(http.ListenAndServe(":8000", router))
}