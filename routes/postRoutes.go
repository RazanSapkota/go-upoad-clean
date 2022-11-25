package routes

import (
	"example/go-api/controllers"

	"github.com/gin-gonic/gin"
)

func PostRoutes(router *gin.Engine) {
	router.POST("/post", controllers.CreatePost)
	router.GET("/post", controllers.GetPosts)
	router.GET("/post/:id", controllers.FetchOnePost)
	router.PUT("/post/:id", controllers.UpdatePost)
}