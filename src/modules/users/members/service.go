package members

import (
	"final-project/src/commons"
	"final-project/src/modules/users"
)

type Service interface {
	RegisterMemberService(member users.RegisterUserDTO) (users.ViewUserDTO, error)
}

type memberService struct {
	userService users.Service
}

func NewService(userService users.Service) Service {
	return &memberService{
		userService,
	}
}

func (service *memberService) RegisterMemberService(member users.RegisterUserDTO) (users.ViewUserDTO, error) {
	registeredMember, err := service.userService.RegisterUserService(member, commons.Roles.Member, "system")

	if err != nil {
		return users.ViewUserDTO{}, err
	}

	return registeredMember, nil
}