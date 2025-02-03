package admins

import (
	"errors"
	"final-project/src/modules/roles"
	"final-project/src/modules/users"
	"final-project/src/utils"
)

type Service interface {
	RegisterUserService(user users.RegisterUserDTO, creator string) (users.ViewUserDTO, error)
	GetAllUserService() ([]users.UserDTO, error)
	GetAllUserByRoleService(role string) ([]users.UserDTO, error)
	GetUserByIdService(userId string) (users.UserDTO, error)
	UpdateUserByIdService(userId string, user users.UserDTO) (users.UserDTO, error)
	ModifyUserStatusByIdService(userId string, status string) (users.UserDTO, error)
	ModifyUserRoleByIdService(userId string, role string) (users.UserDTO, error)
	DeleteUserByIdService(userId string) (users.UserDTO, error)
}

type adminService struct {
	adminRepository Repository
	roleRepository  roles.Repository
	userService     users.Service
}

func NewService(adminRepository Repository, roleRepository roles.Repository, userService users.Service) Service {
	return &adminService{
		adminRepository,
		roleRepository,
		userService,
	}
}

func (service *adminService) RegisterUserService(user users.RegisterUserDTO, creator string) (users.ViewUserDTO, error) {
	registeredMember, err := service.userService.RegisterUserService(user, user.Role, creator)

	if err != nil {
		return users.ViewUserDTO{}, err
	}

	return registeredMember, nil
}

func (service *adminService) GetAllUserService() ([]users.UserDTO, error) {
	allUsers, err := service.adminRepository.GetAllUserRepository()

	if err != nil {
		return []users.UserDTO{}, err
	}

	return allUsers, nil
}

func (service *adminService) GetAllUserByRoleService(role string) ([]users.UserDTO, error) {
	validRole := utils.IsValidRole(role)

	if !validRole {
		return []users.UserDTO{}, errors.New("invalid role")
	}

	roleId, err := service.roleRepository.GetRoleIdByNameRepository(role)

	if err != nil {
		return []users.UserDTO{}, err
	}

	allUsersByRole, err := service.adminRepository.GetAllUserByRoleRepository(roleId)

	if err != nil {
		return []users.UserDTO{}, err
	}

	return allUsersByRole, nil
}

func (service *adminService) GetUserByIdService(userId string) (users.UserDTO, error) {
	user, err := service.adminRepository.GetUserByIdRepository(userId)

	if err != nil {
		return users.UserDTO{}, err
	}

	return user, nil
}

func (service *adminService) UpdateUserByIdService(userId string, user users.UserDTO) (users.UserDTO, error) {
	updatedUser, err := service.adminRepository.UpdateUserByIdRepository(userId, user)

	if err != nil {
		return users.UserDTO{}, err
	}

	return updatedUser, err
}

func (service *adminService) ModifyUserRoleByIdService(userId string, role string) (users.UserDTO, error) {
	validRole := utils.IsValidRole(role)

	if !validRole {
		return users.UserDTO{}, errors.New("invalid role")
	}

	roleId, err := service.roleRepository.GetRoleIdByNameRepository(role)

	if err != nil {
		return users.UserDTO{}, err
	}

	modifiedUser, err := service.adminRepository.ModifyUserRoleByIdRepository(userId, roleId)

	if err != nil {
		return users.UserDTO{}, err
	}

	return modifiedUser, err
}

func (service *adminService) ModifyUserStatusByIdService(userId string, status string) (users.UserDTO, error) {
	validStatus := utils.IsValidStatus(status)

	if !validStatus {
		return users.UserDTO{}, errors.New("invalid status")
	}
	
	modifiedUser, err := service.adminRepository.ModifyUserStatusByIdRepository(userId, status)

	if err != nil {
		return users.UserDTO{}, err
	}

	return modifiedUser, err
}

func (service *adminService) DeleteUserByIdService(userId string) (users.UserDTO, error) {
	deletedUser, err := service.adminRepository.DeleteUserByIdRepository(userId)

	if err != nil {
		return users.UserDTO{}, err
	}

	return deletedUser, err
}
