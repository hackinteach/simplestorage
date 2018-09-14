package ObjectManagement

import (
	. "../Misc"
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
)

func CreateTicket(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	bucketName := vars["bucketName"]
	json.NewEncoder(w).Encode(bucketName)
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