package BucketManagement

import (
	. "../Misc"
	. "../Mongo"
	. "../Structure"
	"github.com/mitchellh/mapstructure"
	"net/http"
	"strings"
)

func CreateBucket(w http.ResponseWriter, r *http.Request) {

	var bucket Bucket
	var tmp map[string]interface{}

	tmp = make(map[string]interface{})

	bucketName := strings.Replace(r.URL.Path,"/","",1)
	tmp["Name"] = bucketName
	tmp["Created"] = GetTime()
	tmp["Modified"] = GetTime()
	mapstructure.Decode(tmp,&bucket)
	if ! CheckBucketExist(bucketName){
		AddBucket(bucket)
	}
}


func DeleteBucket(w http.ResponseWriter, r *http.Request) {

}

func ListBucket(w http.ResponseWriter, r *http.Request) {

}