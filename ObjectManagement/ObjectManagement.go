package ObjectManagement

import (
	"encoding/json"
	"github.com/mitchellh/mapstructure"
	"net/http"
	. "simplestorage/Misc"
	. "simplestorage/Mongo"
	. "simplestorage/Structure"
)

func CreateTicket(w http.ResponseWriter, r *http.Request) {
	bucketName := GetBucketName(r)
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

func UploadPart(w http.ResponseWriter, r *http.Request) {
	bucketName := GetBucketName(r)
	objectName := GetObjectName(r)

	partNumber := r.URL.Query().Get("partNumber")
	valid := ValidatePattern(partNumber,PART_NUM_PATTERN)

	length := r.Header.Get("Content-Length")
	md5 := r.Header.Get("Content-MD5")

	var ret = map[string]string{
		"md5": md5,
		"length": length,
		"partNumber": partNumber,
	}

	w.Header().Set("Content-Type","application/json")

	/* VALIDATE REQUEST */

	if !valid{
		ret["error"] = ERROR["InvalidPartNumber"]
	}

	if length == "" {
		ret["error"] = ERROR["LengthMismatched"]
	}

	if md5 == "" {
		ret["error"] = ERROR["MD5Mismatched"]
	}

	if !FindOjbect(bucketName,objectName) {
		ret["error"] = ERROR["InvalidBucket"]
	}

	if ret["error"] == "" {
		json.NewEncoder(w).Encode(ret)
		w.WriteHeader(400)
	}

	/* PERFORM REQUEST */
	//@TODO Make Part struct and create Object, add to db
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