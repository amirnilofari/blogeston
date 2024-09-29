package routes

import (
	"github.com/amirnilofari/hash-go-mysql/controllers"
	"github.com/amirnilofari/hash-go-mysql/middlewares"
	"github.com/gin-gonic/gin"
)

func PrivateRoutes(router *gin.Engine) {
	private := router.Group("/")
	private.Use(middlewares.JWTAuthMiddleware())
	{
		private.POST("/create-post", controllers.CreatePost)
		private.POST("/posts/:id/create-comment", controllers.CreateComment)
		private.POST("/comments/:comment_id/react", controllers.ReactToComment)

		private.Use(middlewares.RoleMiddleware("admin"))
		{
			private.GET("/users", controllers.GetAllUsers)
			private.POST("/posts/:id/publish", controllers.PublishPost)
		}
	}

}
