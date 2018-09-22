package main

import (
	"bytes"
	"io"
	"log"
	"net/http"
	"os"
	"simplestorage/BucketManagement"
	"simplestorage/Misc"
	"simplestorage/ObjectManagement"
	"strconv"
)

type Request struct {
	bucketName  string
	objectName string
	ExpectedStatusCode int
	Method string
	Query map[string]string
	Endpoint func(w http.ResponseWriter, r *http.Request)
	Headers map[string]string
	Body *bytes.Buffer
}

func CreateBucketRequest(expectedstatuscode int,bucketname string) Request {
	return Request{
		bucketName: bucketname,
		objectName: "",
		ExpectedStatusCode: expectedstatuscode,
		Method: "POST",
		Query: map[string]string{"create": ""},
		Endpoint: BucketManagement.CreateBucket,
		Headers: map[string]string{},
		Body: nil,
	}
}

func DeleteBucketRequest(expectedstatuscode int,bucketname string) Request {
	return Request{
		bucketName: bucketname,
		objectName: "",
		ExpectedStatusCode: expectedstatuscode,
		Method: "DELETE",
		Query: map[string]string{"delete": ""},
		Endpoint: BucketManagement.DeleteBucket,
		Headers: map[string]string{},
		Body: nil,
	}
}

func ListBucketRequest(expectedstatuscode int,bucketname string) Request {
	return Request{
		bucketName: bucketname,
		objectName: "",
		ExpectedStatusCode: expectedstatuscode,
		Method: "GET",
		Query: map[string]string{"list": ""},
		Endpoint: BucketManagement.ListBucket,
		Headers: map[string]string{},
		Body: nil,
	}
}

func CreateObjectRequest(expectedstatuscode int,bucketname string,objectname string) Request {
	return Request{
		bucketName: bucketname,
		objectName: objectname,
		ExpectedStatusCode: expectedstatuscode,
		Method: "POST",
		Query: map[string]string{"crate": ""},
		Endpoint: ObjectManagement.CreateTicket,
		Headers: map[string]string{},
		Body: nil,
	}
}

func UploadPartsObjectRequest(expectedstatuscode int,bucketname string,objectname string,partnumber int,filepath string) Request {
	md5:= Misc.Hash(filepath)

	file, err := os.Open(filepath)
	if err != nil {
		log.Fatalln(err.Error())
	}
	stat, err := os.Stat(filepath)
	if err != nil {
		log.Fatal(err.Error())
	}
	filesize  := strconv.Itoa(int(stat.Size()))
	body := &bytes.Buffer{}
	io.Copy(body, file)
	file.Close()
	return Request{
		bucketName: bucketname,
		objectName: objectname,
		ExpectedStatusCode: expectedstatuscode,
		Method: "PUT",
		Query: map[string]string{"partNumber": strconv.Itoa(partnumber)},
		Endpoint: ObjectManagement.UploadPart,
		Headers: map[string]string{"Content-MD5":md5,"Content-Length": filesize},
		Body: body,
	}
}

func UploadStringPartsObjectRequest(expectedstatuscode int,bucketname string,objectname string,partnumber string,filepath string) Request {
	md5 := Misc.Hash(filepath)

	file, err := os.Open(filepath)
	if err != nil {
		log.Fatalln(err.Error())
	}
	stat, err := os.Stat(filepath)
	if err != nil {
		log.Fatal(err.Error())
	}
	filesize  := strconv.Itoa(int(stat.Size()))
	body := &bytes.Buffer{}
	//writer := multipart.NewWriter(body)
	//part, err := writer.CreateFormFile("binary",filepath)
	//if err != nil {
	//	log.Fatal(err)
	//}
	io.Copy(body, file)
	file.Close()
	return Request{
		bucketName: bucketname,
		objectName: objectname,
		ExpectedStatusCode: expectedstatuscode,
		Method: "PUT",
		Query: map[string]string{"partNumber": partnumber},
		Endpoint: ObjectManagement.UploadPart,
		Headers: map[string]string{"Content-MD5":md5,"Content-Length": filesize},
		Body: body,
	}
}

func CompleteUploadObjectRequest(expectedstatuscode int,bucketname string,objectname string) Request {
	return Request{
		bucketName: bucketname,
		objectName: objectname,
		ExpectedStatusCode: expectedstatuscode,
		Method: "POST",
		Query: map[string]string{"complete":""},
		Endpoint: ObjectManagement.CompleteUpload,
		Headers: map[string]string{},
		Body: nil,
	}
}

func DeletePartObjectRequest(expectedstatuscode int,bucketname string,objectname string,partnumber int) Request {
	return Request{
		bucketName: bucketname,
		objectName: objectname,
		ExpectedStatusCode: expectedstatuscode,
		Method: "DELETE",
		Query: map[string]string{"partNumber":strconv.Itoa(partnumber)},
		Endpoint: ObjectManagement.CompleteUpload,
		Headers: map[string]string{},
		Body: nil,
	}
}

func DeleteStringPartObjectRequest(expectedstatuscode int,bucketname string,objectname string,partnumber string) Request {
	return Request{
		bucketName: bucketname,
		objectName: objectname,
		ExpectedStatusCode: expectedstatuscode,
		Method: "DELETE",
		Query: map[string]string{"partNumber":partnumber},
		Endpoint: ObjectManagement.DeletePart,
		Headers: map[string]string{},
		Body: nil,
	}
}

func DeleteObjectRequest(expectedstatuscode int,bucketname string,objectname string,) Request {
	return Request{
		bucketName: bucketname,
		objectName: objectname,
		ExpectedStatusCode: expectedstatuscode,
		Method: "POST",
		Query: map[string]string{},
		Endpoint: ObjectManagement.DeleteObject,
		Headers: map[string]string{},
		Body: nil,
	}
}


func DownloadObjectRequest(expectedstatuscode int,bucketname string,objectname string,ranges map[int][]int) Request {
	headers := map[string]string{}
	if len(ranges) > 0 {
		headers["Ranges"] = "bytes="
		for key := range ranges{
			if len(ranges[key]) == 1 {
				headers["Ranges"] += strconv.Itoa(ranges[key][0])+"-"
			}else if len(ranges[key]) == 2 {
				headers["Ranges"] += strconv.Itoa(ranges[key][0])+"-"+strconv.Itoa(ranges[key][1])
			}
			headers["Ranges"] += ","
		}
		headers["Ranges"] = headers["Ranges"][:len(headers["Ranges"])-1]
	}
	return Request{
		bucketName: bucketname,
		objectName: objectname,
		ExpectedStatusCode: expectedstatuscode,
		Method: "GET",
		Query: map[string]string{},
		Endpoint: ObjectManagement.DownloadObject,
		Headers: headers,
		Body: nil,
	}
}

func PutMetadataRequest(expectedstatuscode int,bucketname string,objectname string,key string,value string) Request {
	body := &bytes.Buffer{}
	body.Write([]byte(value))
	return Request{
		bucketName: bucketname,
		objectName: objectname,
		ExpectedStatusCode: expectedstatuscode,
		Method: "PUT",
		Query: map[string]string {"metadata":"","key":key},
		Endpoint: ObjectManagement.UpdateMeta,
		Headers: map[string]string{},
		Body: body,
	}
}

func DeleteMetadataKeyRequest(expectedstatuscode int,bucketname string,objectname string,key string) Request {
	return Request{
		bucketName:         bucketname,
		objectName:         objectname,
		ExpectedStatusCode: expectedstatuscode,
		Method:             "DELETE",
		Query:              map[string]string {"metadata":"","key":key},
		Endpoint:           ObjectManagement.DeleteMeta,
		Headers: map[string]string{},
		Body: nil,
	}
}

func GetMetadataWithKeyRequest(expectedstatuscode int,bucketname string,objectname string,key string) Request {
	return Request{
		bucketName:         bucketname,
		objectName:         objectname,
		ExpectedStatusCode: expectedstatuscode,
		Method:             "GET",
		Query:              map[string]string {"metadata":"","key":key},
		Endpoint:           ObjectManagement.GetMetaByKey,
		Headers: map[string]string{},
		Body: nil,
	}
}

func GetAllMetadataRequest(expectedstatuscode int,bucketname string,objectname string) Request {
	return Request{
		bucketName: bucketname,
		objectName: objectname,
		ExpectedStatusCode: expectedstatuscode,
		Method: "GET",
		Query: map[string]string{"metadata":""},
		Endpoint: ObjectManagement.GetMeta,
		Headers: map[string]string{},
		Body: nil,
	}
}