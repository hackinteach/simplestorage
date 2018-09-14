package Misc

import (
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"time"
)

const (
	BucketPath = "buckets/"
)

func ValidatePattern(str string)(bool){
	matched, err := regexp.MatchString("^[1-9][0-9]{0,3}0?",str)
	if err != nil{

	}
	return matched
}

func GetTime()(int64){
	return time.Now().UnixNano() / 1000000
}

func MakeDirectory(name string)(bool){
	var fullPath = filepath.Join(BucketPath,name)
	//log.Printf("path: %s",fullPath)
	err := os.MkdirAll(fullPath, 511)
	exec.Command("chmod","777",fullPath)
	return err == nil
}

func RemoveDirectory(name string)(bool){
	var fullPath = filepath.Join(BucketPath,name)

	err := os.Remove(fullPath)
	log.Print("Removing")
	log.Print(err)
	return err == nil
}

func GetBucketName(r *http.Request)(string) {
	return mux.Vars(r)["bucketName"]
}