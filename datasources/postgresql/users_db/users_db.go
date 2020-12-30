package users_db

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"log"
	"os"
)

const (
	postgresUsersUsername = "postgres_users_username"
	postgresUsersPassword = "postgres_users_password"
	postgresUsersHost     = "postgres_users_host"
	postgresUsersPort     = "postgres_users_port"
	postgresUsersDatabase = "postgres_users_database"
)

var (
	Client *sql.DB

	username = os.Getenv(postgresUsersUsername)
	password = os.Getenv(postgresUsersPassword)
	host     = os.Getenv(postgresUsersHost)
	port     = os.Getenv(postgresUsersPort)
	database = os.Getenv(postgresUsersDatabase)
)

func init() {
	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, username, password, database)

	var err error
	Client, err = sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(nil)
	}

	if err := Client.Ping(); err != nil {
		panic(err)
	}

	log.Println("database successfully configured")
}
