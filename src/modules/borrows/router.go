package borrows

import (
	"final-project/src/commons"
	"final-project/src/commons/middlewares"

	"github.com/gin-gonic/gin"
)

func BorrowRouter(router *gin.Engine) {
	repository := NewRepository()
	service := NewService(repository)
	controller := NewController(service)

	api := router.Group("/api")
	api.Use(middlewares.JwtMiddleware())
	api.Use(middlewares.VerifyRoleMiddleware(commons.Roles.Admin, commons.Roles.Librarian))
	{
		api.POST("/borrow", controller.BorrowBookController)
		api.POST("/return/:borrowId", controller.ReturnBookController)
	}
}
