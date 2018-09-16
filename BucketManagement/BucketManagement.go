package BucketManagement

import (
	"encoding/json"
	"github.com/mitchellh/mapstructure"
	"log"
	"net/http"
	. "simplestorage/Misc"
	. "simplestorage/Mongo"
	. "simplestorage/Structure"
)

func CreateBucket(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type","application/json")
	var bucket Bucket

	bucketName := GetBucketName(r)

	var tmp map[string]interface{}
	tmp = make(map[string]interface{})
	tmp["name"] = bucketName
	tmp["created"] = GetTime()
	tmp["modified"] = GetTime()

	mapstructure.Decode(tmp,&bucket)

	if ! CheckBucketExist(bucketName){
		// Bucket can be create
		log.Printf("Creating %s",bucketName)
		add := AddBucket(bucket)
		log.Printf("Added %s to MongoDB",bucketName)
		mkdir := MakeBucketDirectory(bucketName)

		if add && mkdir {
			var tmpBucket TempBucket
			mapstructure.Decode(tmp,&tmpBucket)
			json.NewEncoder(w).Encode(tmpBucket)
			w.Header().Set("Content-Type","application/json")
			w.WriteHeader(http.StatusOK)
		}else{
			w.WriteHeader(http.StatusBadRequest)
		}
	}else{
		w.WriteHeader(http.StatusBadRequest)
	}
}

func DeleteBucket(w http.ResponseWriter, r *http.Request) {
	var bucketName = GetBucketName(r)
	if CheckBucketExist(bucketName) {
		log.Print("Bucket Exists")
		rm := RemoveDirectory(bucketName)
		del := RemoveBucket(bucketName)

		if ! (rm && del){
			w.WriteHeader(http.StatusBadRequest)
		}
		w.WriteHeader(http.StatusOK)
	}else{
		log.Print("Bucket NOT Exists")
		w.WriteHeader(http.StatusBadRequest)
	}
}

func ListBucket(w http.ResponseWriter, r *http.Request) {
	var bucketName = GetBucketName(r)

	var result map[string]interface{}
	result = make(map[string]interface{})

	var bucket = GetBucket(bucketName)

	if bucket.Name == "" {
		w.WriteHeader(http.StatusBadRequest)
	}

	mapstructure.Decode(bucket,&result)

	objects := GetObjectList(bucketName)

	if objects != nil {
		result["objects"] = objects
	}else{
		result["objects"] = []string{}
	}

	w.Header().Set("Content-Type","application/json")
	json.NewEncoder(w).Encode(result)
	w.WriteHeader(http.StatusOK)
}