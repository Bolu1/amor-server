package models

import(
	"time"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type File struct{
	ID 		primitive.ObjectID		`bson:"_id"`
	Name	string					`json:"name" validate:"required,min=2,max=100"`
	Location	string					`json:"location" validate:"required,min=2,max=100"`
	Size	float64				`json:"size" validate:"required"`
	Type	string					`json:"type"`
	User_id string					`json:"user_id"`
	File_id	string					`json:"fileId"`
	Created_at 	time.Time						`json:"created_at"`
	Updated_at 	time.Time				`json:"updated_at"`
}