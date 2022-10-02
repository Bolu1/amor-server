package routes

import(
	"github.com/gin-gonic/gin"
	controller "amor/controllers"
)

func MessageRoutes(incomingRoutes *gin.Engine){

	incomingRoutes.POST("/message/new", controller.Add())
	incomingRoutes.GET("/message/sent", controller.GetSent())
	incomingRoutes.GET("/message/received", controller.GetReceived())
	incomingRoutes.GET("/message/message/:id", controller.GetOne())
	incomingRoutes.GET("/message/attachment/:id", controller.GetAttachmentByFileId())

}