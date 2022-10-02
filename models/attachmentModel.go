package models

import(
	"time"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Attachment struct{

	ID 		primitive.ObjectID		`bson:"_id"`
	Name string					`json:"name_id"`
	Location	string			`json:"location"`
	Attachment_id	string			`json:"attachment_id"`
	Message_id	string			`json:"message_id"`
	Created_at 	time.Time			`json:"created_at"`
	Updated_at 	time.Time			`json:"updated_at"`
}