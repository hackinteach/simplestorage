package Mongo

import (
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"log"
	"os"
	. "simplestorage/Structure"
)

var (
	session          *mgo.Session
	BucketCollection *mgo.Collection
	ObjectCollection *mgo.Collection
	PartCollection 	 *mgo.Collection
	DB_SERVER		 string
)

const (
	DB = "SimpleStorage"
)

func init() {
	var err error

	env, f := os.LookupEnv("PROD")
	if f && env == "DOCKER"{
		DB_SERVER = "mongodb"
	}else{
		DB_SERVER = "localhost"
	}

	log.Printf("Connecting to MongoDB")

	session, err = mgo.Dial(DB_SERVER)

	if err != nil {
		panic(err)
	}

	log.Printf("DB Connected")

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
	log.Print("Removing bucket")
	var obj []Object
	ObjectCollection.Find(bson.M{"bucket":name}).All(&obj)

	for _,o := range obj {
		PartCollection.RemoveAll(bson.M{"object":o.Name})
	}

	_, oerr := ObjectCollection.RemoveAll(bson.M{"bucket":name})
	err := BucketCollection.Remove(bson.M{"name":name})
	log.Print("Done Removing bucket")
	return err == nil && oerr == nil
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
		if o.Name == objectName {
			return true
		}
	}
	return false
}

func CreatePart(part Part)(bool){
	return PartCollection.Insert(part) == nil
}

func UpdateObject(o Object)(error){
	selector := bson.M{"name":o.Name}
	err := ObjectCollection.Update(selector,o)
	if err != nil {
		log.Print(err.Error())
	}
	return err
}

func PushObjectPart(objectName string, partName int)(error){
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

func FindPart(part int, objectName, bucketName string)(Part){
	selector := bson.M{"number": part,"object": objectName, "bucket":bucketName}
	var p Part
	PartCollection.Find(selector).One(&p)
	return p
}

func GetObject(objectName, bucket string)(Object, bool){
	selector := bson.M{"name":objectName}
	var o Object
	ObjectCollection.Find(selector).One(&o)
	if o.Bucket == bucket{
		return o, true
	}
	return o, false
}

func SetObjectComplete(o Object)(error){
	selector := bson.M{"name":o.Name}
	updater := bson.M{"completed": o.Completed}
	return ObjectCollection.Update(selector,bson.M{"$set":updater})
}

func RemovePart(partNumber int, objectName string)(error){
	selector := bson.M{"number":partNumber, "object": objectName}
	return PartCollection.Remove(selector)
}

func RemoveObjectDB(o Object) error{
	selector := bson.M{"name":o.Name, "bucket": o.Bucket}
	for _,p := range o.Part {
		PartCollection.RemoveAll(bson.M{"number":p})
	}
	return ObjectCollection.Remove(selector)
}