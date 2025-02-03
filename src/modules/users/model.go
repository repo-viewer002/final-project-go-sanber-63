package users

import (
	"time"
)

type UserDTO struct {
	Id               string     `json:"id"`
	Username         string     `json:"username"`
	Password         string     `json:"password"`
	Email            string     `json:"email"`
	First_Name       string     `json:"first_name"`
	Last_Name        string     `json:"last_name"`
	Address          string     `json:"address"`
	Phone_Number     string     `json:"phone_number"`
	Is_Penalized     bool       `json:"is_penalized"`
	Penalty_Duration *time.Time `json:"penalty_duration"`
	Status           string     `json:"status"`
	Role_Id          string     `json:"role_id"`
	Role             string     `json:"role"`
	Created_At       time.Time  `json:"created_at"`
	Created_By       string     `json:"created_by"`
	Modified_At      time.Time  `json:"modified_at"`
	Modified_By      string     `json:"modified_by"`
}

type RegisterUserDTO struct {
	Username     string `json:"username"`
	Password     string `json:"password"`
	Email        string `json:"email"`
	First_Name   string `json:"first_name"`
	Last_Name    string `json:"last_name"`
	Address      string `json:"address"`
	Phone_Number string `json:"phone_number"`
	Role_Id      string `json:"role_id"`
	Role         string `json:"role"`
	Created_By   string `json:"created_by"`
	Modified_By  string `json:"modified_by"`
}

type ViewUserDTO struct {
	Id               string     `json:"id"`
	Username         string     `json:"username"`
	Email            string     `json:"email"`
	First_Name       string     `json:"first_name"`
	Last_Name        string     `json:"last_name"`
	Address          string     `json:"address"`
	Phone_Number     string     `json:"phone_number"`
	Is_Penalized     bool       `json:"is_penalized"`
	Penalty_Duration *time.Time `json:"penalty_duration,omitempty"`
	Status           string     `json:"status"`
	Role             string     `json:"role"`
}

type UpdateUserDTO struct {
	Username     string `json:"username"`
	Password     string `json:"password"`
	Email        string `json:"email"`
	First_Name   string `json:"first_name"`
	Last_Name    string `json:"last_name"`
	Address      string `json:"address"`
	Phone_Number string `json:"phone_number"`
	Role_Id      string `json:"role_id"`
	Created_By   string `json:"created_by"`
	Modified_By  string `json:"modified_by"`
}

type CreateUserRequestDTO struct {
	UserDTO
	Role string `json:"role"`
}

type UpdateUserRequestDTO struct {
	UserDTO
	Role   string `json:"role"`
	Status string `json:"status"`
}
