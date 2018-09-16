package ObjectManagement

import (
	"encoding/json"
	"github.com/mitchellh/mapstructure"
	"io/ioutil"
	"net/http"
	. "simplestorage/Misc"
	. "simplestorage/Mongo"
	. "simplestorage/Structure"
	"strconv"
)

func CreateTicket(w http.ResponseWriter, r *http.Request) {
	bucketName := GetBucketName(r)
	objectName := GetObjectName(r)

	buckExist := CheckBucketExist(bucketName)
	objExist := CheckObjectExist(objectName)

	if buckExist && !objExist {
		MakeObjectDirectory(bucketName, objectName)
		//@TODO update DB
		var object Object

		var temp map[string]interface{}
		temp = make(map[string]interface{})

		temp["name"] = objectName
		temp["bucket"] = bucketName
		temp["completed"] = false
		temp["created"] = GetTime()
		temp["modified"] = GetTime()

		mapstructure.Decode(temp, &object)

		CreateObject(object)

		w.WriteHeader(http.StatusOK)
	} else {
		w.WriteHeader(http.StatusBadRequest)
	}
}

func UploadPart(w http.ResponseWriter, r *http.Request) {
	bucketName := GetBucketName(r)
	objectName := GetObjectName(r)

	partNumber := r.URL.Query().Get("partNumber")
	valid := ValidatePattern(partNumber, PartNumPattern)

	tlength := r.Header.Get("Content-Length")
	length, lengthErr := strconv.Atoi(tlength)
	md5 := r.Header.Get("Content-MD5")

	ret := map[string]interface{}{
		"md5":        md5,
		"length":     length,
		"partNumber": partNumber,
	}

	/* VALIDATE REQUEST */
	if !valid {
		ret["error"] = ERROR["InvalidPartNumber"]
	}

	if lengthErr != nil {
		ret["error"] = ERROR["LengthMismatched"]
	}

	if md5 == "" {
		ret["error"] = ERROR["MD5Mismatched"]
	}

	if !FindOjbect(bucketName, objectName) {
		ret["error"] = ERROR["InvalidBucket"]
	}


	/* PERFORM REQUEST */
	var part Part
	part.Number = partNumber
	part.MD5 = md5
	part.Size = length
	part.Object = objectName

	b := r.Body
	f, _ := ioutil.ReadAll(b)
	defer b.Close()
	checksum, err := WriteFile(f,partNumber,objectName,bucketName)

	if checksum != md5 || err != nil || md5 == ""{
		ret["error"] = ERROR["MD5Mismatched"]
	}
	//log.Printf("checksum %s",checksum)
	if ret["error"] != nil {

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ret)

	}else{

		CreatePart(part)
		UpdateObjectLength(objectName,length)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(ret)

	}
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
