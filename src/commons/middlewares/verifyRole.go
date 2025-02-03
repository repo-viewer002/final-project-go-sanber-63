package middlewares

import (
	"final-project/src/commons/responses"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func VerifyRoleMiddleware(allowedRoles ...string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		claims, exists := ctx.Get("user")

		if !exists {
			responses.GenerateUnauthorizedResponse(ctx, "invalid authorization header format")
			ctx.Abort()
			return
		}

		mapClaims, ok := claims.(jwt.MapClaims)

		if !ok {
			responses.GenerateUnauthorizedResponse(ctx, "unauthorized access: invalid token format")
			ctx.Abort()
			return
		}

		role, roleExists := mapClaims["role"].(string)
		if !roleExists {
			responses.GenerateUnauthorizedResponse(ctx, "unauthorized access: role not found")
			ctx.Abort()
			return
		}

		for _, allowedRole := range allowedRoles {
			if role == allowedRole {
				ctx.Next()
				return
			}
		}

		responses.GenerateForbiddenResponse(ctx, "unauthorized access: insufficient permissions")
		ctx.Abort()
	}
}
