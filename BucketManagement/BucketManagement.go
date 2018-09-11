package BucketManagement

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
)

func CreateBucket(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	bucketName := vars["bucketName"]
	json.NewEncoder(w).Encode(bucketName)
}

func DeleteBucket(w http.ResponseWriter, r *http.Request) {

}

func ListBucket(w http.ResponseWriter, r *http.Request) {

}