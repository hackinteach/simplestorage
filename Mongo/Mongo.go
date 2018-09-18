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

func AddBucket(bucket Bucket) (bool) {
	err := BucketCollection.Insert(bucket)
	return err == nil
}

func RemoveBucket(name string)(bool){
	err := BucketCollection.Remove(bson.M{"name":name})
	return err == nil
}

func GetReturnBucket(name string)(TempBucket){
	var result TempBucket
	BucketCollection.Find(bson.M{"name":name}).One(&result)
	return result
}

func GetObjectList(bucketName string)([]TempObject){
	var result []TempObject
	ObjectCollection.Find(bson.M{"bucket":bucketName}).All(&result)
	return result
}

func CreateObject(object Object)(bool){
	return ObjectCollection.Insert(object) == nil
}

func FindObject(bucketName string, objectName string)(bool){
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

func UpdateObject(o Object)(error){
	selector := bson.M{"object":o.Name}
	return ObjectCollection.Update(selector,o)
}

func UpdateObjectPart(objectName string, partName string)(error){
	selector := bson.M{"name":objectName}
	updater := bson.M{"part":partName}
	return ObjectCollection.Update(selector,bson.M{"$push": updater})
}

func FindParts(object string)([]Part){
	selector := bson.M{"object": object}
	var res []Part
	PartCollection.Find(selector).All(&res)
	return res
}

func FindPart(part string)(Part){
	selector := bson.M{"number": part}
	var p Part
	PartCollection.Find(selector).One(&p)
	return p
}

func GetObject(objectName string)(Object){
	selector := bson.M{"name":objectName}
	var o Object
	ObjectCollection.Find(selector).One(&o)
	return o
}

func SetObjectComplete(objectName string)(error){
	selector := bson.M{"name":objectName}
	return ObjectCollection.Update(selector,bson.M{"completed":true})
}

func RemovePart(partNumber string, objcetName string)(error){
	selector := bson.M{"number":partNumber, "object": objcetName}
	return PartCollection.Remove(selector)
}