package controllers

import(
	"github.com/gin-gonic/gin"
	"context"
	"amor/database"
	"amor/models"
	"net/http"
	"fmt"
	"path/filepath"
	"time"
	"log"
	"math"
	"strconv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
	"github.com/go-playground/validator/v10"
)

var fileCollection *mongo.Collection = database.OpenCollection(database.Client, "file")
var validate = validator.New()

func GetFiles() gin.HandlerFunc{
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

		matchStage := bson.D{{"$match", bson.D{{"user_id", uid}}}}
		groupStage := bson.D{{"$group", bson.D{{"_id", bson.D{{"_id", "null"}}}, {"total_count", bson.D{{"$sum", 1}}}, {"data", bson.D{{"$push", "$$ROOT"}}}}}}
		projectStage := bson.D{
			{
				"$project", bson.D{
					{"_id", 0},
					{"total_count", 1},
					{"file_items", bson.D{{"$slice", []interface{}{"$data", startIndex, recordPerPage}}}},
				}}}

		result, err := fileCollection.Aggregate(ctx, mongo.Pipeline{
			matchStage, groupStage, projectStage})
		defer cancel()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error occured while listing files items"})
			return
		}
		var allFiles []bson.M
		if err = result.All(ctx, &allFiles); err != nil {
			log.Fatal(err)
		}
		c.JSON(http.StatusOK, allFiles[0])
	}
}

func GetFile() gin.HandlerFunc{
	return func(c *gin.Context){
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		id := c.Param("id")
		var file models.File

		err := fileCollection.FindOne(ctx, bson.M{"id": id}).Decode(&file)
		defer cancel()
		if err != nil{
			c.JSON(http.StatusInternalServerError, gin.H{"error":"error occured while fetching the file"})
			return
		}
		c.JSON(http.StatusOK, file)
	}

}


func round(num float64) int {
	return int(num + math.Copysign(0.5, num))
}

func toFixed(num float64, precision int) float64{
	output := math.Pow(10, float64(precision))
	return float64(round(num*output)) / output
}

func AddFile() gin.HandlerFunc{
	return func(c *gin.Context){
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		var dbfile models.File

		file, err := c.FormFile("file")
		if err != nil {
			c.String(http.StatusBadRequest, "get form err: %s", err.Error())
			return
		}

		// parse file size
		size := float64(file.Size)/1000000
		extension := filepath.Ext(file.Filename)
		
		// check if there is space left
		var user models.User
		uid := c.MustGet("uid").(string)

		_ = userCollection.FindOne(ctx, bson.M{"user_id": uid}).Decode(&user)
		userSpace := float64(user.Space_Left)
		if(size > userSpace){
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Insufficient space"})
			return
		}

		updateUserSpace(uid, userSpace-size)


		// parse file name
		date := time.Now().Format(time.RFC850)
		filename := filepath.Base(file.Filename)
		location := "public/"+date+"-"+filename
		if err := c.SaveUploadedFile(file, location); err != nil {
			c.String(http.StatusBadRequest, "upload file err: %s", err.Error())
			return
		}
		
		dbfile.Name = filename
		dbfile.Location = location
		dbfile.Type = extension
		dbfile.User_id = uid
		dbfile.Size = size
		dbfile.Created_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		dbfile.Updated_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		dbfile.ID = primitive.NewObjectID()
		dbfile.File_id = dbfile.ID.Hex()
		
		result, insertErr := fileCollection.InsertOne(ctx, dbfile)
		if insertErr != nil{
			msg := fmt.Sprintf("File was not added")
			c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
			return
		}
		c.JSON(http.StatusOK, result)

		
	}
}

func updateUserSpace(id string, amount float64) (err1 string) {

	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()
	var user models.User

	filter := bson.M{"user_id": id}

	var updateObj primitive.D
	updateObj = append(updateObj, bson.E{"space_left", amount})

	user.Updated_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	updateObj = append(updateObj, bson.E{"updated_at", user.Updated_at})

	upsert := true

		opt := options.UpdateOptions{
			Upsert: &upsert,
		}
		_, err := userCollection.UpdateOne(
			ctx,
			filter,
			bson.D{{"$set", updateObj}},
			&opt,
		)

		if err != nil {
			fmt.Println(err.Error())
			err1 = err.Error()
		}
		return err1

}

func DeleteFile() gin.HandlerFunc{
	return func(c *gin.Context){

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		fileId := c.Param("file_id")
		filter := bson.M{"fileId": fileId}
		result, err := fileCollection.DeleteOne(ctx, filter)
		res := map[string]interface{}{"data": result}

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": err})
			return
		}

		if result.DeletedCount < 1 {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "No data to delete"})
			return
		}

		c.JSON(http.StatusCreated, gin.H{"message": "File deleted successfully", "Data": res})
	}
} 

func ChangeProfile() gin.HandlerFunc{

	return func (c *gin.Context){
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		uid := c.MustGet("uid").(string)
		

		file, err := c.FormFile("file")
		if err != nil {
			c.String(http.StatusBadRequest, "get form err: %s", err.Error())
			return
		}

		// parse file size
		size := float64(file.Size)/1000000
		extension := filepath.Ext(file.Filename)
		if(extension != ".png" && extension != ".jpg"){
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Not an image"})
			return
		}
		
		if(size > 0.7){
			c.JSON(http.StatusInternalServerError, gin.H{"error": "File too large"})
			return
		}

		// parse file name
		date := time.Now().Format(time.RFC850)
		filename := filepath.Base(file.Filename)
		location := "public/"+date+"-"+filename
		if err := c.SaveUploadedFile(file, location); err != nil {
			c.String(http.StatusBadRequest, "upload file err: %s", err.Error())
			return
		}

		// update image

	var user models.User
	filter := bson.M{"user_id": uid}



	var updateObj primitive.D
	updateObj = append(updateObj, bson.E{"avatar", location})

	user.Updated_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	updateObj = append(updateObj, bson.E{"updated_at", user.Updated_at})

	upsert := true

		opt := options.UpdateOptions{
			Upsert: &upsert,
		}
		_, insertErr := userCollection.UpdateOne(
			ctx,
			filter,
			bson.D{{"$set", updateObj}},
			&opt,
		)

		if insertErr != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not update profile"})
			return
		}
		
		c.JSON(http.StatusOK, gin.H{"message":"Profile changed"})

		
	}
}