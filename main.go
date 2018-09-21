package main

import (
	_ "encoding/json"
	"fmt"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"os"
	. "simplestorage/BucketManagement"
	. "simplestorage/ObjectManagement"
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
	var buckObj = fmt.Sprintf("%s%s",bucketPath,objectPath)

	/* Bucket Management */
	router.HandleFunc(bucketPath, CreateBucket).Queries("create","{create}").Methods("POST")
	router.HandleFunc(bucketPath, DeleteBucket).Queries("delete","{delete}").Methods("DELETE")
	router.HandleFunc(bucketPath, ListBucket).Queries("list","{list}").Methods("GET")

	/* Object Management */
	router.HandleFunc(buckObj, CreateTicket).Queries("create","{create}").Methods("POST")
	router.HandleFunc(buckObj, UploadPart).Queries("partNumber","{partNumber}").Methods("PUT")
	router.HandleFunc(buckObj, CompleteUpload).Queries("complete","{complete}").Methods("POST")
	router.HandleFunc(buckObj, DeletePart).Queries("delete","{delete}","partNumber","{partNumber}").Methods("DELETE")
	router.HandleFunc(buckObj, DeleteObject).Queries("delete","{delete}").Methods("DELETE")
	router.HandleFunc(buckObj, GetEndPoints).Methods("GET")
	router.HandleFunc(buckObj, UpdateMeta).Queries("metadata","{metadata}","key","{key}").Methods("PUT")
	router.HandleFunc(buckObj, DeleteMeta).Queries("metadata","{metadata}","key","{key}").Methods("DELETE")
	router.HandleFunc(buckObj, GetEndPoints).Queries("metadata","{metadata}","key","{key}").Methods("GET")
	router.HandleFunc(buckObj, GetEndPoints).Queries("metadata","{metadata}").Methods("GET")

	srv := &http.Server{
		Handler:      handlers.LoggingHandler(os.Stdout, router),
		Addr:         ":8080",
	}

	log.Fatal(srv.ListenAndServe())
	//log.Fatal(http.ListenAndServe("localhost:8080", router))

}