package ObjectManagement

import (
	"encoding/json"
	"fmt"
	"github.com/mitchellh/mapstructure"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"simplestorage/Error"
	. "simplestorage/Misc"
	. "simplestorage/Mongo"
	. "simplestorage/Structure"
	"sort"
	"strconv"
	"strings"
)

func CreateTicket(w http.ResponseWriter, r *http.Request) {
	bucketName := GetBucketName(r)
	objectName := GetObjectName(r)

	if !ValidatePattern(objectName,ObjNamePattern) {
		log.Print("Invalid Object Name")
		w.WriteHeader(http.StatusBadRequest)
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
		temp["byte"] = -1

		mapstructure.Decode(temp, &object)
		//log.Print(object)
		CreateObject(object)

		w.WriteHeader(http.StatusOK)
		return
	} else {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
}

func UploadPart(w http.ResponseWriter, r *http.Request) {

	bucketName := GetBucketName(r)
	objectName := GetObjectName(r)

	pn := r.URL.Query().Get("partNumber")
	partNumber,_ := strconv.Atoi(pn)
	valid := ValidatePattern(pn, PartNumPattern)

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
	} else if lengthErr != nil {
		ret["error"] = Error.ErrorLength
	} else if !FindObject(bucketName, objectName) {
		ret["error"] = Error.ErrorObjectName
	} else if !CheckBucketExist(bucketName) {
		ret["error"] = Error.ErrorBucket
	}

	path := filepath.Join(fmt.Sprintf("%s/%s/%s/%d",BucketPath,bucketName,objectName,partNumber))
	f, cr := os.Create(path)
	defer f.Close()
	if cr != nil {
		log.Print("Error creating file")
		log.Print(cr.Error())
	}
	b := r.Body;
	defer b.Close()
	n, wr := io.Copy(f,b)

	if wr != nil {
		log.Printf("Error writing file\n %s",wr.Error())
	}else{
		log.Printf("Written %d bytes",n)
	}

	checksum := Hash(path)

	if checksum != md5 || md5 == ""{
		ret["error"] = Error.ErrorMD5
		log.Print("MD5 error, removing file")
	}
	//log.Printf("checksum %s",checksum)
	if ret["error"] != nil {
		if err := os.Remove(path) ; err != nil {
			log.Printf("ERR REMOVE: %s",err.Error())
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ret)
		return
	}else{

		o,_ := GetObject(objectName,bucketName)
		part := Part{
			Number : partNumber,
			MD5: md5,
			Size: length,
			Object: objectName,
			Bucket: bucketName,
		}

		if !SearchIntArray(o.Part,partNumber){
			PushObjectPart(objectName,partNumber)
		}
		o = UpdatePart(o,partNumber)
		UpdateObject(o)

		/* PERFORM REQUEST */

		CreatePart(part)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(ret)
		return
	}
}

func CompleteUpload(w http.ResponseWriter, r *http.Request) {
	bucketName := GetBucketName(r)
	objectName := GetObjectName(r)
	o, _ := GetObject(objectName, bucketName)

	ret := map[string]interface{}{
		"name": objectName,
		"eTag" : Etag(o),
		"length": Length(o),
	}

	if bucketName == "" || !CheckBucketExist(bucketName) {
		ret["error"] = 	Error.ErrorBucket
	}else if !FindObject(bucketName,objectName) {
		ret["error"] = Error.ErrorObjectName
	}

	if ret["error"] != nil || ret["error"] != ""{
		log.Print("Cannot complete object")
		w.Header().Set("Content-Type","application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ret)
		return
	}else{
		log.Printf("Completed object : %s",objectName)
		o.Completed = true
		SetObjectComplete(o)
		w.Header().Set("Content-Type","application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(ret)
		return
	}
}

func DeletePart(w http.ResponseWriter, r *http.Request) {
	bucketName := GetBucketName(r)
	objectName := GetObjectName(r)
	pn := r.URL.Query().Get("partNumber")
	partNumber,_ := strconv.Atoi(pn)
	o, found := GetObject(objectName,bucketName)
	//log.Printf("DELETING: %s/%s",bucketName,objectName)
	if !found {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if 	o.Completed ||
		!FindObject(bucketName,objectName) ||
		!ValidatePattern(pn,PartNumPattern) ||
		!SearchIntArray(o.Part,partNumber) {
			//log.Printf("FAILED Partnumber: %s",pn)
			w.WriteHeader(http.StatusBadRequest)
			return
	}

	// Remove from dir
	err := RemovePartFile(bucketName,objectName,pn)
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
	return
}

func DeleteObject(w http.ResponseWriter, r *http.Request) {
	bucketName := GetBucketName(r)
	objectName := GetObjectName(r)

	if FindObject(bucketName,objectName) {
		o,_ := GetObject(objectName,bucketName)
		RemoveObjectDirectory(bucketName,objectName)
		RemoveObjectDB(o)
		w.WriteHeader(http.StatusOK)
		return
	}else{
		w.WriteHeader(http.StatusBadRequest)
		return
	}
}

func DownloadObject(w http.ResponseWriter, r *http.Request) {
	bucketName := GetBucketName(r)
	objectName := GetObjectName(r)

	if !FindObject(bucketName,objectName){
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	hRange := r.Header.Get("Range")
	//eTag := r.Header.Get("ETag")

	o, found := GetObject(objectName,bucketName)

	if !found {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if hRange == "" {
		sort.Ints(o.Part)
		for _, p := range o.Part {
			part := fmt.Sprintf("%d",p)
			path := filepath.Join(BucketPath,"/",bucketName,"/",objectName,"/",part)
			f,err := os.Open(path)
			if err != nil {
				log.Print(err.Error())
			}
			n, cop := io.Copy(w,f)
			if cop != nil {
				log.Print("Error copying file to response")
				log.Print(cop.Error())
			}else{
				log.Printf("Wrote %d to response",n)
			}
			f.Close()
		}

		w.WriteHeader(http.StatusOK)
		return
	}else{
		r,_ := regexp.Compile("([0-9]+)")
		scope := r.FindAllString(hRange,-1)
		from,_ := strconv.Atoi(scope[0])
		var to int
		if len(scope) == 2 {
			to,_ = strconv.Atoi(scope[1])
		}else{
			to = Length(o)
		}

		wrote := 0
		readSoFar := 0

		log.Printf("Range: %d - %d",from,to)
		for _,p:= range o.Part {
			path := filepath.Join(BucketPath,"/",bucketName,"/",objectName,"/",fmt.Sprintf("%d",p))
			f,_ := os.Open(path)
			part := FindPart(p,objectName,bucketName)
			//info, _ := f.Stat()
			//size := info.Size()
			log.Printf("======Part: %d======",p)
			pstart := readSoFar
			pend := readSoFar+part.Size

			log.Printf("Start: %d | End: %d",pstart,pend)

			if pend >= to && pstart <= from {
				actualStart := from-wrote-readSoFar
				//actualEnd := pend-wrote
				//log.Printf("ActualStart: %d",actualStart)
				f.Seek(int64(actualStart),0)
				n,err := io.CopyN(w,f,int64(to-from+1))
				if err != nil{
					log.Print(err.Error())
				}
				//log.Printf("Bytes wrote : %d",n)
				wrote += int(n)
			}else if pend >= from && pend <= to{
				//actualStart := from-pend // from end
				actualStart := from-pstart // from beginning
				//log.Printf("ActualStart: %d",actualStart)
				f.Seek(int64(actualStart),0)
				n,err := io.CopyN(w,f,int64(to-from+1))
				if err != nil{
					log.Print(err.Error())
				}
				//log.Printf("Bytes wrote : %d",n)
				wrote += int(n)
			}

			readSoFar += part.Size
			f.Close()
		}

		w.WriteHeader(http.StatusOK)
		return
	}
}

func UpdateMeta(w http.ResponseWriter, r *http.Request) {
	bucketName := GetBucketName(r)
	objectName := GetObjectName(r)
	//log.Printf("Meta : %s/%s",bucketName,objectName)
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
		return
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
		return
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
		return
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
		return
	}
}

func GetEndPoints(w http.ResponseWriter, r *http.Request) {
	if len(r.URL.RawQuery) > 0 {
		splitted := strings.Split(r.URL.RawQuery,"&")
		if len(splitted) == 1 {
			GetMeta(w,r)
		}else if len(splitted) == 2 {
			GetMetaByKey(w,r)
		}else {
			http.Error(w, "BadRouting", http.StatusNotFound)
			return
		}
	} else{
		DownloadObject(w,r)
	}
}
