package roles

import (
	"final-project/src/commons/middlewares"
	"final-project/src/commons/responses"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type Controller interface {
	CreateRoleController(ctx *gin.Context)
	GetAllRoleController(ctx *gin.Context)
	GetRoleByIdController(ctx *gin.Context)
	UpdateRoleByIdController(ctx *gin.Context)
	DeleteRoleByIdController(ctx *gin.Context)
}

type roleController struct {
	service Service
}

func NewController(service Service) Controller {
	return &roleController{
		service,
	}
}

func (controller *roleController) CreateRoleController(ctx *gin.Context) {
	_, username, _, err := middlewares.GetClaims(ctx)

	if err != nil {
		responses.GenerateUnauthorizedResponse(ctx, err.Error())

		return
	}

	var role Role

	if err := ctx.ShouldBindJSON(&role); err != nil {
		responses.GenerateBadRequestResponse(ctx, err.Error())

		return
	}

	role.Created_By = username
	role.Modified_By = username

	createdRole, err := controller.service.CreateRoleService(role)

	if err != nil {
		responses.GenerateBadRequestResponse(ctx, err.Error())

		return
	}

	responses.GenerateSuccessResponseWithData(ctx, http.StatusCreated, "create role success", createdRole)
}

func (controller *roleController) GetAllRoleController(ctx *gin.Context) {
	role, err := controller.service.GetAllRoleService()

	if err != nil {
		responses.GenerateBadRequestResponse(ctx, err.Error())

		return
	}

	responses.GenerateSuccessResponseWithData(ctx, http.StatusOK, "get all role success", role)
}

func (controller *roleController) GetRoleByIdController(ctx *gin.Context) {
	getId := ctx.Param("id")

	role, err := controller.service.GetRoleByIdService(getId)

	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			responses.GenerateNotFoundResponse(ctx, err.Error())
		} else {
			responses.GenerateBadRequestResponse(ctx, err.Error())
		}

		return
	}

	responses.GenerateSuccessResponseWithData(ctx, http.StatusOK, fmt.Sprintf("get role by id \"%s\" success", getId), role)
}

func (controller *roleController) UpdateRoleByIdController(ctx *gin.Context) {
	_, username, _, err := middlewares.GetClaims(ctx)

	if err != nil {
		responses.GenerateUnauthorizedResponse(ctx, err.Error())

		return
	}

	var role Role

	getId := ctx.Param("id")

	if err := ctx.ShouldBindJSON(&role); err != nil {
		responses.GenerateNotFoundResponse(ctx, err.Error())

		return
	}

	role.Modified_By = username
	updatedRole, err := controller.service.UpdateRoleByIdService(getId, role)

	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			responses.GenerateNotFoundResponse(ctx, err.Error())
		} else {
			responses.GenerateBadRequestResponse(ctx, err.Error())
		}

		return
	}

	responses.GenerateSuccessResponseWithData(ctx, http.StatusOK, fmt.Sprintf("update role by id \"%s\" success", getId), updatedRole)
}

func (controller *roleController) DeleteRoleByIdController(ctx *gin.Context) {
	getId := ctx.Param("id")

	deletedRole, err := controller.service.DeleteRoleByIdService(getId)

	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			responses.GenerateNotFoundResponse(ctx, err.Error())
		} else {
			responses.GenerateBadRequestResponse(ctx, err.Error())
		}

		return
	}

	responses.GenerateSuccessResponseWithData(ctx, http.StatusOK, fmt.Sprintf("delete role by id \"%s\" success", getId), deletedRole)
}
