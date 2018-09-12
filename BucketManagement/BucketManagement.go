package BucketManagement

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
)

type Bucket struct {
	Name string 	`json:"bucketName"`
	Created int64	`json:"created"`
	Modified int64	`json:"modified"`
}

func CreateBucket(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	json.NewEncoder(w).Encode(vars)
}

func DeleteBucket(w http.ResponseWriter, r *http.Request) {

}

func ListBucket(w http.ResponseWriter, r *http.Request) {

}