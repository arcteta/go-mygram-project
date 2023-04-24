package routers

import (
	"go-mygram/controllers"
	"go-mygram/middleware"

	"github.com/gin-gonic/gin"
)

func StartApp() *gin.Engine {
	r := gin.Default()

	userRouter := r.Group("/users")
	{
		userRouter.POST("/register", controllers.UserRegister)
		userRouter.POST("/login", controllers.UserLogin)
	}

	photoRouter := r.Group("/photo")
	{
		photoRouter.Use(middleware.Authentication())
		photoRouter.GET("/", controllers.GetAllPhoto)
		photoRouter.GET("/:id", controllers.GetOnePhoto)
		photoRouter.POST("/", controllers.CreatePhoto)
		photoRouter.PUT("/:id", middleware.Authorization(), controllers.UpdatePhoto)
		photoRouter.DELETE("/:id", middleware.Authorization(), controllers.DeletePhoto)
	}

	commentRouter := r.Group("/comment")
	{
		commentRouter.Use(middleware.Authentication())
		commentRouter.GET("/", controllers.GetAllComment)
		commentRouter.GET("/:id", controllers.GetOneComment)
		commentRouter.POST("/:photoId", controllers.CreateComment)
		commentRouter.PUT("/:id", middleware.Authorization(), controllers.UpdateComment)
		commentRouter.DELETE("/:id", middleware.Authorization(), controllers.DeleteComment)
	}

	// socialMediaRouter := r.Group("/social-media")
	// {
	// 	socialMediaRouter.Use(middleware.Authentication())
	// 	socialMediaRouter.GET("/", controllers.GetAllSocialMedia)
	// 	socialMediaRouter.GET("/:id", controllers.GetOneSocialMedia)
	// 	socialMediaRouter.POST("/", controllers.CreateSocialMedia)
	// 	socialMediaRouter.PUT("/:id", middleware.Authorization(), controllers.UpdateSocialMedia)
	// 	socialMediaRouter.DELETE("/:id", middleware.Authorization(), controllers.DeleteSocialMedia)
	// }
	return r
}
