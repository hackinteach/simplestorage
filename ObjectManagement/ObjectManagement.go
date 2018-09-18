package ObjectManagement

import (
	"encoding/json"
	"github.com/mitchellh/mapstructure"
	"io/ioutil"
	"net/http"
	"regexp"
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
		temp["meta"] = make(map[string]interface{})

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
		UpdateObjectPart(objectName,partNumber)
		o,_ := GetObject(objectName,bucketName)
		o = o.UpdatePart(partNumber)
		UpdateObject(o)
		CreatePart(part)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(ret)

	}
}

func CompleteUpload(w http.ResponseWriter, r *http.Request) {
	bucketName := GetBucketName(r)
	objectName := GetObjectName(r)
	o, _ := GetObject(objectName, bucketName)

	//tl := r.Header.Get("Content-Length")
	//totalLength,_ := strconv.Atoi(tl)
	etag := r.Header.Get("Content-MD5")

	ret := map[string]interface{}{
		"name": objectName,
		"eTag" : etag,
	}

	if !CheckBucketExist(bucketName) {
		ret["error"] = 	Error.ErrorBucket
	}else if !FindObject(bucketName,objectName) {
		ret["error"] = Error.ErrorObjectName
	}else if o.Etag() != etag {
		ret["error"] = Error.ErrorMD5
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
	o, found := GetObject(objectName,bucketName)

	if !found {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

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
	// Update object in Object Collection
	UpdateObject(o)
	// Remove part from Part Collection
	RemovePart(partNumber,objectName)
	w.WriteHeader(http.StatusOK)
}

func DeleteObject(w http.ResponseWriter, r *http.Request) {
	bucketName := GetBucketName(r)
	objectName := GetObjectName(r)

	if FindObject(bucketName,objectName) {
		RemoveObjectDirectory(bucketName,objectName)
		w.WriteHeader(http.StatusOK)
	}

	w.WriteHeader(http.StatusBadRequest)
}

func DownloadObject(w http.ResponseWriter, r *http.Request) {
	bucketName := GetBucketName(r)
	objectName := GetObjectName(r)

	if !FindObject(bucketName,objectName){
		w.WriteHeader(http.StatusBadRequest)
	}

	hRange := r.Header.Get("Range")
	//eTag := r.Header.Get("ETag")

	o, found := GetObject(objectName,bucketName)

	if !found {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if hRange != "" {
		f := o.File()
		w.Write(f)
		w.WriteHeader(http.StatusOK)
	}else{
		r,_ := regexp.Compile("([0-9]+)")
		scope := r.FindAllString(hRange,-1)
		from,_ := strconv.Atoi(scope[0])
		var to int
		if len(scope) == 2 {
			to,_ = strconv.Atoi(scope[1])
		}else{
			to = o.Length()
		}
		f := o.FileRange(int64(from),int64(to))
		w.Write(f)
		w.WriteHeader(http.StatusOK)
	}
}

func UpdateMeta(w http.ResponseWriter, r *http.Request) {
	bucketName := GetBucketName(r)
	objectName := GetObjectName(r)
	key := r.URL.Query().Get("key")
	b := r.Body
	f, _ := ioutil.ReadAll(b)
	defer b.Close()
	value := string(f)
	o, found := GetObject(objectName, bucketName)
	if found {
		o.Meta[key] = value
		UpdateObject(o)
		w.WriteHeader(http.StatusOK)
		return
	}else{
		w.WriteHeader(http.StatusNotFound)
	}
}

func DeleteMeta(w http.ResponseWriter, r *http.Request) {
	bucketName := GetBucketName(r)
	objectName := GetObjectName(r)
	key := r.URL.Query().Get("key")
	o, found := GetObject(objectName, bucketName)
	if found {
		_, ok := o.Meta[key]
		if ok {
			delete(o.Meta,key)
			UpdateObject(o)
		}
		w.WriteHeader(http.StatusOK)
		return
	}else{
		w.WriteHeader(http.StatusNotFound)
	}
}

func GetMetaByKey(w http.ResponseWriter, r *http.Request) {
	bucketName := GetBucketName(r)
	objectName := GetObjectName(r)
	key := r.URL.Query().Get("key")

	o, found := GetObject(objectName, bucketName)
	if found {
		_,ok := o.Meta[key]
		w.Header().Set("Content-Type","application/json")
		w.WriteHeader(http.StatusOK)
		if ok {
			json.NewEncoder(w).Encode(o.Meta[key])
		}
		return
	}else{
		w.WriteHeader(http.StatusNotFound)
	}
}

func GetMeta(w http.ResponseWriter, r *http.Request) {
	bucketName := GetBucketName(r)
	objectName := GetObjectName(r)

	o, found := GetObject(objectName, bucketName)
	if found {
		w.Header().Set("Content-Type","application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(o.Meta)
		return
	}else{
		w.WriteHeader(http.StatusNotFound)
	}
}
