package users

import (
	"final-project/src/commons/middlewares"
	"final-project/src/commons/responses"
	"final-project/src/utils"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type Controller interface {
	ViewProfileController(ctx *gin.Context)
	UpdateProfileController(ctx *gin.Context)
}

type userController struct {
	service Service
}

func NewController(service Service) Controller {
	return &userController{
		service,
	}
}

func (controller *userController) ViewProfileController(ctx *gin.Context) {
	id, _, _, err := middlewares.GetClaims(ctx)

	if err != nil {
		responses.GenerateUnauthorizedResponse(ctx, err.Error())
	}

	profile, err := controller.service.ViewProfileService(id)

	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			responses.GenerateNotFoundResponse(ctx, err.Error())
		} else {
			responses.GenerateBadRequestResponse(ctx, err.Error())
		}

		return
	}

	responses.GenerateSuccessResponseWithData(ctx, http.StatusOK, "view profile success", profile)
}

func (controller *userController) UpdateProfileController(ctx *gin.Context) {
	id, username, role, err := middlewares.GetClaims(ctx)

	if err != nil {
		responses.GenerateUnauthorizedResponse(ctx, err.Error())

		return
	}

	var user UpdateUserDTO

	if err := ctx.ShouldBindJSON(&user); err != nil {
		responses.GenerateNotFoundResponse(ctx, err.Error())

		return
	}

	utils.GenerateDataModifier(role, username, &user.Modified_By)

	updatedProfile, err := controller.service.UpdateProfileService(id, user)

	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			responses.GenerateNotFoundResponse(ctx, err.Error())
		} else {
			responses.GenerateBadRequestResponse(ctx, err.Error())
		}

		return
	}

	responses.GenerateSuccessResponseWithData(ctx, http.StatusOK, "update profile success", updatedProfile)
}
