package users

import (
	"errors"
	"final-project/src/modules/roles"
	"final-project/src/utils"
)

type Service interface {
	RegisterUserService(user RegisterUserDTO, role string, creator string) (ViewUserDTO, error)
	ViewProfileService(userId string) (ViewUserDTO, error)
	UpdateProfileService(userId string, user UpdateUserDTO) (ViewUserDTO, error)
}

type userService struct {
	userRepository Repository
	roleRepository roles.Repository
}

func NewService(userRepository Repository, roleRepository roles.Repository) Service {
	return &userService{
		userRepository, 
		roleRepository,
	}
}

func (service *userService) RegisterUserService(user RegisterUserDTO, role string, creator string) (ViewUserDTO, error) {
	validRole := utils.IsValidRole(role)

	if !validRole {
		return ViewUserDTO{}, errors.New("invalid role")
	}

	roleId, err := service.roleRepository.GetRoleIdByNameRepository(role)

	if err != nil {
		return ViewUserDTO{}, err
	}

	user.Role_Id = roleId
	user.Created_By = creator
	user.Modified_By = user.Created_By
	registeredUser, err := service.userRepository.RegisterUserRepository(user)

	if err != nil {
		return ViewUserDTO{}, err
	}

	return registeredUser, nil
}

func (service *userService) ViewProfileService(userId string) (ViewUserDTO, error) {
	user, err := service.userRepository.ViewProfileRepository(userId)

	if err != nil {
		return ViewUserDTO{}, err
	}

	return user, nil
}

func (service *userService) UpdateProfileService(userId string, user UpdateUserDTO) (ViewUserDTO, error) {
	updatedUser, err := service.userRepository.UpdateProfileRepository(userId, user)

	if err != nil {
		return ViewUserDTO{}, err
	}

	return updatedUser, err
}
