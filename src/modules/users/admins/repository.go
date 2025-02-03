package admins

import (
	"database/sql"
	"final-project/src/configs/database"
	"final-project/src/modules/users"
	"fmt"
)

type Repository interface {
	GetAllUserRepository() ([]users.UserDTO, error)
	GetAllUserByRoleRepository(roleId string) ([]users.UserDTO, error)
	GetUserByIdRepository(userId string) (users.UserDTO, error)
	UpdateUserByIdRepository(userId string, user users.UserDTO) (users.UserDTO, error)
	ModifyUserRoleByIdRepository(userId string, roleId string) (users.UserDTO, error)
	ModifyUserStatusByIdRepository(userId string, status string) (users.UserDTO, error)
	DeleteUserByIdRepository(userId string) (users.UserDTO, error)
}

type adminRepository struct{}

func NewRepository() Repository {
	return &adminRepository{}
}

func (repository *adminRepository) GetAllUserRepository() ([]users.UserDTO, error) {
	var allUsers []users.UserDTO

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
      roles.name AS role,
			users.created_at,     
			users.created_by,     
			users.modified_at,
			users.modified_by
    FROM 
			users 
    LEFT JOIN 
      roles ON users.role_id = roles.id
	`

	rows, err := database.DB.Query(query)

	if err != nil {
		return []users.UserDTO{}, err
	}

	defer rows.Close()

	for rows.Next() {
		var user users.UserDTO

		err = rows.Scan(&user.Id, &user.Username, &user.Email, &user.First_Name, &user.Last_Name, &user.Address, &user.Phone_Number, &user.Is_Penalized, &user.Penalty_Duration, &user.Status, &user.Role, &user.Created_At, &user.Created_By, &user.Modified_At, &user.Modified_By)

		if err != nil {
			return []users.UserDTO{}, err
		}

		allUsers = append(allUsers, user)
	}

	return allUsers, nil
}

func (repository *adminRepository) GetAllUserByRoleRepository(roleId string) ([]users.UserDTO, error) {
	var allUsers []users.UserDTO

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
      roles.name AS role,
			users.created_at,     
			users.created_by,     
			users.modified_at,
			users.modified_by
    FROM 
			users 
    LEFT JOIN 
      roles ON users.role_id = roles.id
		WHERE
			users.role_id = $1
	`

	rows, err := database.DB.Query(query, roleId)

	if err != nil {
		return []users.UserDTO{}, err
	}

	defer rows.Close()

	for rows.Next() {
		var user users.UserDTO

		err = rows.Scan(&user.Id, &user.Username, &user.Email, &user.First_Name, &user.Last_Name, &user.Address, &user.Phone_Number, &user.Is_Penalized, &user.Penalty_Duration, &user.Status, &user.Role, &user.Created_At, &user.Created_By, &user.Modified_At, &user.Modified_By)

		if err != nil {
			return []users.UserDTO{}, err
		}

		allUsers = append(allUsers, user)
	}

	return allUsers, nil
}

func (repository *adminRepository) GetUserByIdRepository(userId string) (users.UserDTO, error) {
	var user users.UserDTO

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
      roles.name AS role,
			users.created_at,     
			users.created_by,     
			users.modified_at,
			users.modified_by
    FROM 
			users 
    LEFT JOIN 
      roles ON users.role_id = roles.id
		WHERE users.id = $1
	`

	err := database.DB.QueryRow(query, userId).
		Scan(&user.Id, &user.Username, &user.Email, &user.First_Name, &user.Last_Name, &user.Address, &user.Phone_Number, &user.Is_Penalized, &user.Penalty_Duration, &user.Status, &user.Role)

	if err != nil {
		if err == sql.ErrNoRows {
			return users.UserDTO{}, fmt.Errorf("failed to view profile, user with id \"%s\" not found", userId)
		}

		return users.UserDTO{}, err
	}

	return user, nil
}

func (repository *adminRepository) UpdateUserByIdRepository(userId string, user users.UserDTO) (users.UserDTO, error) {
	fmt.Println("userId :", userId)
	fmt.Println("user :", user)
	query := `
		UPDATE users 
		SET
			username = COALESCE(NULLIF($2, ''), username),      
			email = COALESCE(NULLIF($3, ''), email),  
			password = COALESCE(NULLIF($4, ''), password),
			first_name = COALESCE(NULLIF($5, ''), first_name),       
			last_name = COALESCE(NULLIF($6, ''), last_name),      
			address = COALESCE(NULLIF($7, ''), address),       
			phone_number = COALESCE(NULLIF($8, ''), phone_number),
			is_penalized = COALESCE($9, is_penalized),
			penalty_duration = COALESCE($10, penalty_duration),
			status = COALESCE(NULLIF($11, ''), status),
			role_id = CASE 
                WHEN $12 = '' THEN role_id 
                ELSE COALESCE($12::uuid, role_id)
              END,
			modified_by = COALESCE(NULLIF($13, ''), modified_by)
		WHERE id = $1
		RETURNING 
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
			(SELECT name FROM roles WHERE roles.id = users.role_id) AS role,
			users.created_at,     
			users.created_by,     
			users.modified_at,
			users.modified_by
	`

	err := database.DB.QueryRow(query, userId, user.Username, user.Email, user.Password, user.First_Name, user.Last_Name, user.Address, user.Phone_Number, user.Is_Penalized, user.Penalty_Duration, user.Status, user.Role_Id, user.Modified_By).
		Scan(&user.Id, &user.Username, &user.Email, &user.First_Name, &user.Last_Name, &user.Address, &user.Phone_Number, &user.Is_Penalized, &user.Penalty_Duration, &user.Status, &user.Role, &user.Created_At, &user.Created_By, &user.Modified_At, &user.Modified_By)

	if err != nil {
		if err == sql.ErrNoRows {
			return user, fmt.Errorf("failed updating user, user with id \"%s\" not found", userId)
		}

		return users.UserDTO{}, err
	}

	return user, nil
}

func (repository *adminRepository) ModifyUserRoleByIdRepository(userId string, roleId string) (users.UserDTO, error) {
	var modifiedUser users.UserDTO

	query := `
		UPDATE users 
		SET
			role_id = $2
		WHERE
			id = $1
		RETURNING 
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
			(SELECT name FROM roles WHERE roles.id = users.role_id) AS role,
			users.created_at,     
			users.created_by,     
			users.modified_at,
			users.modified_by
	`

	err := database.DB.QueryRow(query, userId, roleId).
		Scan(&modifiedUser.Id, &modifiedUser.Username, &modifiedUser.Email, &modifiedUser.First_Name, &modifiedUser.Last_Name, &modifiedUser.Address, &modifiedUser.Phone_Number, &modifiedUser.Is_Penalized, &modifiedUser.Penalty_Duration, &modifiedUser.Status, &modifiedUser.Role, &modifiedUser.Created_At, &modifiedUser.Created_By, &modifiedUser.Modified_At, &modifiedUser.Modified_By)

	if err != nil {
		if err == sql.ErrNoRows {
			return modifiedUser, fmt.Errorf("failed modifying user role, user with id \"%s\" not found", userId)
		}

		return users.UserDTO{}, err
	}

	return modifiedUser, nil
}

func (repository *adminRepository) ModifyUserStatusByIdRepository(userId string, status string) (users.UserDTO, error) {
	var modifiedUser users.UserDTO

	query := `
		UPDATE users 
		SET
			status = $2
		WHERE
			id = $1
		RETURNING
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
			(SELECT name FROM roles WHERE roles.id = users.role_id) AS role,
			users.created_at,     
			users.created_by,     
			users.modified_at,
			users.modified_by
	`

	err := database.DB.QueryRow(query, userId, status).
		Scan(&modifiedUser.Id, &modifiedUser.Username, &modifiedUser.Email, &modifiedUser.First_Name, &modifiedUser.Last_Name, &modifiedUser.Address, &modifiedUser.Phone_Number, &modifiedUser.Is_Penalized, &modifiedUser.Penalty_Duration, &modifiedUser.Status, &modifiedUser.Role, &modifiedUser.Created_At, &modifiedUser.Created_By, &modifiedUser.Modified_At, &modifiedUser.Modified_By)

	if err != nil {
		if err == sql.ErrNoRows {
			return modifiedUser, fmt.Errorf("failed modifying user status, user with id \"%s\" not found", userId)
		}

		return users.UserDTO{}, err
	}

	return modifiedUser, nil
}

func (repository *adminRepository) DeleteUserByIdRepository(userId string) (users.UserDTO, error) {
	var deletedUser users.UserDTO

	query := `
		DELETE FROM users 
		WHERE id = $1 
		RETURNING
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
			(SELECT name FROM roles WHERE roles.id = users.role_id) AS role,
			users.created_at,     
			users.created_by,     
			users.modified_at,
			users.modified_by
	`

	err := database.DB.QueryRow(query, userId).
		Scan(&deletedUser.Id, &deletedUser.Username, &deletedUser.Email, &deletedUser.First_Name, &deletedUser.Last_Name, &deletedUser.Address, &deletedUser.Phone_Number, &deletedUser.Is_Penalized, &deletedUser.Penalty_Duration, &deletedUser.Status, &deletedUser.Role, &deletedUser.Created_At, &deletedUser.Created_By, &deletedUser.Modified_At, &deletedUser.Modified_By)

	if err != nil {
		if err == sql.ErrNoRows {
			return deletedUser, fmt.Errorf("failed deleting user, user with id \"%s\" not found", userId)
		}

		return users.UserDTO{}, err
	}

	return deletedUser, nil
}
