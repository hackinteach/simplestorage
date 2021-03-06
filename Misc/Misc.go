package Misc

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
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
	err := os.MkdirAll(fullPath, 0755)
	//exec.Command("chmod", "777", fullPath).Run()
	//log.Print("MKDIR")
	return err == nil
}

func MakeObjectDirectory(bucket string, name string) (bool) {
	var fullPath = filepath.Join(BucketPath, bucket, "/", name)
	err := os.MkdirAll(fullPath, 0755)
	log.Print("Created object dir")
	if err != nil {
		log.Print(err.Error())
	}
	return err == nil

}

func RemoveBucketDirectory(name string) (bool) {
	var fullPath = filepath.Join(BucketPath, name)

	err := os.RemoveAll(fullPath)
	return err == nil
}

func GetBucketName(r *http.Request) (string) {
	split := strings.Split(r.URL.Path,"/")
	if len(split) < 2 {
		log.Print("Bucket name not specified")
		return ""
	}else if ValidatePattern(split[1],Structure.BuckNamePattern){
		return split[1]
	}
	return ""
}

func GetObjectName(r *http.Request) (string) {
	split := strings.Split(r.URL.Path,"/")
	if len(split) < 2 {
		log.Print("Bucket name not specified")
		return ""
	}else if ValidatePattern(split[2],Structure.ObjNamePattern){
		return split[2]
	}
	return ""
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

func SearchIntArray(arr []int, search int) (bool) {
	for _, st := range arr {
		if st == search {
			return true
		}
	}
	return false
}

func RemoveItem(arr []int, item int) (Result []int) {
	for i, elm := range arr {
		if elm == item {
			return remove(arr, i)
		}
	}
	return arr
}

func remove(s []int, i int) []int {
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

	sort.SliceStable(parts, func (i,j int) bool {
		return parts[i].Number < parts[i].Number
	})

	for _,p := range parts {
		log.Printf("Part Number: %d %s",p.Number, p.MD5)
		md5 = append(md5, p.MD5)
	}
	j := strings.Join(md5,"")
	md5s,_ := hex.DecodeString(j)
	hasher.Write(md5s)
	hashed := hex.EncodeToString(hasher.Sum(nil)[:16])
	return fmt.Sprintf("%s-%d",hashed,len(md5))
}

func Length(o Structure.Object) int {
	length := 0
	for _,v := range o.Part {
		p := Mongo.FindPart(v, o.Name, o.Bucket)
		length += p.Size
	}
	return length
}

func UpdatePart(o Structure.Object, part int) Structure.Object {
	if !SearchIntArray(o.Part,part){
		o.Part = append(o.Part,part)
		sort.Ints(o.Part)
	}
	return o
}

func File(o Structure.Object) []byte {
	var ret []byte
	bucket := o.Bucket
	parts := o.Part
	name := o.Name
	for _,p := range parts {
		path := filepath.Join(BucketPath,"/",bucket,"/",name,"/",fmt.Sprintf("%d",p))
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
		path := filepath.Join(BucketPath, "/", bucket, "/", name, "/", fmt.Sprintf("%d",p))
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

func Hash(path string)string {
	f,err := os.Open(path)

	if err != nil {
		return ""
	}

	defer f.Close()
	hasher := md5.New()

	if _,err := io.Copy(hasher, f); err != nil {
		return ""
	}

	sum := hasher.Sum(nil)
	//log.Printf("%x",sum)
	return fmt.Sprintf("%x",sum)
}