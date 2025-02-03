package genres

import (
	"final-project/src/commons/middlewares"
	"final-project/src/commons/responses"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type Controller interface {
	CreateGenreController(ctx *gin.Context)
	GetAllGenreController(ctx *gin.Context)
	GetGenreByIdController(ctx *gin.Context)
	UpdateGenreByIdController(ctx *gin.Context)
	DeleteGenreByIdController(ctx *gin.Context)
}

type genreController struct {
	service Service
}

func NewController(service Service) Controller {
	return &genreController{
		service,
	}
}

func (controller *genreController) CreateGenreController(ctx *gin.Context) {
	_, username, _, err := middlewares.GetClaims(ctx)

	if err != nil {
		responses.GenerateUnauthorizedResponse(ctx, err.Error())

		return
	}

	var genre Genre

	if err := ctx.ShouldBindJSON(&genre); err != nil {
		responses.GenerateBadRequestResponse(ctx, err.Error())

		return
	}

	genre.Created_By = username
	genre.Modified_By = username

	createdGenre, err := controller.service.CreateGenreService(genre)

	if err != nil {
		responses.GenerateBadRequestResponse(ctx, err.Error())

		return
	}

	responses.GenerateSuccessResponseWithData(ctx, http.StatusCreated, "create genre success", createdGenre)
}

func (controller *genreController) GetAllGenreController(ctx *gin.Context) {
	name := ctx.Query("name")

	genre, err := controller.service.GetAllGenreService(name)

	if err != nil {
		responses.GenerateBadRequestResponse(ctx, err.Error())

		return
	}

	responses.GenerateSuccessResponseWithData(ctx, http.StatusOK, "get all genre success", genre)
}

func (controller *genreController) GetGenreByIdController(ctx *gin.Context) {
	getId := ctx.Param("id")

	genre, err := controller.service.GetGenreByIdService(getId)

	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			responses.GenerateNotFoundResponse(ctx, err.Error())
		} else {
			responses.GenerateBadRequestResponse(ctx, err.Error())
		}

		return
	}

	responses.GenerateSuccessResponseWithData(ctx, http.StatusOK, fmt.Sprintf("get genre by id \"%s\" success", getId), genre)
}

func (controller *genreController) UpdateGenreByIdController(ctx *gin.Context) {
	_, username, _, err := middlewares.GetClaims(ctx)

	if err != nil {
		responses.GenerateUnauthorizedResponse(ctx, err.Error())

		return
	}

	var genre Genre

	getId := ctx.Param("id")

	if err := ctx.ShouldBindJSON(&genre); err != nil {
		responses.GenerateNotFoundResponse(ctx, err.Error())

		return
	}

	genre.Modified_By = username
	updatedGenre, err := controller.service.UpdateGenreByIdService(getId, genre)

	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			responses.GenerateNotFoundResponse(ctx, err.Error())
		} else {
			responses.GenerateBadRequestResponse(ctx, err.Error())
		}

		return
	}

	responses.GenerateSuccessResponseWithData(ctx, http.StatusOK, fmt.Sprintf("update genre by id \"%s\" success", getId), updatedGenre)
}

func (controller *genreController) DeleteGenreByIdController(ctx *gin.Context) {
	getId := ctx.Param("id")

	deletedGenre, err := controller.service.DeleteGenreByIdService(getId)

	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			responses.GenerateNotFoundResponse(ctx, err.Error())
		} else {
			responses.GenerateBadRequestResponse(ctx, err.Error())
		}

		return
	}

	responses.GenerateSuccessResponseWithData(ctx, http.StatusOK, fmt.Sprintf("delete genre by id \"%s\" success", getId), deletedGenre)
}
