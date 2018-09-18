package Misc

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
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
	"simplestorage/Mongo"
	"simplestorage/Structure"
	"sort"
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

func RemoveBucketDirectory(name string) (bool) {
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
	checksum := fmt.Sprintf("%x", hash.Sum([]byte{}))
	return checksum, nil
}

func RemovePartFile(bucket string, object string, part string) (Error error) {
	path := fmt.Sprintf("%s/%s/%s/%s", BucketPath, bucket, object, part)
	return os.Remove(path)
}

func SearchStringArray(arr []string, search string) (bool) {
	for _, st := range arr {
		if st == search {
			return true
		}
	}
	return false
}

func RemoveItem(arr []string, item string) (Result []string) {
	for i, elm := range arr {
		if elm == item {
			return remove(arr, i)
		}
	}
	return arr
}

func remove(s []string, i int) []string {
	s[len(s)-1], s[i] = s[i], s[len(s)-1]
	return s[:len(s)-1]
}

func RemoveObjectDirectory(bucket string, object string) {
	path := filepath.Join(BucketPath, "/", bucket, "/", object)
	os.RemoveAll(path)
}

func GetFile(path string)[]byte{
	f,_ := os.Open(path)
	var ret []byte
	f.Read(ret)
	return ret
}


func Etag(o Structure.Object) string {
	hasher := md5.New()
	parts := Mongo.FindParts(o.Name)
	var md5 []string

	for _,p := range parts {
		md5 = append(md5, p.MD5)
	}

	md5s := strings.Join(md5,"")
	hasher.Write([]byte(md5s))
	hashed := hex.EncodeToString(hasher.Sum(nil))
	return fmt.Sprintf("%s-%d",hashed,len(md5))
}

func Length(o Structure.Object) int {
	length := 0
	for _,v := range o.Part {
		p := Mongo.FindPart(v)
		length += p.Size
	}
	return length
}

func UpdatePart(o Structure.Object, part string) Structure.Object {
	o.Part = append(o.Part,part)
	sort.Strings(o.Part)
	return o
}

func File(o Structure.Object) []byte {
	var ret []byte
	bucket := o.Bucket
	parts := o.Part
	name := o.Name
	for _,p := range parts {
		path := filepath.Join(BucketPath,"/",bucket,"/",name,"/",p)
		f := GetFile(path)
		for _,v := range f {
			ret = append(ret,v)
		}
	}
	return ret
}

func FileRange(o Structure.Object, from, to int64) []byte{
	var res []byte
	bucket := o.Bucket
	parts := o.Part
	name := o.Name
	count := int64(0)
	for _,p := range parts {
		path := filepath.Join(BucketPath, "/", bucket, "/", name, "/", p)
		f,_ := os.Open(path)
		info,_ := f.Stat()
		size := info.Size()
		count += size
		if count >= from {
			if count <= to{
				f.ReadAt(res,from-count)
			}else{
				var tmp []byte
				f.Read(tmp)
				tmp = tmp[:to-count]
				for _,v := range tmp{
					res = append(res,v)
				}
			}
		}
	}
	return res
}