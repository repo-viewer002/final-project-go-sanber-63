package borrows

import (
	"final-project/src/commons/middlewares"
	"final-project/src/commons/responses"
	"final-project/src/utils"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Controller interface {
	BorrowBookController(ctx *gin.Context)
	ReturnBookController(ctx *gin.Context)
}

type borrowController struct {
	service Service
}

func NewController(service Service) Controller {
	return &borrowController{
		service,
	}
}

func (controller *borrowController) BorrowBookController(ctx *gin.Context) {
	_, username, role, err := middlewares.GetClaims(ctx)
	
	fmt.Println(username, role)
	if err != nil {
		responses.GenerateUnauthorizedResponse(ctx, err.Error())
		return
	}

	var borrow Borrow
	if err := ctx.ShouldBindJSON(&borrow); err != nil {
		responses.GenerateBadRequestResponse(ctx, err.Error())
		return
	}

	if len(borrow.Books) < 1 {
		responses.GenerateBadRequestResponse(ctx, "please input book ids to borrow")
		return
	}

	utils.GenerateDataModifier(role, username, &borrow.Created_By)

	createdBook, err := controller.service.BorrowBookService(borrow)
	if err != nil {
		responses.GenerateBadRequestResponse(ctx, err.Error())
		return
	}

	responses.GenerateSuccessResponseWithData(ctx, http.StatusCreated, "borrow books success", createdBook)
}

func (controller *borrowController) ReturnBookController(ctx *gin.Context) {
	borrowId := ctx.Param("borrowId")

	createdBook, err := controller.service.ReturnBookService(borrowId)
	if err != nil {
		responses.GenerateBadRequestResponse(ctx, err.Error())
		return
	}

	responses.GenerateSuccessResponseWithData(ctx, http.StatusCreated, "return books success", createdBook)
}
