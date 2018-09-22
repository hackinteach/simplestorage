package Structure

import (
	"gopkg.in/mgo.v2/bson"
)

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
	Name 		string			`bson:"name" json:"name" mapstructure:"name"`
	Bucket		string			`bson:"bucket" json:"bucket" mapstructure:"bucket"`
	Completed	bool			`bson:"completed" json:"completed" mapstructure:"completed"`
	Created 	int64			`bson:"created" json:"created" mapstructure:"created"`
	Modified 	int64			`bson:"modified" json:"modified" mapstructure:"modified"`
	Part		[]int			`bson:"part" json:"part" mapstructure:"part"`
	Meta		map[string]interface{} `bson:"meta" json:"meta" mapstructure:"meta"`
}

type TempObject struct {
	Name 		string			`bson:"name" json:"name" mapstructure:"name"`
	ETag		string			`bson:"etag" json:"eTag"  mapstructure:"eTag"`
	Created 	int64			`json:"created" mapstructure:"created"`
	Modified 	int64			`json:"modified" mapstructure:"modified"`
}

type Part struct {
	ID				bson.ObjectId	`bson:"_id,omitempty"`
	Number 			int				`bson:"number" json:"number"`
	MD5 			string			`bson:"md5" json:"md5"`
	Size			int				`bson:"size" json:"size"`
	Object			string			`bson:"object" json:"object"`
	Bucket			string 			`bson:"bucket" json:"bucket"`
}

const PartNumPattern = `^([1-9][0-9]{0,3}|10000)$`
const ObjNamePattern = "(^[-_a-zA-Z0-9]{1}[-._a-zA-Z0-9]*[-_a-zA-Z0-9]{1}$|^[-_a-zA-Z0-9]{1}$)"
const BuckNamePattern = "^[-_a-zA-Z0-9]+$"