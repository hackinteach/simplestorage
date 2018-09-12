package BucketManagement

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"gopkg.in/mgo.v2/bson"
	"net/http"
)

type Bucket struct {
	ID bson.ObjectId `bson:"_id,omitempty"`
	Name string 	`bson:"name" json:"bucketName"`
	Created int64	`bson:"created" json:"created"`
	Modified int64	`bson:"modified" json:"modified"`
}

func CreateBucket(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	json.NewEncoder(w).Encode(vars)
}

func DeleteBucket(w http.ResponseWriter, r *http.Request) {

}

func ListBucket(w http.ResponseWriter, r *http.Request) {

}