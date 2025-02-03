package main

import (
	"final-project/src/commons"
	"final-project/src/commons/middlewares"
	"final-project/src/configs/database"
	"final-project/src/modules/auth"
	"final-project/src/modules/books"
	"final-project/src/modules/borrows"
	"final-project/src/modules/genres"
	"final-project/src/modules/roles"
	"final-project/src/modules/users"
	"final-project/src/modules/users/admins"
	"final-project/src/modules/users/librarians"
	"final-project/src/modules/users/members"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	database.InitializeDB()

	router := gin.Default()
	router.Use(middlewares.Log())

	router.GET("/", indexController)

	roles.RoleRouter(router)

	auth.AuthRouter(router)

	users.UserRouter(router)
	members.MemberRouter(router)
	librarians.LibrarianRouter(router)
	admins.AdminRouter(router)

	genres.GenreRouter(router)
	books.BookRouter(router)
	borrows.BorrowRouter(router)

	router.Run(fmt.Sprintf(":%d", commons.PORT))
}

func indexController(ctx *gin.Context) {
	scheme := "http"
	if ctx.Request.TLS != nil {
		scheme = "https"
	}

	host := ctx.Request.Host

	apiDocumentation := scheme + "://" + host + "/swagger/index.html"

	ctx.Data(http.StatusOK, "text/html", []byte(`
			<!DOCTYPE html>
			<html lang="en">
			<head>
				<meta charset="UTF-8">
				<meta name="viewport" content="width=device-width, initial-scale=1.0">
				<title>API Documentation</title>
			</head>
			<body>
				<h1>Sanbercodes Golang Backend Development Batch 63 | Quiz 3</h1>
				<a href="`+apiDocumentation+`" target="_blank">API Documentation</a></br>
				<a href="`+commons.REPOSITORY+`" target="_blank">Github Repository</a></br>
				<a href="`+commons.ENDPOINT+`" target="_blank">Railway Deployment URL</a>
			</body>
			</html>
		`))
}
