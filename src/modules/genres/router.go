package genres

import (
	"final-project/src/commons"
	"final-project/src/commons/middlewares"

	"github.com/gin-gonic/gin"
)

func GenreRouter(router *gin.Engine) {
	repository := NewRepository()
	service := NewService(repository)
	controller := NewController(service)

	api := router.Group("/api/genres")
	api.Use(middlewares.JwtMiddleware())

	api.GET("", controller.GetAllGenreController)
	api.GET("/:id", controller.GetGenreByIdController)

	api.Use(middlewares.VerifyRoleMiddleware(commons.Roles.Admin, commons.Roles.Librarian))
	{
		api.POST("", controller.CreateGenreController)
		api.PUT("/:id", controller.UpdateGenreByIdController)
		api.DELETE("/:id", controller.DeleteGenreByIdController)
	}
}
