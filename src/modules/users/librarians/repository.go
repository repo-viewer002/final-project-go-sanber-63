package librarians

import (
	"final-project/src/configs/database"
	"final-project/src/modules/users"
)

type Repository interface {
	GetAllMemberRepository(memberRoleId string) ([]users.ViewUserDTO, error)
}

type memberRepository struct{}

func NewRepository() Repository {
	return &memberRepository{}
}

func (repository *memberRepository) GetAllMemberRepository(memberRoleId string) ([]users.ViewUserDTO, error) {
	var members []users.ViewUserDTO

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
		WHERE users.role_id = $1
	`

	rows, err := database.DB.Query(query, memberRoleId)

	if err != nil {
		return []users.ViewUserDTO{}, err
	}

	defer rows.Close()

	for rows.Next() {
		var member users.ViewUserDTO

		err = rows.Scan(&member.Id, &member.Username, &member.Email, &member.First_Name, &member.Last_Name, &member.Address, &member.Phone_Number, &member.Is_Penalized, &member.Penalty_Duration, &member.Status, &member.Role)

		if err != nil {
			return []users.ViewUserDTO{}, err
		}

		members = append(members, member)
	}

	return members, nil
}
