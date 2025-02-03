package auth

import (
	"final-project/src/commons"
	"final-project/src/commons/responses"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type Controller interface {
	LoginController(ctx *gin.Context)
}

type authController struct {
	service Service
}

func NewController(service Service) Controller {
	return &authController{
		service,
	}
}

func (controller *authController) LoginController(ctx *gin.Context) {
	var credentials Credentials
	if err := ctx.ShouldBindJSON(&credentials); err != nil {
		responses.GenerateBadRequestResponse(ctx, err.Error())

		return
	}

	token, role, err := controller.service.LoginService(credentials)

	if err != nil {
		if strings.Contains(err.Error(), "invalid credentials") {
			responses.GenerateUnauthorizedResponse(ctx, err.Error())
		} else {
			responses.GenerateBadRequestResponse(ctx, err.Error())
		}
		return
	}

	switch role {
	case commons.Roles.Admin:
		responses.GenerateSuccessResponseWithData(ctx, http.StatusOK, "admin login success", token)
	case commons.Roles.Librarian:
		responses.GenerateSuccessResponseWithData(ctx, http.StatusOK, "librarian login success", token)
	case commons.Roles.Member:
		responses.GenerateSuccessResponseWithData(ctx, http.StatusOK, "member login success", token)
	}
}
