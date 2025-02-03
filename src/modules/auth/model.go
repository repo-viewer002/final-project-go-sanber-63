package auth

type Credentials struct {
	Identifier string `json:"identifier"`
	Password   string `json:"password"`
}

type ValidUser struct {
	Id       string `json:"id"`
	Username string `json:"username"`
	Email    string `json:"enmail"`
	Password string `json:"password"`
	Role     string `json:"role"`
}
