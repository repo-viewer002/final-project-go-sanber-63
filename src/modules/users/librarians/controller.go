package librarians

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
	CreateMemberController(ctx *gin.Context)
	GetAllMemberController(ctx *gin.Context)
	GetMemberByIdController(ctx *gin.Context)
	UpdateMemberByIdController(ctx *gin.Context)
}

type memberController struct {
	service Service
}

func NewController(service Service) Controller {
	return &memberController{
		service,
	}
}

func (controller *memberController) CreateMemberController(ctx *gin.Context) {
	_, username, _, err := middlewares.GetClaims(ctx)

	if err != nil {
		responses.GenerateUnauthorizedResponse(ctx, err.Error())

		return
	}

	var member users.RegisterUserDTO

	if err := ctx.ShouldBindJSON(&member); err != nil {
		responses.GenerateBadRequestResponse(ctx, err.Error())

		return
	}

	hashedPassword, err := utils.HashPassword(member.Password)

	if err != nil {
		responses.GenerateBadRequestResponse(ctx, err.Error())
	}

	member.Password = hashedPassword
	createdMember, err := controller.service.CreateMemberService(member, username)

	if err != nil {
		responses.GenerateBadRequestResponse(ctx, err.Error())

		return
	}

	responses.GenerateSuccessResponseWithData(ctx, http.StatusCreated, "librarian create member success", createdMember)
}

func (controller *memberController) GetAllMemberController(ctx *gin.Context) {
	members, err := controller.service.GetAllMemberService()

	if err != nil {
		responses.GenerateBadRequestResponse(ctx, err.Error())

		return
	}

	responses.GenerateSuccessResponseWithData(ctx, http.StatusOK, "librarian get all member success", members)
}

func (controller *memberController) GetMemberByIdController(ctx *gin.Context) {
	getId := ctx.Param("memberId")

	member, err := controller.service.GetMemberByIdService(getId)

	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			responses.GenerateNotFoundResponse(ctx, err.Error())
		} else {
			responses.GenerateBadRequestResponse(ctx, err.Error())
		}

		return
	}

	responses.GenerateSuccessResponseWithData(ctx, http.StatusOK, fmt.Sprintf("librarian get member by id \"%s\" success", getId), member)
}

func (controller *memberController) UpdateMemberByIdController(ctx *gin.Context) {
	_, username, _, err := middlewares.GetClaims(ctx)

	if err != nil {
		responses.GenerateUnauthorizedResponse(ctx, err.Error())

		return
	}

	var member users.UpdateUserDTO

	getId := ctx.Param("memberId")

	if err := ctx.ShouldBindJSON(&member); err != nil {
		responses.GenerateNotFoundResponse(ctx, err.Error())

		return
	}

	member.Modified_By = fmt.Sprintf("librarian %s", username)
	updatedMember, err := controller.service.UpdateMemberByIdService(getId, member)

	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			responses.GenerateNotFoundResponse(ctx, err.Error())
		} else {
			responses.GenerateBadRequestResponse(ctx, err.Error())
		}

		return
	}

	responses.GenerateSuccessResponseWithData(ctx, http.StatusOK, fmt.Sprintf("librarian update member by id \"%s\" success", getId), updatedMember)
}
