package Mongo

import (
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"log"
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
			log.Print("Found")
			return true
		}
	}
	return false
}

func AddBucket(bucket Bucket) (bool) {
	err := BucketCollection.Insert(bucket)
	return err == nil
}

func RemoveBucket(name string)(bool){
	err := BucketCollection.Remove(bson.M{"name":name})
	return err == nil
}