package books

import (
	"final-project/src/commons/middlewares"
	"final-project/src/commons/responses"
	"final-project/src/utils"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type Controller interface {
	CreateBookController(ctx *gin.Context)
	GetAllBookController(ctx *gin.Context)
	GetAllBookByGenreController(ctx *gin.Context)
	GetBookByIdController(ctx *gin.Context)
	UpdateBookByIdController(ctx *gin.Context)
	DeleteBookByIdController(ctx *gin.Context)
}

type bookController struct {
	service Service
}

func NewController(service Service) Controller {
	return &bookController{
		service,
	}
}

func (controller *bookController) CreateBookController(ctx *gin.Context) {
	_, username, role, err := middlewares.GetClaims(ctx)
	if err != nil {
		responses.GenerateUnauthorizedResponse(ctx, err.Error())
		return
	}

	var book Book
	if err := ctx.ShouldBindJSON(&book); err != nil {
		responses.GenerateBadRequestResponse(ctx, err.Error())
		return
	}

	utils.GenerateDataModifier(role, username, &book.Created_By)
	utils.GenerateDataModifier(role, username, &book.Modified_By)

	createdBook, err := controller.service.CreateBookService(book)
	if err != nil {
		responses.GenerateBadRequestResponse(ctx, err.Error())
		return
	}

	responses.GenerateSuccessResponseWithData(ctx, http.StatusCreated, "create book success", createdBook)
}

func (controller *bookController) GetAllBookController(ctx *gin.Context) {
	genreSearchTypeQuery := ctx.Query("genre_search_type")
	genresQuery := ctx.Query("genres")

	if genreSearchTypeQuery != "" && genreSearchTypeQuery != "any" && genreSearchTypeQuery != "all" {
		responses.GenerateBadRequestResponse(ctx, "Invalid search type. Use 'any' or 'all'")
		return
	}

	genres := strings.Split(genresQuery, ",")

	var searchBook = SearchBook{
		Name:              ctx.DefaultQuery("name", ""),
		Authors:           ctx.DefaultQuery("authors", ""),
		Publisher:         ctx.DefaultQuery("publisher", ""),
		Publish_Year:      ctx.DefaultQuery("publish_year", ""),
		Genre_Search_Type: genreSearchTypeQuery,
		Genres:            genres,
	}

	book, err := controller.service.GetAllBookService(searchBook)

	if err != nil {
		responses.GenerateBadRequestResponse(ctx, err.Error())

		return
	}

	responses.GenerateSuccessResponseWithData(ctx, http.StatusOK, "get all book success", book)
}

func (controller *bookController) GetAllBookByGenreController(ctx *gin.Context) {
	searchTypeQuery := ctx.Query("condition")
	genresQuery := ctx.Query("genres")

	if searchTypeQuery != "any" && searchTypeQuery != "all" {
		responses.GenerateBadRequestResponse(ctx, "Invalid search type. Use 'any' or 'all'")
		return
	}

	genres := strings.Split(genresQuery, ",")

	book, err := controller.service.GetAllBookByGenreService(searchTypeQuery, genres...)

	if err != nil {
		responses.GenerateBadRequestResponse(ctx, err.Error())

		return
	}

	switch searchTypeQuery {
	case "any":
		responses.GenerateSuccessResponseWithData(ctx, http.StatusOK, "get all book by matching any genre of "+genresQuery+" success", book)
	case "all":
		responses.GenerateSuccessResponseWithData(ctx, http.StatusOK, "get all book by matching all genre of genre "+genresQuery+" success", book)
	default:
		responses.GenerateBadRequestResponse(ctx, "invalid search condition")
	}
}

func (controller *bookController) GetBookByIdController(ctx *gin.Context) {
	getId := ctx.Param("bookId")

	book, err := controller.service.GetBookByIdService(getId)

	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			responses.GenerateNotFoundResponse(ctx, err.Error())
		} else {
			responses.GenerateBadRequestResponse(ctx, err.Error())
		}

		return
	}

	responses.GenerateSuccessResponseWithData(ctx, http.StatusOK, fmt.Sprintf("get book by id \"%s\" success", getId), book)
}

func (controller *bookController) UpdateBookByIdController(ctx *gin.Context) {
	_, username, role, err := middlewares.GetClaims(ctx)

	if err != nil {
		responses.GenerateUnauthorizedResponse(ctx, err.Error())

		return
	}

	var book Book

	getId := ctx.Param("bookId")

	if err := ctx.ShouldBindJSON(&book); err != nil {
		responses.GenerateNotFoundResponse(ctx, err.Error())

		return
	}

	utils.GenerateDataModifier(role, username, &book.Modified_By)
	updatedBook, err := controller.service.UpdateBookByIdService(getId, book)

	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			responses.GenerateNotFoundResponse(ctx, err.Error())
		} else {
			responses.GenerateBadRequestResponse(ctx, err.Error())
		}

		return
	}

	responses.GenerateSuccessResponseWithData(ctx, http.StatusOK, fmt.Sprintf("update book by id \"%s\" success", getId), updatedBook)
}

func (controller *bookController) DeleteBookByIdController(ctx *gin.Context) {
	getId := ctx.Param("bookId")

	deletedBook, err := controller.service.DeleteBookByIdService(getId)

	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			responses.GenerateNotFoundResponse(ctx, err.Error())
		} else {
			responses.GenerateBadRequestResponse(ctx, err.Error())
		}

		return
	}

	responses.GenerateSuccessResponseWithData(ctx, http.StatusOK, fmt.Sprintf("delete book by id \"%s\" success", getId), deletedBook)
}
