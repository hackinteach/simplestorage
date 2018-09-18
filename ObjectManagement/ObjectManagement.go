package ObjectManagement

import (
	"encoding/json"
	"github.com/mitchellh/mapstructure"
	"io/ioutil"
	"net/http"
	"simplestorage/Error"
	. "simplestorage/Misc"
	. "simplestorage/Mongo"
	. "simplestorage/Structure"
	"strconv"
)

func CreateTicket(w http.ResponseWriter, r *http.Request) {
	bucketName := GetBucketName(r)
	objectName := GetObjectName(r)

	if !ValidatePattern(objectName,ObjNamePattern) {
		w.WriteHeader(200)
		return
	}

	buckExist := CheckBucketExist(bucketName)
	objExist := FindObject(bucketName,objectName)

	if buckExist && !objExist {
		MakeObjectDirectory(bucketName, objectName)
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
		ret["error"] = Error.ErrorPartNumber
	}

	if lengthErr != nil {
		ret["error"] = Error.ErrorLength
	}

	if md5 == "" {
		ret["error"] = Error.ErrorMD5
	}

	if !FindObject(bucketName, objectName) {
		ret["error"] = Error.ErrorObjectName
	}

	if !CheckBucketExist(bucketName) {
		ret["error"] = Error.ErrorBucket
	}


	/* PERFORM REQUEST */
	UpdateObjectPart(objectName,partNumber)
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
		ret["error"] = Error.ErrorMD5
	}
	//log.Printf("checksum %s",checksum)
	if ret["error"] != nil {

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ret)

	}else{

		CreatePart(part)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(ret)

	}
}

func CompleteUpload(w http.ResponseWriter, r *http.Request) {
	bucketName := GetBucketName(r)
	objectName := GetObjectName(r)
	o := GetObject(objectName)

	tl := r.Header.Get("Content-Length")
	totalLength,_ := strconv.Atoi(tl)
	etag := r.Header.Get("Content-MD5")

	ret := map[string]interface{}{
		"name": objectName,
		"eTag" : etag,
		"length": totalLength,
	}

	if !CheckBucketExist(bucketName) {
		ret["error"] = 	Error.ErrorBucket
	}else if !FindObject(bucketName,objectName) {
		ret["error"] = Error.ErrorObjectName
	}else if o.Etag() != etag {
		ret["error"] = Error.ErrorMD5
	}else if o.Length() != totalLength {
		ret["error"] = Error.ErrorLength
	}

	if ret["error"] != nil {
		w.Header().Set("Content-Type","application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ret)
	}else{
		o.Completed = true
		SetObjectComplete(o.Name)
		w.Header().Set("Content-Type","application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(ret)
	}
}

func DeletePart(w http.ResponseWriter, r *http.Request) {
	bucketName := GetBucketName(r)
	objectName := GetObjectName(r)
	partNumber := r.URL.Query().Get("partNumber")
	o := GetObject(objectName)

	if 	o.Completed ||
		!FindObject(bucketName,objectName) ||
		partNumber == "" ||
		!SearchStringArray(o.Part,partNumber) {
			w.WriteHeader(http.StatusBadRequest)
			return
	}

	// Remove from dir
	err := RemovePartFile(bucketName,objectName,partNumber)
	if err != nil {
		w.Write([]byte("Cannot remove part, please try again"))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Remove from DB
	o.Part = RemoveItem(o.Part,partNumber)
	UpdateObject(o)
	RemovePart(partNumber,objectName)
	w.WriteHeader(http.StatusOK)
}

func DeleteObject(w http.ResponseWriter, r *http.Request) {
	bucketName := GetBucketName(r)
	objectName := GetObjectName(r)

	if FindObject(bucketName,objectName) {

	}
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
