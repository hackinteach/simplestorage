package main

import (
	. "./BucketManagement"
	. "./ObjectManagement"
	_ "encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

const (
	//bucketPath = `/{bucketName: ^[0-9 | a-z | \- | _ ]*}`
	//objectPath = `/{objectName:^(?!\.)[0-9 | a-z | \- | _ | \.]*}`
	bucketPath = `/{bucketName}`
	objectPath = `/{objectName}`
)

// our main function
func main() {
	router := mux.NewRouter()

	/* Bucket Management */
	router.HandleFunc(bucketPath, CreateBucket).Queries("create","{create}").Methods("POST")
	router.HandleFunc(bucketPath, DeleteBucket).Queries("delete","{delete}").Methods("DELETE")
	router.HandleFunc(bucketPath, ListBucket).Queries("list","{list}").Methods("GET")

	/* Object Management */
	router.HandleFunc(fmt.Sprintf("%s%s",bucketPath,objectPath), CreateTicket).Queries("create","{create}").Methods("POST")
	router.HandleFunc(fmt.Sprintf("%s%s",bucketPath,objectPath), UploadAll).Queries("partNumber","{partNumber}").Methods("PUT")

	log.Fatal(http.ListenAndServe(":8000", router))
}