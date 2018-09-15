package ObjectManagement

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/mitchellh/mapstructure"
	"log"
	"net/http"
	. "simplestorage/Misc"
	. "simplestorage/Mongo"
	. "simplestorage/Structure"
)

func CreateTicket(w http.ResponseWriter, r *http.Request) {
	bucketName := GetBucketName(r)
	log.Printf("Bucketname: %s",bucketName)
	objectName := GetObjectName(r)

	buckExist := CheckBucketExist(bucketName)
	objExist := CheckObjectExist(objectName)

	if buckExist && !objExist {
		MakeObjectDirectory(bucketName,objectName)
		//@TODO update DB
		var object Object

		var temp map[string]interface{}
		temp = make(map[string]interface{})

		temp["name"] = objectName
		temp["bucket"] = bucketName

		mapstructure.Decode(temp,&object)

		CreateObject(object)

		w.Header().Set("Content-Type","application/json")
		json.NewEncoder(w).Encode(object)
		w.WriteHeader(200)
	}else{
		w.WriteHeader(400)
	}
}

func UploadAll(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	partNumber := r.URL.Query().Get("partNumber")
	valid := ValidatePattern(partNumber)
	if !valid{

	}
	json.NewEncoder(w).Encode(vars)
}

func CompleteUpload(w http.ResponseWriter, r *http.Request) {

}

func DeletePart(w http.ResponseWriter, r *http.Request) {

}

func DeleteObject(w http.ResponseWriter, r *http.Request) {

}

func DownloadObject(w http.ResponseWriter, r *http.Request) {

}

func UpdateMeta(w http.ResponseWriter, r *http.Request) {

}

func DeleteMeta(w http.ResponseWriter, r *http.Request) {

}

func GetMetaByKey(w http.ResponseWriter, r *http.Request) {

}

func GetMeta(w http.ResponseWriter, r *http.Request) {

}