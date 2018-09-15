package Structure

import "gopkg.in/mgo.v2/bson"

type Bucket struct {
	ID	 		bson.ObjectId	`bson:"_id,omitempty"`
	Name 		string 			`bson:"name" json:"name"`
	Created 	int64			`bson:"created" json:"created"`
	Modified 	int64			`bson:"modified" json:"modified"`
	Objects		[]string		`bson:"objects" json:"objects"`
}
// for returning purpose
type TempBucket struct {
	Name 		string 			`json:"name" mapstructure:"name"`
	Created 	int64			`json:"created" mapstructure:"created"`
	Modified 	int64			`json:"modified" mapstructure:"modified"`
}

type Object struct {
	ID			bson.ObjectId	`bson:"_id,omitempty"`
	Name 		string			`bson:"name" json:"name"`
	ETag		string			`bson:"etag" json:"etag"`
	Parts		[]Part			`bson:"parts" json:"part"`
	Length		int				`bson:"length" json:"length"`
	Bucket		string			`bson:"bucket" json:"bucket"`
	Created 	int64			`bson:"created" json:"created"`
	Modified 	int64			`bson:"modified" json:"modified"`
}

type TempObject struct {
	Name 		string			`bson:"name" json:"name" mapstructure:"name"`
	ETag		string			`bson:"etag" json:"eTag"  mapstructure:"eTag"`
	Created 	int64			`json:"created" mapstructure:"created"`
	Modified 	int64			`json:"modified" mapstructure:"modified"`
}

type Part struct {
	ID			bson.ObjectId	`bson:"_id,omitempty"`
	Number 		int				`bson:"number" json:"number"`
	MD5 		string			`bson:"md5" json:"md5"`
	Size		int				`bson:"size" json:"size"`
}

var ERROR = map[string]string {
	"LengthMismatched":	"LengthMismatched",
	"MD5Mismatched": "MD5Mismatched",
	"InvalidPartNumber": "InvalidPartNumber",
	"InvalidObjectName": "InvalidObjectName",
	"InvalidBucket": "InvalidBucket",
}