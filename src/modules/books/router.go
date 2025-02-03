package books

import (
	"final-project/src/commons"
	"final-project/src/commons/middlewares"

	"github.com/gin-gonic/gin"
)

func BookRouter(router *gin.Engine) {
	repository := NewRepository()
	service := NewService(repository)
	controller := NewController(service)

	api := router.Group("/api/books")
	api.Use(middlewares.JwtMiddleware())

	api.GET("", controller.GetAllBookController)
	api.GET("/genres", controller.GetAllBookByGenreController)
	api.GET("/:bookId", controller.GetBookByIdController)

	api.Use(middlewares.VerifyRoleMiddleware(commons.Roles.Admin, commons.Roles.Librarian))
	{
		api.POST("", controller.CreateBookController)
		api.PUT("/:bookId", controller.UpdateBookByIdController)
		api.DELETE("/:bookId", controller.DeleteBookByIdController)
	}
}
