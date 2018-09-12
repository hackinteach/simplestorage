package Mongo

import (
	. "../BucketManagement"
	"gopkg.in/mgo.v2"
)

var (
	session          *mgo.Session
	bucketCollection *mgo.Collection
	objectCollection *mgo.Collection
)

const (
	DB = "SimpleStorage"
)

func init() {
	var err error
	session, err = mgo.Dial("localhost")

	if err != nil {
		panic(err)
	}

	db := session.DB(DB)
	bucketCollection = db.C("Bucket")
	objectCollection = db.C("Object")
}

func AddBucket(bucket Bucket) (bool) {
	err := bucketCollection.Insert(bucket)
	return err != nil
}
