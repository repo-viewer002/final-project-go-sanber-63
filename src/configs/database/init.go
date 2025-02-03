package database

import (
	"database/sql"
	"final-project/src/commons"
	"fmt"

	_ "github.com/lib/pq"
)

var (
	DB  *sql.DB
	err error
)

func InitializeDB() {
	fmt.Printf("Database Information\nDatabase : Postgres\nHost : %s\nPort : %s\nUser : %s\nDatabase Name : %s\n", commons.DB_HOST, commons.DB_PORT, commons.DB_USER, commons.DB_NAME)
	fmt.Println("\nConnecting to database...")

	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s", commons.DB_HOST, commons.DB_PORT, commons.DB_USER, commons.DB_PASSWORD, commons.DB_NAME, commons.DB_SSL_MODE)

	DB, err = sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}

	err = DB.Ping()
	if err != nil {
		panic(err)
	}

	DBMigrate(DB)

	fmt.Println("\nSuccessfully connected to database")
}
