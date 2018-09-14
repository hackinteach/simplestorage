package Structure

import "gopkg.in/mgo.v2/bson"

type Bucket struct {
	ID	 		bson.ObjectId	`bson:"_id,omitempty"`
	Name 		string 			`bson:"name" json:"name"`
	Created 	int64			`bson:"created" json:"created"`
	Modified 	int64			`bson:"modified" json:"modified"`
}
// for returning purpose
type TempBucket struct {
	Name 		string 			`json:"name"`
	Created 	int64			`json:"created"`
	Modified 	int64			`json:"modified"`
}

type Object struct {
	ID			bson.ObjectId	`bson:"_id"`
	Name 		string			`bson:"name",json:"name"`
	ETag		string			`bson:"etag",json:"etag"`
	Parts		[]Part			`bson:"parts",json:"part"`
	Length		int				`bson:"length",json:"length"`
}

type Part struct {
	ID			bson.ObjectId	`bson:"_id"`
	Number 		int				`bson:"number",json:"number"`
	MD5 		string			`bson:"md5",json:"md5"`
	Size		int				`bson:"size",json:"size"`
}