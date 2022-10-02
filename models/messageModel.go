package models

import(
	"time"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Message struct{

	ID 		primitive.ObjectID		`bson:"_id"`
	Message string					`json:"message"`
	Parent_id string					`json:"parent_id"`
	Sender_id	string			`json:"sender_id"`
	Receiver_id	string			`json:"receiver_id"`
	Message_id	string			`json:"message_id"`
	Created_at 	time.Time			`json:"created_at"`
	Updated_at 	time.Time			`json:"updated_at"`
}