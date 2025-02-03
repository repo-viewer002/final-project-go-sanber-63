package auth

import (
	"github.com/gin-gonic/gin"
)

func AuthRouter(router *gin.Engine) {
	authRepository := NewRepository()

	authService := NewService(authRepository)
	authController := NewController(authService)

	api := router.Group("/api")
	api.POST("/login", authController.LoginController)
}
