package Misc

import (
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"time"
)

const (
	BucketPath = "buckets/"
)

func ValidatePattern(str string, pat string)(bool){
	matched, _ := regexp.MatchString(pat,str)
	return matched
}

func GetTime()(int64){
	return time.Now().UnixNano() / 1000000
}

func MakeBucketDirectory(name string)(bool){
	var fullPath = filepath.Join(BucketPath,name)
	//log.Printf("path: %s",fullPath)
	err := os.MkdirAll(fullPath, 511)
	exec.Command("chmod","777",fullPath)
	return err == nil
}

func MakeObjectDirectory(bucket string, name string)(bool){
	var fullPath = filepath.Join(BucketPath,bucket,"/",name)
	err := os.Mkdir(fullPath,511)
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
	return strings.ToLower(mux.Vars(r)["bucketName"])
}

func GetObjectName(r *http.Request)(string){
	return strings.ToLower(mux.Vars(r)["objectName"])
}