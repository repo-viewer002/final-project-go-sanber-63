package auth

import (
	"database/sql"
	"errors"
	"final-project/src/configs/database"
)

type Repository interface {
	ValidateUsernameAndEmail(identifier string) (ValidUser, error)
}

type authRepository struct{}

func NewRepository() Repository {
	return &authRepository{}
}

func (repository *authRepository) ValidateUsernameAndEmail(identifier string) (ValidUser, error) {
	var user ValidUser

	query := `
		SELECT
			users.id,
			users.username,
			users.email,
			users.password,
			roles.name AS role 
		FROM 
			users
		LEFT JOIN 
			roles ON users.role_id = roles.id
		WHERE 
			username = $1
		OR
			email = $1
	`

	err := database.DB.QueryRow(query, identifier).
		Scan(&user.Id, &user.Username, &user.Email, &user.Password, &user.Role)

	if err != nil {
		if err == sql.ErrNoRows {
			return ValidUser{}, errors.New("invalid credentials")
		}
		return ValidUser{}, err
	}

	return user, err
}
