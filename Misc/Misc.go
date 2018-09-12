package Misc

import (
	"regexp"
	"time"
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