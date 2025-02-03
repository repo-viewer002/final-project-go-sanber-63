package commons

import (
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type RoleName struct {
	Admin     string
	Librarian string
	Member    string
}

type UserStatuses struct {
	Active      string
	Deactivated string
	Suspended   string
}

var (
	PORT                   int
	DB_HOST                string
	DB_PORT                string
	DB_USER                string
	DB_PASSWORD            string
	DB_NAME                string
	DB_SSL_MODE            string
	ENDPOINT               string
	REPOSITORY             string
	JWT_SECRET_KEY         []byte
	PENALTY_AMOUNT_PER_DAY int
)

var Roles = RoleName{
	Admin:     "admin",
	Librarian: "librarian",
	Member:    "member",
}

var UserStatus = UserStatuses{
	Active:      "active",
	Deactivated: "deactivated",
	Suspended:   "suspended",
}

func init() {
	if os.Getenv("RAILWAY_ENVIRONMENT_NAME") == "" {
		err := godotenv.Load("../.env")

		if err != nil {
			panic("Error loading .env file, " + err.Error())
		}
	}

	var getPort = os.Getenv("PORT")

	portNumber, err := strconv.Atoi(getPort)
	if err != nil {
		panic("Invalid PORT value (int expected) : " + getPort)
	}

	var getJwtSecretKey = os.Getenv("JWT_SECRET_KEY")
	var jwtSecretKey = []byte(getJwtSecretKey)

	PORT = portNumber
	DB_HOST = os.Getenv("DB_HOST")
	DB_PORT = os.Getenv("DB_PORT")
	DB_USER = os.Getenv("DB_USER")
	DB_PASSWORD = os.Getenv("DB_PASSWORD")
	DB_NAME = os.Getenv("DB_NAME")
	DB_SSL_MODE = os.Getenv("DB_SSL_MODE")
	ENDPOINT = os.Getenv("ENDPOINT")
	REPOSITORY = os.Getenv("REPOSITORY")
	JWT_SECRET_KEY = jwtSecretKey
	PENALTY_AMOUNT_PER_DAY = 1000
}
