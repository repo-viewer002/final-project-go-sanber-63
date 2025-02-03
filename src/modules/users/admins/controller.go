package admins

import (
	"final-project/src/commons/middlewares"
	"final-project/src/commons/responses"
	"final-project/src/modules/users"
	"final-project/src/utils"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type Controller interface {
	RegisterUserController(ctx *gin.Context)
	GetAllUserController(ctx *gin.Context)
	GetAllUserByRoleController(ctx *gin.Context)
	GetUserByIdController(ctx *gin.Context)
	UpdateUserByIdController(ctx *gin.Context)
	ModifyUserStatusByIdController(ctx *gin.Context)
	ModifyUserRoleByIdController(ctx *gin.Context)
	DeleteUserByIdController(ctx *gin.Context)
}

type adminController struct {
	service Service
}

func NewController(service Service) Controller {
	return &adminController{
		service,
	}
}

func (controller *adminController) RegisterUserController(ctx *gin.Context) {
	_, username, _, err := middlewares.GetClaims(ctx)

	if err != nil {
		responses.GenerateUnauthorizedResponse(ctx, err.Error())

		return
	}

	var user users.RegisterUserDTO

	if err := ctx.ShouldBindJSON(&user); err != nil {
		responses.GenerateBadRequestResponse(ctx, err.Error())

		return
	}

	hashedPassword, err := utils.HashPassword(user.Password)

	if err != nil {
		responses.GenerateBadRequestResponse(ctx, err.Error())
	}

	user.Password = hashedPassword
	creator := "admin " + username

	createdMember, err := controller.service.RegisterUserService(user, creator)

	if err != nil {
		responses.GenerateBadRequestResponse(ctx, err.Error())

		return
	}

	responses.GenerateSuccessResponseWithData(ctx, http.StatusCreated, fmt.Sprintf("admin create user with role (%s) success", user.Role), createdMember)
}

func (controller *adminController) GetAllUserController(ctx *gin.Context) {
	users, err := controller.service.GetAllUserService()

	if err != nil {
		responses.GenerateBadRequestResponse(ctx, err.Error())

		return
	}

	responses.GenerateSuccessResponseWithData(ctx, http.StatusOK, "admin get all user success", users)
}

func (controller *adminController) GetAllUserByRoleController(ctx *gin.Context) {
	role := ctx.Param("role")

	users, err := controller.service.GetAllUserByRoleService(role)

	if err != nil {
		responses.GenerateBadRequestResponse(ctx, err.Error())

		return
	}

	responses.GenerateSuccessResponseWithData(ctx, http.StatusOK, fmt.Sprintf("admin get all users with role \"%s\" success", role), users)
}

func (controller *adminController) GetUserByIdController(ctx *gin.Context) {
	id := ctx.Param("userId")

	member, err := controller.service.GetUserByIdService(id)

	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			responses.GenerateNotFoundResponse(ctx, err.Error())
		} else {
			responses.GenerateBadRequestResponse(ctx, err.Error())
		}

		return
	}

	responses.GenerateSuccessResponseWithData(ctx, http.StatusOK, fmt.Sprintf("admin get user by id \"%s\" success", id), member)
}

func (controller *adminController) UpdateUserByIdController(ctx *gin.Context) {
	_, username, _, err := middlewares.GetClaims(ctx)

	if err != nil {
		responses.GenerateUnauthorizedResponse(ctx, err.Error())

		return
	}

	var user users.UserDTO

	id := ctx.Param("id")

	if err := ctx.ShouldBindJSON(&user); err != nil {
		responses.GenerateNotFoundResponse(ctx, err.Error())

		return
	}

	user.Modified_By = fmt.Sprintf("admin %s", username)
	updatedMember, err := controller.service.UpdateUserByIdService(id, user)

	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			responses.GenerateNotFoundResponse(ctx, err.Error())
		} else {
			responses.GenerateBadRequestResponse(ctx, err.Error())
		}

		return
	}

	responses.GenerateSuccessResponseWithData(ctx, http.StatusOK, fmt.Sprintf("admin update user by id \"%s\" success", id), updatedMember)
}

func (controller *adminController) ModifyUserRoleByIdController(ctx *gin.Context) {
	id := ctx.Param("id")

	var user users.UserDTO

	if err := ctx.ShouldBindJSON(&user); err != nil {
		responses.GenerateBadRequestResponse(ctx, err.Error())

		return
	}

	modifiedMember, err := controller.service.ModifyUserRoleByIdService(id, user.Role)

	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			responses.GenerateNotFoundResponse(ctx, err.Error())
		} else {
			responses.GenerateBadRequestResponse(ctx, err.Error())
		}

		return
	}

	responses.GenerateSuccessResponseWithData(ctx, http.StatusOK, fmt.Sprintf("modifying user role by id \"%s\" success", id), modifiedMember)
}

func (controller *adminController) ModifyUserStatusByIdController(ctx *gin.Context) {
	id := ctx.Param("id")

	var user users.UserDTO

	if err := ctx.ShouldBindJSON(&user); err != nil {
		responses.GenerateBadRequestResponse(ctx, err.Error())

		return
	}

	modifiedMember, err := controller.service.ModifyUserStatusByIdService(id, user.Status)

	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			responses.GenerateNotFoundResponse(ctx, err.Error())
		} else {
			responses.GenerateBadRequestResponse(ctx, err.Error())
		}

		return
	}

	responses.GenerateSuccessResponseWithData(ctx, http.StatusOK, fmt.Sprintf("modifying user status by id \"%s\" success", id), modifiedMember)
}

func (controller *adminController) DeleteUserByIdController(ctx *gin.Context) {
	id := ctx.Param("id")

	deletedMember, err := controller.service.DeleteUserByIdService(id)

	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			responses.GenerateNotFoundResponse(ctx, err.Error())
		} else {
			responses.GenerateBadRequestResponse(ctx, err.Error())
		}

		return
	}

	responses.GenerateSuccessResponseWithData(ctx, http.StatusOK, fmt.Sprintf("delete member by id \"%s\" success", id), deletedMember)
}
