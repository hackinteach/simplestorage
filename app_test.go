package main

import (
	"errors"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"
)

// not used since we can set url query
func CreateURL(request Request) string{
	if len(request.Query) == 0 {
		return request.bucketName + request.objectName
	}
	url := request.bucketName + request.objectName + "?"
	for key := range request.Query {
		if request.Query[key] != "" {
			url += key + "=" + request.Query[key]
		}else{
			url += key
		}
		url += "&"
	}
	url = url[:len(url)-1]
	return url
}

func ExecuteRequest(request Request) error{
	//url := CreateURL(request) //create url is not used cus we can set the query by ourselves lmao
	//request.Body.Read

	req, err := http.NewRequest(request.Method, request.bucketName + request.objectName, nil)
	if err != nil {
		return err
	}

	if request.Body != nil {
		req.Body = ioutil.NopCloser(strings.NewReader(request.Body.String()))
	}

	if len(request.Headers) > 0{
		for key := range request.Headers  {
			req.Header.Set(key,request.Headers[key])
		}
	}
	qry := req.URL.Query()
	for key,value := range request.Query {
		qry.Add(key,value)
	}
	req.URL.RawQuery = qry.Encode()
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(request.Endpoint)
	handler.ServeHTTP(rr, req)
	status := rr.Code
	if status != request.ExpectedStatusCode  {
		headers := ""
		headers += req.URL.String() + "\n"
		headers += req.URL.Query().Encode() + "\n"
		headers += req.URL.String() + "\n"

		for key := range req.Header  {
			headers += key + " " + req.Header.Get("key") + "\n"
		}
		return errors.New(
			"Request Header:\n" + headers + "\n" +
				"Error:" + rr.Body.String()+"\n" +
				" Got:" + strconv.Itoa(status) + "\n"+
				" Want:" + strconv.Itoa(request.ExpectedStatusCode)+"\n" +
				"Method: " + req.Method)
	}
	return nil
}

func HandleError(t *testing.T,err error){
	if err != nil {
		t.Errorf(err.Error())
	}
}

func TestBucketEndpoints(t *testing.T){
	HandleError(t,ExecuteRequest(DeleteBucketRequest(http.StatusBadRequest,"/automated")))
	HandleError(t,ExecuteRequest(ListBucketRequest(http.StatusBadRequest,"/automated")))
	HandleError(t,ExecuteRequest(CreateBucketRequest(http.StatusOK,"/automated")))
	HandleError(t,ExecuteRequest(ListBucketRequest(http.StatusOK,"/automated")))
	HandleError(t,ExecuteRequest(DeleteBucketRequest(http.StatusOK,"/automated")))
	HandleError(t,ExecuteRequest(DeleteBucketRequest(http.StatusBadRequest,"/automated")))
	HandleError(t,ExecuteRequest(ListBucketRequest(http.StatusBadRequest,"/automated")))
}

func TestFlow(t *testing.T){
	filepathsmall1 := "./Test/testdata/application_test.csv-aa"
	filepathsmall2 := "./Test/testdata/application_test.csv-ab"
	filepathsmall3 := "./Test/testdata/application_test.csv-ac"
	filepathlarge1 := "./Test/testdata/application_test.csv--aa"
	filepathlarge2 := "./Test/testdata/application_test.csv--ab"
	filepathlarge3 := "./Test/testdata/application_test.csv--ac"
	filepathlarge4 := "./Test/testdata/application_test.csv--ad"
	filepathlarge5 := "./Test/testdata/application_test.csv--ae"
	filepathlarge6 := "./Test/testdata/application_test.csv--af"
	filepathlarge7 := "./Test/testdata/application_test.csv--ag"
	filepathvsmall1 := "./Test/testdata/simple.txt"
	bucket1 := CreateBucketRequest(http.StatusOK,"/automated")
	bucket2 := CreateBucketRequest(http.StatusOK,"/automated2")
	bucket3 := CreateBucketRequest(http.StatusOK,"/automated3")
	bucket4 := CreateBucketRequest(http.StatusOK,"/_automated-")
	bucket5 := CreateBucketRequest(http.StatusOK,"/_automated")
	bucket6 := CreateBucketRequest(http.StatusOK,"/automated_")
	bucket7 := CreateBucketRequest(http.StatusOK,"/a-utomated-")
	bucket8 := CreateBucketRequest(http.StatusOK,"/-au-tomated-")
	bucket9 := CreateBucketRequest(http.StatusOK,"/autom--a-t-ed")
	badbucket1 := CreateBucketRequest(http.StatusBadRequest,"/automated")
	badbucket2 := CreateBucketRequest(http.StatusBadRequest,"/automated2")
	badbucket3 := CreateBucketRequest(http.StatusBadRequest,"/automated3")
	badbucketname1 := CreateBucketRequest(http.StatusBadRequest,"/.automated")
	badbucketname2 := CreateBucketRequest(http.StatusBadRequest,"/automated2.")
	badbucketname3 := CreateBucketRequest(http.StatusBadRequest,"/a$$utomate$$d3")
	badbucketname4 := CreateBucketRequest(http.StatusBadRequest,"/automate&d3")
	badbucketname5 := CreateBucketRequest(http.StatusBadRequest,"/(automate*d3")
	badbucketname6 := CreateBucketRequest(http.StatusBadRequest,"/(automa....te*d3")
	deletebucket1 := DeleteBucketRequest(http.StatusOK,"/automated")
	deletebucket2 := DeleteBucketRequest(http.StatusOK,"/automated2")
	deletebucket3 := DeleteBucketRequest(http.StatusOK,"/automated3")
	deletebucket4 := DeleteBucketRequest(http.StatusOK,"/_automated-")
	deletebucket5 := DeleteBucketRequest(http.StatusOK,"/_automated")
	deletebucket6 := DeleteBucketRequest(http.StatusOK,"/automated_")
	deletebucket7 := DeleteBucketRequest(http.StatusOK,"/a-utomated-")
	deletebucket8 := DeleteBucketRequest(http.StatusOK,"/-au-tomated-")
	deletebucket9 := DeleteBucketRequest(http.StatusOK,"/autom--a-t-ed")
	baddeletebucket1 := DeleteBucketRequest(http.StatusBadRequest,"/automated")
	baddeletebucket2 := DeleteBucketRequest(http.StatusBadRequest,"/automated2")
	baddeletebucket3 := DeleteBucketRequest(http.StatusBadRequest,"/automated3")
	object11 := CreateObjectRequest(http.StatusOK, "/automated","/automatedsmall")
	object12 := CreateObjectRequest(http.StatusOK, "/automated","/automatedlarge")
	object13 := CreateObjectRequest(http.StatusOK, "/automated","/automatedvsmall")
	object14 := CreateObjectRequest(http.StatusOK, "/automated","/aut.omat--edsmall")
	object15 := CreateObjectRequest(http.StatusOK, "/automated","/au.tom_ated..large")
	object16 := CreateObjectRequest(http.StatusOK, "/automated","/auto...mat-edv.small")
	badobject11 := CreateObjectRequest(http.StatusBadRequest, "/automated","/automatedsmall")
	badobject12 := CreateObjectRequest(http.StatusBadRequest, "/automated","/automatedlarge")
	badobject13 := CreateObjectRequest(http.StatusBadRequest, "/automated","/automatedvsmall")
	badobjectname11 := CreateObjectRequest(http.StatusBadRequest, "/automated","/.au..tomated..small")
	badobjectname12 := CreateObjectRequest(http.StatusBadRequest, "/automated","/.automatedsmall.")
	badobjectname13 := CreateObjectRequest(http.StatusBadRequest, "/automated","/.automatedsmall")
	badobjectname14 := CreateObjectRequest(http.StatusBadRequest, "/automated","/automatedsma$ll$$")
	badobjectname15 := CreateObjectRequest(http.StatusBadRequest, "/automated","/aut@$omate$dsmall")
	badobjectname16 := CreateObjectRequest(http.StatusBadRequest, "/automated","/$aut$@#omate$dsmall")
	badobjectname17 := CreateObjectRequest(http.StatusBadRequest, "/automated","/$automateds$#$mall")
	badobjectname18 := CreateObjectRequest(http.StatusBadRequest, "/automated","/.d")
	badobjectname19 := CreateObjectRequest(http.StatusBadRequest, "/automated","/a.")
	deleteobject11 := DeleteObjectRequest(http.StatusOK,"/automated","/automatedsmall")
	deleteobject12 := DeleteObjectRequest(http.StatusOK,"/automated","/automatedlarge")
	deleteobject13 := DeleteObjectRequest(http.StatusOK,"/automated","/automatedvsmall")
	baddeleteobject11 := DeleteObjectRequest(http.StatusBadRequest,"/automated","/automatedsmall")
	baddeleteobject12 := DeleteObjectRequest(http.StatusBadRequest,"/automated","/automatedlarge")
	baddeleteobject13 := DeleteObjectRequest(http.StatusBadRequest,"/automated","/automatedvsmall")
	smallp1 := UploadPartsObjectRequest(http.StatusOK,"/automated","/automatedsmall",1,filepathsmall1)
	smallp2 := UploadPartsObjectRequest(http.StatusOK,"/automated","/automatedsmall",400,filepathsmall2)
	smallp3 := UploadPartsObjectRequest(http.StatusOK,"/automated","/automatedsmall",10000,filepathsmall3)
	largep1 := UploadPartsObjectRequest(http.StatusOK,"/automated","/automatedlarge",1,filepathlarge1)
	largep2 := UploadPartsObjectRequest(http.StatusOK,"/automated","/automatedlarge",2,filepathlarge2)
	largep3 := UploadPartsObjectRequest(http.StatusOK,"/automated","/automatedlarge",3,filepathlarge3)
	largep4 := UploadPartsObjectRequest(http.StatusOK,"/automated","/automatedlarge",4,filepathlarge4)
	largep5 := UploadPartsObjectRequest(http.StatusOK,"/automated","/automatedlarge",5,filepathlarge5)
	largep6 := UploadPartsObjectRequest(http.StatusOK,"/automated","/automatedlarge",6,filepathlarge6)
	largep7 := UploadPartsObjectRequest(http.StatusOK,"/automated","/automatedlarge",7,filepathlarge7)
	vsmallp1 := UploadPartsObjectRequest(http.StatusOK,"/automated","/automatedvsmall",4000,filepathvsmall1)
	badsmallp1 := UploadPartsObjectRequest(http.StatusBadRequest,"/automated","/automatedsmall",-1,filepathsmall1)
	badsmallp2 := UploadPartsObjectRequest(http.StatusBadRequest,"/automated","/automatedsmall",10001,filepathsmall2)
	badsmallp3 := UploadPartsObjectRequest(http.StatusBadRequest,"/automated","/automatedsmall",-1929,filepathsmall3)
	badpartnumsmallp1 := UploadStringPartsObjectRequest(http.StatusBadRequest,"/automated","/automatedsmall","dewoidewuiwdw",filepathsmall1)
	badpartnumsmallp2 := UploadStringPartsObjectRequest(http.StatusBadRequest,"/automated","/automatedsmall","3234r9vr",filepathsmall2)
	badpartnumsmallp3 := UploadStringPartsObjectRequest(http.StatusBadRequest,"/automated","/automatedsmall","nothisshitagain",filepathsmall3)
	deletesmallp1 := DeletePartObjectRequest(http.StatusOK,"/automated","/automatedsmall",1)
	deletesmallp2 := DeletePartObjectRequest(http.StatusOK,"/automated","/automatedsmall",400)
	deletesmallp3 := DeletePartObjectRequest(http.StatusOK,"/automated","/automatedsmall",10000)
	baddeletesmallp1 := DeletePartObjectRequest(http.StatusBadRequest,"/automated","/automatedsmall",-1)
	baddeletesmallp2 := DeletePartObjectRequest(http.StatusBadRequest,"/automated","/automatedsmall",10002)
	baddeletesmallp3 := DeletePartObjectRequest(http.StatusBadRequest,"/automated","/automatedsmall",-3443)
	badpartnumdeletesmallp1 := DeleteStringPartObjectRequest(http.StatusBadRequest,"/automated","/automatedsmall","dewoidewuiwdw")
	badpartnumdeletesmallp2 := DeleteStringPartObjectRequest(http.StatusBadRequest,"/automated","/automatedsmall","3204de")
	badpartnumdeletesmallp3 := DeleteStringPartObjectRequest(http.StatusBadRequest,"/automated","/automatedsmall","3.3-fr")
	completeSmall :=  CompleteUploadObjectRequest(http.StatusOK,"/automated","/automatedsmall")
	completeLarge :=  CompleteUploadObjectRequest(http.StatusOK,"/automated","/automatedlarge")
	completeVSmall :=  CompleteUploadObjectRequest(http.StatusOK,"/automated","/automatedvsmall")
	badCompleteSmall :=  CompleteUploadObjectRequest(http.StatusBadRequest,"/automated","/automatedsmall")
	badCompleteLarge :=  CompleteUploadObjectRequest(http.StatusBadRequest,"/automated","/automatedlarge")
	badCompleteVSmall :=  CompleteUploadObjectRequest(http.StatusBadRequest,"/automated","/automatedvsmall")
	metaDataSmall11 := PutMetadataRequest(http.StatusOK,"/automated","/automatedsmall","somekey1","somevalue1")
	metaDataSmall12 := PutMetadataRequest(http.StatusOK,"/automated","/automatedsmall","somekey1","somevalue2")
	metaDataSmall21 := PutMetadataRequest(http.StatusOK,"/automated","/automatedsmall","somekey2","somevalue1")
	metaDataSmall22 := PutMetadataRequest(http.StatusOK,"/automated","/automatedsmall","somekey2","somevalue2")
	badMetaDataSmall11 := PutMetadataRequest(http.StatusNotFound,"/automated","/automatedsmall","somekey1","somevalue1")
	badMetaDataSmall12 := PutMetadataRequest(http.StatusNotFound,"/automated","/automatedsmall","somekey1","somevalue2")
	badMetaDataSmall21 := PutMetadataRequest(http.StatusNotFound,"/automated","/automatedsmall","somekey2","somevalue1")
	badMetaDataSmall22 := PutMetadataRequest(http.StatusNotFound,"/automated","/automatedsmall","somekey2","somevalue2")
	deleteMetaDataSmall1 := DeleteMetadataKeyRequest(http.StatusOK,"/automated","/automatedsmall","somekey1")
	deleteMetaDataSmall2 := DeleteMetadataKeyRequest(http.StatusOK,"/automated","/automatedsmall","somekey2")
	badDeleteMetaDataSmall1 := DeleteMetadataKeyRequest(http.StatusNotFound,"/automated","/automatedsmall","somekey1")
	badDeleteMetaDataSmall2 := DeleteMetadataKeyRequest(http.StatusNotFound,"/automated","/automatedsmall","somekey2")
	getMetDataSmall1 :=  GetMetadataWithKeyRequest(http.StatusOK,"/automated","/automatedsmall","somekey1")
	getMetDataSmall2 :=  GetMetadataWithKeyRequest(http.StatusOK,"/automated","/automatedsmall","somekey2")
	badGetMetDataSmall1 :=  GetMetadataWithKeyRequest(http.StatusNotFound,"/automated","/automatedsmall","somekey1")
	badGetMetDataSmall2 :=  GetMetadataWithKeyRequest(http.StatusNotFound,"/automated","/automatedsmall","somekey2")
	getAllMetadataSmall := GetAllMetadataRequest(http.StatusOK,"/automated","/automatedsmall")
	getAllMetadataLarge := GetAllMetadataRequest(http.StatusOK,"/automated","/automatedlarge")
	getAllMetadataVSmall := GetAllMetadataRequest(http.StatusOK,"/automated","/automatedvsmall")
	badGetAllMetadataSmall := GetAllMetadataRequest(http.StatusNotFound,"/automated","/automatedsmall")
	badGetAllMetadataLarge := GetAllMetadataRequest(http.StatusNotFound,"/automated","/automatedlarge")
	badGetAllMetadataVSmall := GetAllMetadataRequest(http.StatusNotFound,"/automated","/automatedvsmall")
	badGetAllMetadataAnyhow1 := GetAllMetadataRequest(http.StatusNotFound,"/automated","/dowkfefe")
	badGetAllMetadataAnyhow2 := GetAllMetadataRequest(http.StatusNotFound,"/automated","/34r-rhgt")
	badGetAllMetadataAnyhow3 := GetAllMetadataRequest(http.StatusNotFound,"/automated","/dowkfefe")
	badGetAllMetadataAnyhow4 := GetAllMetadataRequest(http.StatusNotFound,"/ivrjfrfr","/autmomatedsmall")
	badGetAllMetadataAnyhow5 := GetAllMetadataRequest(http.StatusNotFound,"/wijfdirf-34rgt","/autmomatedsmall")


	HandleError(t,ExecuteRequest(badGetAllMetadataLarge))
	HandleError(t,ExecuteRequest(badGetAllMetadataSmall))
	HandleError(t,ExecuteRequest(badGetAllMetadataVSmall))
	HandleError(t,ExecuteRequest(badGetMetDataSmall1))
	HandleError(t,ExecuteRequest(badGetMetDataSmall2))
	HandleError(t,ExecuteRequest(badDeleteMetaDataSmall1))
	HandleError(t,ExecuteRequest(badDeleteMetaDataSmall2))
	HandleError(t,ExecuteRequest(baddeletebucket1))
	HandleError(t,ExecuteRequest(baddeletebucket2))
	HandleError(t,ExecuteRequest(baddeletebucket3))
	HandleError(t,ExecuteRequest(badobject11))
	HandleError(t,ExecuteRequest(badobject12))
	HandleError(t,ExecuteRequest(badobject13))
	HandleError(t,ExecuteRequest(bucket1))
	HandleError(t,ExecuteRequest(bucket2))
	HandleError(t,ExecuteRequest(bucket3))
	HandleError(t,ExecuteRequest(bucket4))
	HandleError(t,ExecuteRequest(bucket5))
	HandleError(t,ExecuteRequest(bucket6))
	HandleError(t,ExecuteRequest(bucket7))
	HandleError(t,ExecuteRequest(bucket8))
	HandleError(t,ExecuteRequest(bucket9))
	HandleError(t,ExecuteRequest(badbucket1))
	HandleError(t,ExecuteRequest(badbucket2))
	HandleError(t,ExecuteRequest(badbucket3))
	HandleError(t,ExecuteRequest(baddeleteobject11))
	HandleError(t,ExecuteRequest(baddeleteobject12))
	HandleError(t,ExecuteRequest(baddeleteobject13))
	HandleError(t,ExecuteRequest(object11))
	HandleError(t,ExecuteRequest(badobject11))
	HandleError(t,ExecuteRequest(object12))
	HandleError(t,ExecuteRequest(badobject12))
	HandleError(t,ExecuteRequest(object13))
	HandleError(t,ExecuteRequest(object14))
	HandleError(t,ExecuteRequest(object15))
	HandleError(t,ExecuteRequest(object16))
	HandleError(t,ExecuteRequest(badobject13))
	HandleError(t,ExecuteRequest(deleteobject11))
	HandleError(t,ExecuteRequest(deleteobject12))
	HandleError(t,ExecuteRequest(deleteobject13))
	HandleError(t,ExecuteRequest(deletebucket1))
	HandleError(t,ExecuteRequest(deletebucket2))
	HandleError(t,ExecuteRequest(deletebucket3))

	HandleError(t,ExecuteRequest(baddeletebucket1))
	HandleError(t,ExecuteRequest(baddeletebucket2))
	HandleError(t,ExecuteRequest(baddeletebucket3))

	HandleError(t,ExecuteRequest(badGetAllMetadataLarge))
	HandleError(t,ExecuteRequest(badGetAllMetadataSmall))
	HandleError(t,ExecuteRequest(badGetAllMetadataVSmall))
	HandleError(t,ExecuteRequest(badGetMetDataSmall1))
	HandleError(t,ExecuteRequest(badGetMetDataSmall2))
	HandleError(t,ExecuteRequest(badDeleteMetaDataSmall1))
	HandleError(t,ExecuteRequest(badDeleteMetaDataSmall2))
	HandleError(t,ExecuteRequest(badMetaDataSmall11))
	HandleError(t,ExecuteRequest(badMetaDataSmall12))
	HandleError(t,ExecuteRequest(badMetaDataSmall21))
	HandleError(t,ExecuteRequest(badMetaDataSmall22))
	HandleError(t,ExecuteRequest(badCompleteLarge))
	HandleError(t,ExecuteRequest(badCompleteSmall))
	HandleError(t,ExecuteRequest(badCompleteVSmall))
	HandleError(t,ExecuteRequest(bucket1))
	HandleError(t,ExecuteRequest(bucket2))
	HandleError(t,ExecuteRequest(bucket3))
	HandleError(t,ExecuteRequest(badMetaDataSmall11))
	HandleError(t,ExecuteRequest(badMetaDataSmall12))
	HandleError(t,ExecuteRequest(badMetaDataSmall21))
	HandleError(t,ExecuteRequest(badMetaDataSmall22))
	HandleError(t,ExecuteRequest(object11))
	HandleError(t,ExecuteRequest(object12))
	HandleError(t,ExecuteRequest(object13))
	HandleError(t,ExecuteRequest(badCompleteLarge))
	HandleError(t,ExecuteRequest(badCompleteSmall))
	HandleError(t,ExecuteRequest(badCompleteVSmall))
	HandleError(t,ExecuteRequest(smallp2))
	HandleError(t,ExecuteRequest(smallp1))
	HandleError(t,ExecuteRequest(smallp3))
	HandleError(t,ExecuteRequest(smallp1))
	HandleError(t,ExecuteRequest(smallp2))
	HandleError(t,ExecuteRequest(largep1))
	HandleError(t,ExecuteRequest(largep2))
	HandleError(t,ExecuteRequest(largep3))
	HandleError(t,ExecuteRequest(largep5))
	HandleError(t,ExecuteRequest(largep4))
	HandleError(t,ExecuteRequest(largep1))
	HandleError(t,ExecuteRequest(largep7))
	HandleError(t,ExecuteRequest(largep6))
	HandleError(t,ExecuteRequest(vsmallp1))
	HandleError(t,ExecuteRequest(badsmallp1))
	HandleError(t,ExecuteRequest(badsmallp2))
	HandleError(t,ExecuteRequest(badsmallp3))
	HandleError(t,ExecuteRequest(badpartnumsmallp1))
	HandleError(t,ExecuteRequest(badpartnumsmallp2))
	HandleError(t,ExecuteRequest(badpartnumsmallp3))
	HandleError(t,ExecuteRequest(baddeletesmallp1))
	HandleError(t,ExecuteRequest(baddeletesmallp2))
	HandleError(t,ExecuteRequest(baddeletesmallp3))
	HandleError(t,ExecuteRequest(badpartnumdeletesmallp1))
	HandleError(t,ExecuteRequest(badpartnumdeletesmallp2))
	HandleError(t,ExecuteRequest(badpartnumdeletesmallp3))
	HandleError(t,ExecuteRequest(deletesmallp1))
	HandleError(t,ExecuteRequest(deletesmallp2))
	HandleError(t,ExecuteRequest(deletesmallp3))
	HandleError(t,ExecuteRequest(smallp2))
	HandleError(t,ExecuteRequest(smallp1))
	HandleError(t,ExecuteRequest(smallp3))
	HandleError(t,ExecuteRequest(smallp1))
	HandleError(t,ExecuteRequest(smallp2))
	HandleError(t,ExecuteRequest(getAllMetadataSmall))
	HandleError(t,ExecuteRequest(getAllMetadataLarge))
	HandleError(t,ExecuteRequest(getAllMetadataVSmall))

	HandleError(t,ExecuteRequest(metaDataSmall11))
	HandleError(t,ExecuteRequest(metaDataSmall12))
	HandleError(t,ExecuteRequest(metaDataSmall21))
	HandleError(t,ExecuteRequest(metaDataSmall22))


	HandleError(t,ExecuteRequest(completeLarge))
	HandleError(t,ExecuteRequest(completeSmall))
	HandleError(t,ExecuteRequest(completeVSmall))
	HandleError(t,ExecuteRequest(badCompleteLarge))
	HandleError(t,ExecuteRequest(badCompleteSmall))
	HandleError(t,ExecuteRequest(badCompleteVSmall))

	HandleError(t,ExecuteRequest(metaDataSmall11))
	HandleError(t,ExecuteRequest(metaDataSmall12))
	HandleError(t,ExecuteRequest(metaDataSmall21))
	HandleError(t,ExecuteRequest(metaDataSmall22))

	HandleError(t,ExecuteRequest(getMetDataSmall1))
	HandleError(t,ExecuteRequest(getMetDataSmall2))

	HandleError(t,ExecuteRequest(deleteMetaDataSmall1))
	HandleError(t,ExecuteRequest(deleteMetaDataSmall2))

	HandleError(t,ExecuteRequest(getMetDataSmall1))
	HandleError(t,ExecuteRequest(getMetDataSmall2))

	HandleError(t,ExecuteRequest(metaDataSmall11))
	HandleError(t,ExecuteRequest(metaDataSmall12))
	HandleError(t,ExecuteRequest(metaDataSmall21))
	HandleError(t,ExecuteRequest(metaDataSmall22))

	HandleError(t,ExecuteRequest(getAllMetadataSmall))
	HandleError(t,ExecuteRequest(getAllMetadataLarge))
	HandleError(t,ExecuteRequest(getAllMetadataVSmall))
	HandleError(t,ExecuteRequest(badGetAllMetadataAnyhow1))
	HandleError(t,ExecuteRequest(badGetAllMetadataAnyhow2))
	HandleError(t,ExecuteRequest(badGetAllMetadataAnyhow3))
	HandleError(t,ExecuteRequest(badGetAllMetadataAnyhow4))
	HandleError(t,ExecuteRequest(badGetAllMetadataAnyhow5))


	HandleError(t,ExecuteRequest(deletebucket1))
	HandleError(t,ExecuteRequest(deletebucket2))
	HandleError(t,ExecuteRequest(deletebucket3))
	HandleError(t,ExecuteRequest(deletebucket4))
	HandleError(t,ExecuteRequest(deletebucket5))
	HandleError(t,ExecuteRequest(deletebucket6))
	HandleError(t,ExecuteRequest(deletebucket7))
	HandleError(t,ExecuteRequest(deletebucket8))
	HandleError(t,ExecuteRequest(deletebucket9))
	HandleError(t,ExecuteRequest(badbucketname1))
	HandleError(t,ExecuteRequest(badbucketname2))
	HandleError(t,ExecuteRequest(badbucketname3))
	HandleError(t,ExecuteRequest(badbucketname4))
	HandleError(t,ExecuteRequest(badbucketname5))
	HandleError(t,ExecuteRequest(badbucketname6))
	HandleError(t,ExecuteRequest(badobjectname11))
	HandleError(t,ExecuteRequest(badobjectname12))
	HandleError(t,ExecuteRequest(badobjectname13))
	HandleError(t,ExecuteRequest(badobjectname14))
	HandleError(t,ExecuteRequest(badobjectname15))
	HandleError(t,ExecuteRequest(badobjectname16))
	HandleError(t,ExecuteRequest(badobjectname17))
	HandleError(t,ExecuteRequest(badobjectname18))
	HandleError(t,ExecuteRequest(badobjectname19))
}
