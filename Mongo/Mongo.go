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
	PartCollection 	 *mgo.Collection
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
	PartCollection   = db.C("Part")
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

func CheckObjectExist(name string)(bool){
	var objects []Object
	ObjectCollection.Find(bson.M{"name":name}).All(&objects)

	for _,b := range objects{
		if b.Name == name {
			//log.Print("Found")
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

func GetBucket(name string)(TempBucket){
	var result TempBucket
	BucketCollection.Find(bson.M{"name":name}).One(&result)
	return result
}

func GetObjectList(bucketName string)([]TempObject){
	// @TODO Get object list from DB
	var result []TempObject
	ObjectCollection.Find(bson.M{"bucket":bucketName}).All(&result)
	return result
}

func CreateObject(object Object)(bool){
	return ObjectCollection.Insert(object) == nil
}

func FindOjbect(bucketName string, objectName string)(bool){
	var result []Object
	ObjectCollection.Find(bson.M{"bucket":bucketName}).All(&result)

	for _,o := range result {
		if o.Name == objectName {}
		return true
	}
	return false
}

func CreatePart(part Part)(bool){
	return PartCollection.Insert(part) == nil
}