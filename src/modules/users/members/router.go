package members

import (
	"final-project/src/modules/roles"
	"final-project/src/modules/users"

	"github.com/gin-gonic/gin"
)

func MemberRouter(router *gin.Engine) {
	roleRepository := roles.NewRepository()
	userRepository := users.NewRepository()
	userService := users.NewService(userRepository, roleRepository)

	memberService := NewService(userService)
	memberController := NewController(memberService)

	api := router.Group("/api")
	api.POST("/register", memberController.RegisterMemberController)
}
