package responses

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func GenerateBadRequestResponse(ctx *gin.Context, message string) {
	ctx.AbortWithStatusJSON(
		http.StatusBadRequest,
		GenerateFailMessage(message),
	)
}

func GenerateUnauthorizedResponse(ctx *gin.Context, message string) {
	ctx.AbortWithStatusJSON(
		http.StatusUnauthorized,
		GenerateFailMessage(message),
	)
}

func GenerateForbiddenResponse(ctx *gin.Context, message string) {
	ctx.AbortWithStatusJSON(
		http.StatusForbidden,
		GenerateFailMessage(message),
	)
}

func GenerateNotFoundResponse(ctx *gin.Context, message string) {
	ctx.AbortWithStatusJSON(
		http.StatusNotFound,
		GenerateFailMessage(message),
	)
}
