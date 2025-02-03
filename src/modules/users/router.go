package users

import (
	"final-project/src/commons/middlewares"
	"final-project/src/modules/roles"

	"github.com/gin-gonic/gin"
)

func UserRouter(router *gin.Engine) {
	roleRepository := roles.NewRepository()
	userRepository := NewRepository()
	userService := NewService(userRepository, roleRepository)
	userController := NewController(userService)

	api := router.Group("/api")

	api.Use(middlewares.JwtMiddleware())
	{
		api.GET("/profile", userController.ViewProfileController)
		api.PUT("/profile", userController.UpdateProfileController)
	}
}
