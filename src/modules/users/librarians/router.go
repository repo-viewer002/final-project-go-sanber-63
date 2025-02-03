package librarians

import (
	"final-project/src/commons"
	"final-project/src/commons/middlewares"
	"final-project/src/modules/roles"
	"final-project/src/modules/users"

	"github.com/gin-gonic/gin"
)

func LibrarianRouter(router *gin.Engine) {
	roleRepository := roles.NewRepository()
	roleService := roles.NewService(roleRepository)
	userRepository := users.NewRepository()
	userService := users.NewService(userRepository, roleRepository)
	librarianRepository := NewRepository()
	librarianService := NewService(librarianRepository, userService, roleService)
	librarianController := NewController(librarianService)

	api := router.Group("/api/members")
	api.Use(middlewares.JwtMiddleware())
	api.Use(middlewares.VerifyRoleMiddleware(commons.Roles.Admin, commons.Roles.Librarian))
	{
		api.POST("/", librarianController.CreateMemberController)
		api.GET("/", librarianController.GetAllMemberController)
		api.GET("/:memberId", librarianController.GetMemberByIdController)
		api.PUT("/:memberId", librarianController.UpdateMemberByIdController)
	}
}
