package Misc

import (
	"os"
	"path/filepath"
	"regexp"
	"time"
)

const (
	BucketPath = "../buckets"
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
	var fullPath = filepath.Join("../data",name)

	err := os.MkdirAll(fullPath, 666)

	return err != nil
}