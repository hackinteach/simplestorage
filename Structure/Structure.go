package Structure

import "gopkg.in/mgo.v2/bson"

type Bucket struct {
	ID	 		bson.ObjectId	`bson:"_id,omitempty"`
	Name 		string 			`bson:"name" json:"name"`
	Created 	int64			`bson:"created" json:"created"`
	Modified 	int64			`bson:"modified" json:"modified"`
}

type TempBucket struct {
	Name 		string 			`json:"name"`
	Created 	int64			`json:"created"`
	Modified 	int64			`json:"modified"`
}