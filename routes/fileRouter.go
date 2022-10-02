package routes

import(
	"github.com/gin-gonic/gin"
	controller "amor/controllers"
)


func FileRoutes(incomingRoutes *gin.Engine){

	incomingRoutes.GET("/files", controller.GetFiles())
	incomingRoutes.GET("/file/:id", controller.GetFile())
	incomingRoutes.POST("/file", controller.AddFile())
	incomingRoutes.DELETE("/file/:id", controller.DeleteFile())
	incomingRoutes.PATCH("/user/profile", controller.ChangeProfile())
}