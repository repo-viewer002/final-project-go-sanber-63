package roles

type Service interface {
	CreateRoleService(role Role) (Role, error)
	GetAllRoleService() ([]Role, error)
	GetRoleByIdService(roleId string) (Role, error)
	GetRoleIdByNameRepository(name string) (string, error)
	UpdateRoleByIdService(roleId string, role Role) (Role, error)
	DeleteRoleByIdService(roleId string) (Role, error)
}

type roleService struct {
	repository Repository
}

func NewService(repository Repository) Service {
	return &roleService{
		repository,
	}
}

func (service *roleService) CreateRoleService(role Role) (Role, error) {
	createdRole, err := service.repository.CreateRoleRepository(role)

	if err != nil {
		return Role{}, err
	}

	return createdRole, nil
}

func (service *roleService) GetAllRoleService() ([]Role, error) {
	role, err := service.repository.GetAllRoleRepository()

	if err != nil {
		return []Role{}, err
	}

	return role, nil
}

func (service *roleService) GetRoleByIdService(roleId string) (Role, error) {
	role, err := service.repository.GetRoleByIdRepository(roleId)

	if err != nil {
		return Role{}, err
	}

	return role, nil
}

func (service *roleService) GetRoleIdByNameRepository(name string) (string, error) {
	role, err := service.repository.GetRoleIdByNameRepository(name)

	if err != nil {
		return "", err
	}

	return role, nil
}

func (service *roleService) UpdateRoleByIdService(roleId string, role Role) (Role, error) {
	updatedRole, err := service.repository.UpdateRoleByIdRepository(roleId, role)

	if err != nil {
		return Role{}, err
	}

	return updatedRole, err
}

func (service *roleService) DeleteRoleByIdService(roleId string) (Role, error) {
	deletedRole, err := service.repository.DeleteRoleByIdRepository(roleId)

	if err != nil {
		return Role{}, err
	}

	return deletedRole, err
}
