package middlewares

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

func Log() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		fmt.Printf("Request endpoint : %s %s\n", ctx.Request.Method, ctx.Request.URL.Path)
		ctx.Next()
	}
}
