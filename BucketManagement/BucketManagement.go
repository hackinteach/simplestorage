package BucketManagement

import (
	"encoding/json"
	"github.com/mitchellh/mapstructure"
	"log"
	"net/http"
	. "simplestorage/Misc"
	. "simplestorage/Mongo"
	. "simplestorage/Structure"
	"strings"
)

func CreateBucket(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type","application/json")
	var bucket Bucket

	bucketName := strings.Replace(r.URL.Path,"/","",1)

	var tmp map[string]interface{}
	tmp = make(map[string]interface{})
	tmp["name"] = bucketName
	tmp["created"] = GetTime()
	tmp["modified"] = GetTime()

	mapstructure.Decode(tmp,&bucket)



	if ! CheckBucketExist(bucketName){
		// Bucket can be create
		log.Printf("Creating %s",bucketName)
		AddBucket(bucket)
		log.Printf("Added %s to MongoDB",bucketName)
		MakeDirectory(bucketName)

		var tmpBucket TempBucket
		mapstructure.Decode(tmp,&tmpBucket)
		json.NewEncoder(w).Encode(tmpBucket)
		w.WriteHeader(200)
	}else{
		log.Printf("Bucket already exists")
		w.WriteHeader(400)
	}
}


func DeleteBucket(w http.ResponseWriter, r *http.Request) {

}

func ListBucket(w http.ResponseWriter, r *http.Request) {

}