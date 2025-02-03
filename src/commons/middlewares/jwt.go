package middlewares

import (
	"errors"
	"final-project/src/commons"
	"final-project/src/commons/responses"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func JwtMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		tokenString, err := getTokenFromHeader(ctx)

		if err != nil {
			responses.GenerateUnauthorizedResponse(ctx, err.Error())

			ctx.Abort()

			return
		}

		token, tokenValidationError := verifyToken(tokenString)

		if tokenValidationError != nil {
			responses.GenerateUnauthorizedResponse(ctx, tokenValidationError.Error())

			ctx.Abort()

			return
		}

		if claims, ok := token.Claims.(jwt.MapClaims); !ok {
			responses.GenerateUnauthorizedResponse(ctx, "invalid token")

			ctx.Abort()

			return
		} else {
			ctx.Set("user", claims)
		}

		ctx.Next()
	}
}

func CreateToken(id string, username string, email string, role string) (string, error) {
	claims := jwt.NewWithClaims(jwt.SigningMethodHS256,
		jwt.MapClaims{
			"sub":      id,
			"username": username,
			"email":    email,
			"role":     role,
			"iss":      "libraryApiServer",
			"aud":      "libraryApiClient",
			"exp":      time.Now().Add(time.Hour).Unix(),
			"iat":      time.Now().Unix(),
		})

	token, err := claims.SignedString(commons.JWT_SECRET_KEY)

	if err != nil {
		return "", err
	}

	return token, nil
}

func verifyToken(tokenString string) (*jwt.Token, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return commons.JWT_SECRET_KEY, nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, err
	}

	return token, nil
}

func getTokenFromHeader(ctx *gin.Context) (string, error) {
	authHeader := ctx.GetHeader("Authorization")

	if authHeader == "" {
		return authHeader, errors.New("authorization header is required")
	}

	splitAuthHeader := strings.Split(authHeader, " ")
	if len(splitAuthHeader) != 2 || splitAuthHeader[0] != "Bearer" {
		return authHeader, errors.New("invalid authorization header format")
	}

	return splitAuthHeader[1], nil
}

// return is for : id, username, roleId, error
func GetClaims(ctx *gin.Context) (string, string, string, error) {
	claims, exists := ctx.Get("user")

	if !exists {
		return "", "", "", errors.New("invalid authorization header format")
	}

	mapClaims, ok := claims.(jwt.MapClaims)

	if !ok {
		responses.GenerateUnauthorizedResponse(ctx, "unauthorized access: invalid token format")
		ctx.Abort()
		return "", "", "", errors.New("invalid token format")
	}

	id := mapClaims["sub"].(string)
	role := mapClaims["role"].(string)
	username := mapClaims["username"].(string)

	return id, username, role, nil
}
