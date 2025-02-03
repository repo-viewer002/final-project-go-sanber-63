package roles

import (
	"final-project/src/commons"
	"final-project/src/commons/middlewares"

	"github.com/gin-gonic/gin"
)

func RoleRouter(router *gin.Engine) {
	repository := NewRepository()
	service := NewService(repository)
	controller := NewController(service)

	api := router.Group("/api/roles")
	api.Use(middlewares.JwtMiddleware())
	api.Use(middlewares.VerifyRoleMiddleware(commons.Roles.Admin))
	{
		api.POST("", controller.CreateRoleController)
		api.GET("", controller.GetAllRoleController)
		api.GET("/:id", controller.GetRoleByIdController)
		api.PUT("/:id", controller.UpdateRoleByIdController)
		api.DELETE("/:id", controller.DeleteRoleByIdController)
	}
}
