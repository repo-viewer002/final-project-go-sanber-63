package librarians

import (
	"final-project/src/commons"
	"final-project/src/modules/roles"
	"final-project/src/modules/users"
)

type Service interface {
	CreateMemberService(member users.RegisterUserDTO, librarianUsername string) (users.ViewUserDTO, error)
	GetAllMemberService() ([]users.ViewUserDTO, error)
	GetMemberByIdService(memberId string) (users.ViewUserDTO, error)
	UpdateMemberByIdService(memberId string, user users.UpdateUserDTO) (users.ViewUserDTO, error)
}

type memberService struct {
	repository  Repository
	userService users.Service
	roleService roles.Service
}

func NewService(repository Repository, userService users.Service, roleService roles.Service) Service {
	return &memberService{
		repository,
		userService,
		roleService,
	}
}

func (service *memberService) CreateMemberService(member users.RegisterUserDTO, librarianUsername string) (users.ViewUserDTO, error) {
	memberRoleId, err := service.roleService.GetRoleIdByNameRepository(commons.Roles.Member)

	if err != nil {
		return users.ViewUserDTO{}, err
	}

	creator := commons.Roles.Librarian + " " + librarianUsername
	member.Role_Id = memberRoleId
	createdMember, err := service.userService.RegisterUserService(member, commons.Roles.Member, creator)

	if err != nil {
		return users.ViewUserDTO{}, err
	}

	return createdMember, nil
}

func (service *memberService) GetAllMemberService() ([]users.ViewUserDTO, error) {
	memberRoleId, err := service.roleService.GetRoleIdByNameRepository(commons.Roles.Member)

	if err != nil {
		return []users.ViewUserDTO{}, err
	}

	members, err := service.repository.GetAllMemberRepository(memberRoleId)

	if err != nil {
		return []users.ViewUserDTO{}, err
	}

	return members, nil
}

func (service *memberService) GetMemberByIdService(memberId string) (users.ViewUserDTO, error) {
	user, err := service.userService.ViewProfileService(memberId)

	if err != nil {
		return users.ViewUserDTO{}, err
	}

	return user, nil
}

func (service *memberService) UpdateMemberByIdService(memberId string, user users.UpdateUserDTO) (users.ViewUserDTO, error) {
	updatedUser, err := service.userService.UpdateProfileService(memberId, user)

	if err != nil {
		return users.ViewUserDTO{}, err
	}

	return updatedUser, err
}
