package admins

import (
	"final-project/src/commons"
	"final-project/src/commons/middlewares"
	"final-project/src/modules/roles"
	"final-project/src/modules/users"

	"github.com/gin-gonic/gin"
)

func AdminRouter(router *gin.Engine) {
	adminRepository := NewRepository()
	roleRepository := roles.NewRepository()
	userRepository := users.NewRepository()
	userService := users.NewService(userRepository, roleRepository)

	adminService := NewService(adminRepository, roleRepository, userService)
	adminController := NewController(adminService)

	api := router.Group("/api/admins")
	api.Use(middlewares.JwtMiddleware())
	api.Use(middlewares.VerifyRoleMiddleware(commons.Roles.Admin))
	{
		api.POST("/users", adminController.RegisterUserController)
		api.GET("/users", adminController.GetAllUserController)
		api.GET("/users/role/:role", adminController.GetAllUserByRoleController)
		api.GET("/users/:id", adminController.GetUserByIdController)
		api.PUT("/users/:id", adminController.UpdateUserByIdController)
		api.PUT("/users/:id/role", adminController.ModifyUserRoleByIdController)
		api.PUT("/users/:id/status", adminController.ModifyUserStatusByIdController)
		api.DELETE("/users/:id", adminController.DeleteUserByIdController)
	}
}
