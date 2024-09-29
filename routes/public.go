package routes

import (
	"github.com/amirnilofari/hash-go-mysql/controllers"
	"github.com/gin-gonic/gin"
)

func PublicRoutes(router *gin.Engine) {
	router.POST("/register", controllers.Register)
	router.POST("/login", controllers.Login)
	router.GET("/posts", controllers.GetPosts)
	router.GET("/posts/:id", controllers.GetPost)
	router.GET("/posts/:id/comments", controllers.GetComments)
}
