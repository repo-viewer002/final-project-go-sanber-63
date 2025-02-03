package members

import (
	"final-project/src/commons/responses"
	"final-project/src/modules/users"
	"final-project/src/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Controller interface {
	RegisterMemberController(ctx *gin.Context)
}

type memberController struct {
	service Service
}

func NewController(service Service) Controller {
	return &memberController{
		service,
	}
}

func (controller *memberController) RegisterMemberController(ctx *gin.Context) {
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
	createdMember, err := controller.service.RegisterMemberService(member)

	if err != nil {
		responses.GenerateBadRequestResponse(ctx, err.Error())

		return
	}

	responses.GenerateSuccessResponseWithData(ctx, http.StatusCreated, "register member success", createdMember)
}
