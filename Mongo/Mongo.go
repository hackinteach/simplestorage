package Mongo

import (
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	. "simplestorage/Structure"
)

var (
	session          *mgo.Session
	BucketCollection *mgo.Collection
	ObjectCollection *mgo.Collection
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
	BucketCollection = db.C("Bucket")
	ObjectCollection = db.C("Object")
}

func CheckBucketExist(bucketName string)(bool){
	var buckets []Bucket
	BucketCollection.Find(bson.M{"name":bucketName}).All(&buckets)

	for _,b := range buckets{
		if b.Name == bucketName {
			return true
		}
	}
	return false
}

func AddBucket(bucket Bucket) (bool) {
	err := BucketCollection.Insert(bucket)
	return err != nil
}

