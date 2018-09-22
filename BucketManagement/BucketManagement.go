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

	if !ValidatePattern(bucketName,BuckNamePattern){
		log.Print("Invalid Bucket name!")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

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
		//log.Printf("Added %s to MongoDB",bucketName)
		mkdir := MakeBucketDirectory(bucketName)

		if add && mkdir {
			var tmpBucket TempBucket
			mapstructure.Decode(tmp,&tmpBucket)
			w.Header().Set("Content-Type","application/json")
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(tmpBucket)
			log.Printf("Bucket %s created",bucketName)
			return
		}else{
			log.Printf("Errro creating bucket %s",bucketName)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	}else{
		log.Printf("Bucket %s not found",bucketName)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
}

func DeleteBucket(w http.ResponseWriter, r *http.Request) {
	var bucketName = GetBucketName(r)
	if CheckBucketExist(bucketName) {
		//log.Print("Bucket Exists")
		rm := RemoveBucketDirectory(bucketName)
		del := RemoveBucket(bucketName)
		log.Printf("Removing bucket %s",bucketName)
		if ! (rm && del){
			w.WriteHeader(http.StatusBadRequest)
			log.Printf("Error removing %s",bucketName)
			return
		}
		log.Printf("Bucket %s removed",bucketName)
		w.WriteHeader(http.StatusOK)
		return
	}else{
		//log.Print("Bucket NOT Exists")
		log.Printf("Bucket %s not found",bucketName)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
}

func ListBucket(w http.ResponseWriter, r *http.Request) {
	var bucketName = GetBucketName(r)

	var result map[string]interface{}
	result = make(map[string]interface{})

	var bucket = GetReturnBucket(bucketName)

	if bucket.Name == "" {
		log.Print("Invalid bucket name")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	mapstructure.Decode(bucket,&result)

	objects := GetObjectList(bucketName)

	if objects != nil {
		result["objects"] = objects
	}else{
		result["objects"] = []string{}
	}
	log.Printf("Listing bucket %s",bucketName)
	w.Header().Set("Content-Type","application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(result)
	return
}