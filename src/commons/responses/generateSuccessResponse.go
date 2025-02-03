package responses

import (
	"github.com/gin-gonic/gin"
)

func GenerateSuccessResponse(ctx *gin.Context, status int, message string) {
	ctx.JSON(
		status,
		GenerateSuccessMessage(message),
	)
}

func GenerateSuccessResponseWithData(ctx *gin.Context, status int, message string, data interface{}) {
	ctx.JSON(
		status,
		GenerateSuccessMessageWithData(message, data),
	)
}
