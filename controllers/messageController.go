package controllers

import(
	"github.com/gin-gonic/gin"
	"context"
	"amor/database"
	"amor/models"
	"net/http"
	"fmt"
	// "path/filepath"
	"time"
	"log"
	// "math"
	"strconv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/bson/primitive"
	// "go.mongodb.org/mongo-driver/mongo/options"
	// "github.com/go-playground/validator/v10"
)


var messageCollection *mongo.Collection = database.OpenCollection(database.Client, "message")
var attachmentCollection *mongo.Collection = database.OpenCollection(database.Client, "attachment")

func Add() gin.HandlerFunc{

	return func(c *gin.Context){
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()
		var message models.Message
		// var attachment models.Attachment

		if err := c.BindJSON(&message); err != nil{
			fmt.Println(err.Error())
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		uid := c.MustGet("uid").(string)

		message.Created_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		message.Updated_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		message.ID = primitive.NewObjectID()
		message.Message_id = message.ID.Hex()
		message.Sender_id = uid

		_, insertErr := messageCollection.InsertOne(ctx, message)
		if insertErr != nil{
			msg := fmt.Sprintf("Message was not sent")
			c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
			return
		}

		// file, err := c.FormFile("file")
		// if err != nil {
		// 	c.String(http.StatusBadRequest, "get form err: %s", err.Error())
		// 	return
		// }
		// if(file != nil){

		// 	size := float64(file.Size)/1000000
		// if(size > 0.1){
		// 	c.JSON(http.StatusInternalServerError, gin.H{"error": "Attachment too large"})
		// return
		// }

		// date := time.Now().Format(time.RFC850)
		// filename := filepath.Base(file.Filename)
		// location := "public/"+date+"-"+filename
		// if err := c.SaveUploadedFile(file, location); err != nil {
		// 	c.String(http.StatusBadRequest, "upload file err: %s", err.Error())
		// 	return
		// }

		// attachment.Name = filename
		// attachment.Location = location
		// attachment.Message_id = message.ID.Hex()
		// attachment.Created_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		// attachment.Updated_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		// attachment.ID = primitive.NewObjectID()
		// attachment.Attachment_id = attachment.ID.Hex()
		
		// _, insertErr := attachmentCollection.InsertOne(ctx, attachment)
		// if insertErr != nil{
		// 	msg := fmt.Sprintf("File was not added")
		// 	c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
		// 	return
		// }
		// }

		c.JSON(http.StatusOK, gin.H{"message":"sent"})


	}

}

func GetSent() gin.HandlerFunc{
	return func(c *gin.Context){

		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)

		recordPerPage, err := strconv.Atoi(c.Query("recordPerPage"))
		if err != nil || recordPerPage < 1 {
			recordPerPage = 10
		}

		page, err := strconv.Atoi(c.Query("page"))
		if err != nil || page < 1 {
			page = 1
		}


		uid := c.MustGet("uid").(string)

		startIndex := (page - 1) * recordPerPage
		startIndex, err = strconv.Atoi(c.Query("startIndex"))

		matchStage := bson.D{{"$match", bson.D{{"sender_id", uid}}}}
		groupStage := bson.D{{"$group", bson.D{{"_id", bson.D{{"_id", "null"}}}, {"total_count", bson.D{{"$sum", 1}}}, {"data", bson.D{{"$push", "$$ROOT"}}}}}}
		projectStage := bson.D{
			{
				"$project", bson.D{
					{"_id", 0},
					{"total_count", 1},
					{"messages", bson.D{{"$slice", []interface{}{"$data", startIndex, recordPerPage}}}},
				}}}

				
				sortStage := bson.D{{"$sort", bson.D{{"_id", -1}}}}
		result, err := messageCollection.Aggregate(ctx, mongo.Pipeline{
			matchStage,sortStage,  groupStage, projectStage})
		defer cancel()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error occured while listing messages items"})
			return
		}
		var allFiles []bson.M
		if err = result.All(ctx, &allFiles); err != nil {
			log.Fatal(err)
		}
		c.JSON(http.StatusOK, allFiles[0])
	}
}

func GetReceived() gin.HandlerFunc{
	return func(c *gin.Context){

		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)

		recordPerPage, err := strconv.Atoi(c.Query("recordPerPage"))
		if err != nil || recordPerPage < 1 {
			recordPerPage = 10
		}

		page, err := strconv.Atoi(c.Query("page"))
		if err != nil || page < 1 {
			page = 1
		}


		uid := c.MustGet("uid").(string)

		startIndex := (page - 1) * recordPerPage
		startIndex, err = strconv.Atoi(c.Query("startIndex"))

		matchStage := bson.D{{"$match", bson.D{{"receiver_id", uid}}}}
		groupStage := bson.D{{"$group", bson.D{{"_id", bson.D{{"_id", "null"}}}, {"total_count", bson.D{{"$sum", 1}}}, {"data", bson.D{{"$push", "$$ROOT"}}}}}}
		projectStage := bson.D{
			{
				"$project", bson.D{
					{"_id", 0},
					{"total_count", 1},
					{"messages", bson.D{{"$slice", []interface{}{"$data", startIndex, recordPerPage}}}},
				}}}

				
				sortStage := bson.D{{"$sort", bson.D{{"_id", -1}}}}
		result, err := messageCollection.Aggregate(ctx, mongo.Pipeline{
			matchStage,sortStage,  groupStage, projectStage})
		defer cancel()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error occured while listing messages items"})
			return
		}
		var allFiles []bson.M
		if err = result.All(ctx, &allFiles); err != nil {
			log.Fatal(err)
		}
		c.JSON(http.StatusOK, allFiles[0])
	}
}

func GetOne() gin.HandlerFunc{

	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		id := c.Param("id")
		var message models.Message

		err := messageCollection.FindOne(ctx, bson.M{"message_id": id}).Decode(&message)

		defer cancel()
		if err != nil{
			fmt.Println(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error occured while lisiting user messages"})
			return
		}
		c.JSON(http.StatusOK, message)
	}
}

func GetAttachmentByFileId() gin.HandlerFunc{

	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		id := c.Param("id")
		var attachment models.Attachment

		err := attachmentCollection.FindOne(ctx, bson.M{"message_id": id}).Decode(&attachment)

		defer cancel()
		if err != nil{
			fmt.Println(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error occured while lisiting user attachments"})
			return
		}
		c.JSON(http.StatusOK, attachment)
	}
}