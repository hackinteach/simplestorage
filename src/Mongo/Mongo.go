package Mongo

import (
	"fmt"
	"gopkg.in/mgo.v2"
)

func Connect(
	info mgo.DialInfo){
		connectString := fmt.Sprintf("mongodb://%s:%s@%s:@%d/%s",user,password,host,port,db)
	session, err := mgo.Dial(connectString)
	if err != nil{

	}

}
