package users

import (
	"database/sql"
	"final-project/src/configs/database"
	"fmt"
)

type Repository interface {
	RegisterUserRepository(user RegisterUserDTO) (ViewUserDTO, error)
	ViewProfileRepository(id string) (ViewUserDTO, error)
	UpdateProfileRepository(id string, user UpdateUserDTO) (ViewUserDTO, error)
}

type userRepository struct{}

func NewRepository() Repository {
	return &userRepository{}
}

func (repository *userRepository) RegisterUserRepository(user RegisterUserDTO) (ViewUserDTO, error) {
	query := `
		INSERT INTO users
		(
			username,      
			password,    
			email,       
			first_name,       
			last_name,      
			address,       
			phone_number,      
			role_id,
			created_by,      
			modified_by     
		)
		VALUES
		($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
		RETURNING 
			id,
			username,      
			email,       
			first_name,       
			last_name,      
			address,       
			phone_number,
			is_penalized,
			penalty_duration,
			status,
			(SELECT name FROM roles WHERE roles.id = users.role_id) AS role;
	`

	var createdUser ViewUserDTO

	err := database.DB.QueryRow(query, user.Username, user.Password, user.Email, user.First_Name, user.Last_Name, user.Address, user.Phone_Number, user.Role_Id, user.Created_By, user.Modified_By).
		Scan(&createdUser.Id, &createdUser.Username, &createdUser.Email, &createdUser.First_Name, &createdUser.Last_Name, &createdUser.Address, &createdUser.Phone_Number, &createdUser.Is_Penalized, &createdUser.Penalty_Duration, &createdUser.Status, &createdUser.Role)

	if err != nil {
		return ViewUserDTO{}, err
	}

	return createdUser, err
}

func (repository *userRepository) ViewProfileRepository(id string) (ViewUserDTO, error) {
	var user ViewUserDTO

	query := `
		SELECT 
      users.id,
			users.username,      
			users.email,       
			users.first_name,       
			users.last_name,      
			users.address,       
			users.phone_number,
			users.is_penalized,
			users.penalty_duration,
			users.status,
      roles.name AS role
    FROM 
			users 
    LEFT JOIN 
      roles ON users.role_id = roles.id
		WHERE users.id = $1
	`

	err := database.DB.QueryRow(query, id).
		Scan(&user.Id, &user.Username, &user.Email, &user.First_Name, &user.Last_Name, &user.Address, &user.Phone_Number, &user.Is_Penalized, &user.Penalty_Duration, &user.Status, &user.Role)

	if err != nil {
		if err == sql.ErrNoRows {
			return ViewUserDTO{}, fmt.Errorf("failed to view profile, user with id \"%s\" not found", id)
		}

		return ViewUserDTO{}, err
	}

	return user, nil
}

func (repository *userRepository) UpdateProfileRepository(id string, user UpdateUserDTO) (ViewUserDTO, error) {
	query := `
		UPDATE users 
		SET 
			username = COALESCE(NULLIF($2, ''), username),      
			password = COALESCE(NULLIF($3, ''), password),
			email = COALESCE(NULLIF($4, ''), email),       
			first_name = COALESCE(NULLIF($5, ''), first_name),       
			last_name = COALESCE(NULLIF($6, ''), last_name),      
			address = COALESCE(NULLIF($7, ''), address),       
			phone_number = COALESCE(NULLIF($8, ''), phone_number),     
			modified_by = $9
		WHERE id = $1 
		RETURNING 
			id,
			username,      
			email,       
			first_name,       
			last_name,      
			address,       
			phone_number,
			is_penalized,
			penalty_duration,
			status,
			(SELECT name FROM roles WHERE roles.id = users.role_id) AS role;
	`

	var updatedUser ViewUserDTO

	err := database.DB.QueryRow(query, id, user.Username, user.Password, user.Email, user.First_Name, user.Last_Name, user.Address, user.Phone_Number, user.Modified_By).
		Scan(&updatedUser.Id, &updatedUser.Username, &updatedUser.Email, &updatedUser.First_Name, &updatedUser.Last_Name, &updatedUser.Address, &updatedUser.Phone_Number, &updatedUser.Is_Penalized, &updatedUser.Penalty_Duration, &updatedUser.Status, &updatedUser.Role)

	if err != nil {
		if err == sql.ErrNoRows {
			return ViewUserDTO{}, fmt.Errorf("failed updating profile, user with id \"%s\" not found", id)
		}

		return ViewUserDTO{}, err
	}

	return updatedUser, nil
}
