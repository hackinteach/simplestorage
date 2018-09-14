package Misc

import (
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
	return err != nil
}