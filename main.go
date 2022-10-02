package main

import(
	"os"
	"amor/database"
	"net/http"
	"amor/server"
	"log"
	routes "amor/routes"
	middleware "amor/middleware"
	"go.mongodb.org/mongo-driver/mongo"
	"github.com/gin-gonic/gin"
)



var fileCollection *mongo.Collection = database.OpenCollection(database.Client, "file")

func main(){
	port := os.Getenv("PORT")
	wport := os.Getenv("PORT")

	if port == ""{
		port = "8000"
	}

	if wport == ""{
		wport = "8080"
	}

	go func(){
		server.AllRooms.Init()

		http.HandleFunc("/create", server.CreateRoomHandler)
		http.HandleFunc("/join", server.JoinRoomRequestHandler)
	
		log.Println("Starting server on port 8000")
		err := http.ListenAndServe(":8080", nil)
		if err != nil{
			log.Fatal(err)
		}
	}()

	router := gin.New()
	router.Static("/public", "./public")
	router.MaxMultipartMemory = 8 << 20
	router.Use(gin.Logger())
	routes.UserRoutes(router)
	router.Use(middleware.Authentication())

	routes.FileRoutes(router)
	routes.MessageRoutes(router)

	router.Run(":" + port)
}