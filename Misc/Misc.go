package Misc

import (
	"bytes"
	"crypto/md5"
	"fmt"
	"github.com/gorilla/mux"
	"io"
	"io/ioutil"
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

func ValidatePattern(str string, pat string) (bool) {
	matched, _ := regexp.MatchString(pat, str)
	return matched
}

func GetTime() (int64) {
	return time.Now().UnixNano() / 1000000
}

func MakeBucketDirectory(name string) (bool) {
	var fullPath = filepath.Join(BucketPath, name)
	//log.Printf("path: %s",fullPath)
	err := os.MkdirAll(fullPath, 0777)
	//exec.Command("chmod", "777", fullPath)
	return err == nil
}

func MakeObjectDirectory(bucket string, name string) (bool) {
	var fullPath = filepath.Join(BucketPath, bucket, "/", name)
	err := os.Mkdir(fullPath, 511)
	exec.Command("chmod", "777", fullPath)
	return err == nil
}

func RemoveDirectory(name string) (bool) {
	var fullPath = filepath.Join(BucketPath, name)

	err := os.RemoveAll(fullPath)
	return err == nil
}

func GetBucketName(r *http.Request) (string) {
	return strings.ToLower(mux.Vars(r)["bucketName"])
}

func GetObjectName(r *http.Request) (string) {
	return strings.ToLower(mux.Vars(r)["objectName"])
}

/**
Write file to /{bucketName}/{objectName}/{filename}
md5 along the fly and return md5 checksum
 */
func WriteFile(f []byte, filename string, object string, bucket string) (MD5 string, Error error) {
	hash := md5.New()

	permission := os.FileMode(0755).Perm()
	path := filepath.Join(BucketPath, "/", bucket, "/", object, "/", filename)
	//file, _ := ioutil.ReadAll(f)

	err := ioutil.WriteFile(path, f, permission)
	if err != nil {
		log.Print(err)
		return "", err
	}
	reader := bytes.NewReader(f)
	if _, err := io.Copy(hash, reader); err != nil {
		return "", err
	}
	checksum := fmt.Sprintf("%x",hash.Sum([]byte{}))
	return checksum, nil
}

func RemovePartFile(bucket string, object string, part string)(Error error){
	path := fmt.Sprintf("%s/%s/%s/%s",BucketPath,bucket,object,part)
	return os.Remove(path)
}

func SearchStringArray(arr []string, search string)(bool){
	for _,st := range arr {
		if st == search {
			return true
		}
	}
	return false
}

func RemoveItem(arr []string, item string)(Result []string){
	for i, elm := range arr{
		if elm == item {
			return remove(arr,i)
		}
	}
	return arr
}

func remove(s []string, i int) []string {
	s[len(s)-1], s[i] = s[i], s[len(s)-1]
	return s[:len(s)-1]
}