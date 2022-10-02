package utils

import(

    "fmt"
    "time"
    "math/rand"
	"amor/database"
	"context"
	"amor/models"
	"go.mongodb.org/mongo-driver/bson"
	// "go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var userCollection *mongo.Collection = database.OpenCollection(database.Client, "user")

func EmailSuggestion(email string) [] string{
	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	var suggestions [] string
	// suggestions = append(suggestions, "4", "hola")
	for len(suggestions) < 3{
		var emailFormat = email+GenerateRandomString(3, 2)+"@amor.com"
		count, _:= userCollection.CountDocuments(ctx, bson.M{"email": emailFormat})
		defer cancel()
		if count == 0 {
			suggestions = append(suggestions, emailFormat)
		}
		

	}
	return suggestions
}

func GetUserByEmail(email string) models.User{

	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)

	var user models.User

	_ = userCollection.FindOne(ctx, bson.M{"user_id": email}).Decode(&user)
	defer cancel()

	fmt.Println(user)
	return user
}

func GenerateRandomString(length int, condition int ) (string){

	var letters [] rune
	if(condition == 0){
	letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
	}
	if(condition == 1){
		letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
		}
	if(condition == 2){
		letters = []rune("0123456789")
		}
	if(condition == 1){
		letters = []rune("abcdefghijklmnopqrstuvwxyz")
		}
		b := make([]rune, length)
		for i := range b {
			b[i] = letters[rand.Intn(len(letters))]
		}
		fmt.Println(string(b))
		return string(b)
}