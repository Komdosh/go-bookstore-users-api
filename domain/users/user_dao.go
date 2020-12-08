package users

import (
	"fmt"
	"github.com/Komdosh/go-bookstore-users-api/datasources/postgresql/users_db"
	"github.com/Komdosh/go-bookstore-users-api/utils/date"
	"github.com/Komdosh/go-bookstore-users-api/utils/errors"
	"strings"
)

const (
	queryInsertUser = "INSERT INTO users_db.users(first_name, last_name, email, date_created) VALUES ($1, $2, $3, $4)  RETURNING id;"
)

var (
	usersDB = make(map[int64]*User)
)

func (user *User) Get() *errors.RestErr {

	if err := users_db.Client.Ping(); err != nil {
		panic(err)
	}

	result := usersDB[user.Id]
	if result == nil {
		return errors.NewNotFoundError(fmt.Sprintf("user %d not found", user.Id))
	}

	user.Id = result.Id
	user.FirstName = result.FirstName
	user.LastName = result.LastName
	user.Email = result.Email
	user.DateCreated = result.DateCreated

	return nil
}

func (user *User) Save() *errors.RestErr {
	stmt, err := users_db.Client.Prepare(queryInsertUser)
	if err != nil {
		return errors.NewInternalServerError(err.Error())
	}
	defer stmt.Close()

	user.DateCreated = date.GetNowString()

	var userId int64
	err = stmt.QueryRow(user.FirstName, user.LastName, user.Email, user.DateCreated).Scan(&userId)
	if err != nil {
		if strings.Contains(err.Error(), "uindex") {
			return errors.NewBadRequestError(err.Error())
		}
		return errors.NewInternalServerError(
			fmt.Sprintf("error when trying to save user: %s", err.Error()),
		)
	}

	user.Id = userId

	return nil
}
