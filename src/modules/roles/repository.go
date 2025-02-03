package roles

import (
	"database/sql"
	"final-project/src/configs/database"
	"fmt"
)

type Repository interface {
	CreateRoleRepository(role Role) (Role, error)
	GetAllRoleRepository() ([]Role, error)
	GetRoleByIdRepository(id string) (Role, error)
	GetRoleIdByNameRepository(name string) (string, error)
	UpdateRoleByIdRepository(id string, role Role) (Role, error)
	DeleteRoleByIdRepository(id string) (Role, error)
}

type roleRepository struct{}

func NewRepository() Repository {
	return &roleRepository{}
}

func (repository *roleRepository) CreateRoleRepository(role Role) (Role, error) {
	query := `
		INSERT INTO roles
		(
			name,
			description,
			created_by, 
			modified_by
		)
		VALUES
		($1, $2, $3, $4)
		RETURNING *
	`

	err := database.DB.QueryRow(query, role.Name, role.Description, role.Created_By, role.Modified_By).
		Scan(&role.Id, &role.Name, &role.Description, &role.Created_At, &role.Created_By, &role.Modified_At, &role.Modified_By)

	if err != nil {
		return Role{}, err
	}

	return role, err
}

func (repository *roleRepository) GetAllRoleRepository() ([]Role, error) {
	var roles []Role

	query := "SELECT * FROM roles"

	rows, err := database.DB.Query(query)

	if err != nil {
		return []Role{}, err
	}

	defer rows.Close()

	for rows.Next() {
		var role Role

		err = rows.Scan(&role.Id, &role.Name, &role.Description, &role.Created_At, &role.Created_By, &role.Modified_At, &role.Modified_By)

		if err != nil {
			return []Role{}, err
		}

		roles = append(roles, role)
	}

	return roles, nil
}

func (repository *roleRepository) GetRoleByIdRepository(id string) (Role, error) {
	var role Role

	query := `
		SELECT * FROM roles 
		WHERE id = $1
	`

	err := database.DB.QueryRow(query, id).
		Scan(&role.Id, &role.Name, &role.Description, &role.Created_At, &role.Created_By, &role.Modified_At, &role.Modified_By)

	if err != nil {
		if err == sql.ErrNoRows {
			return Role{}, fmt.Errorf("failed to get role data, role with id \"%s\" not found", id)
		}

		return Role{}, err
	}

	return role, nil
}

func (repository *roleRepository) GetRoleIdByNameRepository(name string) (string, error) {
	var role Role

	query := `
		SELECT id FROM roles 
		WHERE name = $1
	`

	err := database.DB.QueryRow(query, name).
		Scan(&role.Id)

	if err != nil {
		if err == sql.ErrNoRows {
			return "", fmt.Errorf("failed to get role data, role with name \"%s\" not found", name)
		}

		return "", err
	}

	return role.Id, nil
}

func (repository *roleRepository) UpdateRoleByIdRepository(id string, role Role) (Role, error) {
	query := `
		UPDATE roles 
		SET 
			name = $2,
			description = $3,
			modified_by = $4 
		WHERE id = $1 
		RETURNING *
	`

	err := database.DB.QueryRow(query, id, role.Name, role.Description, role.Modified_By).
		Scan(&role.Id, &role.Name, &role.Description, &role.Created_At, &role.Created_By, &role.Modified_At, &role.Modified_By)

	if err != nil {
		if err == sql.ErrNoRows {
			return role, fmt.Errorf("failed updating role, role with id \"%s\" not found", id)
		}

		return Role{}, err
	}

	return role, nil
}

func (repository *roleRepository) DeleteRoleByIdRepository(id string) (Role, error) {
	var deletedRole Role

	query := `
		DELETE FROM roles 
		WHERE id = $1 
		RETURNING *
	`

	err := database.DB.QueryRow(query, id).
		Scan(&deletedRole.Id, &deletedRole.Name, &deletedRole.Description, &deletedRole.Created_At, &deletedRole.Created_By, &deletedRole.Modified_At, &deletedRole.Modified_By)

	if err != nil {
		if err == sql.ErrNoRows {
			return deletedRole, fmt.Errorf("failed deleting role, role with id \"%s\" not found", id)
		}

		return Role{}, err
	}

	return deletedRole, nil
}
