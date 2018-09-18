package Structure

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"gopkg.in/mgo.v2/bson"
	"simplestorage/Mongo"
	"strings"
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
	Name 		string			`bson:"name" json:"name"`
	Bucket		string			`bson:"bucket" json:"bucket"`
	Completed	bool			`bson:"completed" json:"completed"`
	Created 	int64			`bson:"created" json:"created"`
	Modified 	int64			`bson:"modified" json:"modified"`
	Part		[]string		`bson:"part" json:"part"`
}

type TempObject struct {
	Name 		string			`bson:"name" json:"name" mapstructure:"name"`
	ETag		string			`bson:"etag" json:"eTag"  mapstructure:"eTag"`
	Created 	int64			`json:"created" mapstructure:"created"`
	Modified 	int64			`json:"modified" mapstructure:"modified"`
}

type Part struct {
	ID			bson.ObjectId	`bson:"_id,omitempty"`
	Number 		string			`bson:"number" json:"number"`
	MD5 		string			`bson:"md5" json:"md5"`
	Size		int				`bson:"size" json:"size"`
	Object		string			`bson:"object" json:"object"`
}

const PartNumPattern = `^([1-9][0-9]{0,3}|10000)$`
const ObjNamePattern = `^(?!\.).[ \. | \_ | -| a-z | 0-9]*`
const BuckNamePattern = `(^(?!\.)([a-z|1-9|\-|\_]){2,})`

func (o Object) Etag() string {
	hasher := md5.New()
	parts := Mongo.FindParts(o.Name)
	var md5 []string

	for _,p := range parts {
		md5 = append(md5, p.MD5)
	}

	md5s := strings.Join(md5,"")
	hasher.Write([]byte(md5s))
	hashed := hex.EncodeToString(hasher.Sum(nil))
	return fmt.Sprintf("%s-%d",hashed,len(md5))
}

func (o Object) Length() int {
	length := 0
	for _,v := range o.Part {
		p := Mongo.FindPart(v)
		length += p.Size
	}
	return length
}