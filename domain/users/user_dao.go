package users

import (
	"fmt"
	"github.com/Komdosh/go-bookstore-users-api/datasources/postgresql/users_db"
	"github.com/Komdosh/go-bookstore-users-api/utils/date"
	"github.com/Komdosh/go-bookstore-users-api/utils/errors"
	"strings"
)

const (
	queryInsertUser     = "INSERT INTO users_db.users(first_name, last_name, email, date_created) VALUES ($1, $2, $3, $4) RETURNING id;"
	querySelectUserById = "SELECT id, first_name, last_name, email, date_created FROM users_db.users WHERE id = $1;"
	errorNoRows         = "no rows in result set"
	errorNotUnique      = "uindex"
)

var (
	usersDB = make(map[int64]*User)
)

func (user *User) Get() *errors.RestErr {
	stmt, err := users_db.Client.Prepare(querySelectUserById)
	if err != nil {
		return errors.NewInternalServerError(err.Error())
	}
	defer stmt.Close()

	if err := stmt.QueryRow(user.Id).Scan(&user.Id, &user.FirstName, &user.LastName, &user.Email, &user.DateCreated); err != nil {

		if strings.Contains(err.Error(), errorNoRows) {
			return errors.NewNotFoundError(fmt.Sprintf("user %d not found", user.Id))
		}
		return errors.NewInternalServerError(fmt.Sprintf("error when trying to get user with id: %d", user.Id))
	}

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
		if strings.Contains(err.Error(), errorNotUnique) {
			return errors.NewBadRequestError(err.Error())
		}
		return errors.NewInternalServerError(
			fmt.Sprintf("error when trying to save user: %s", err.Error()),
		)
	}

	user.Id = userId

	return nil
}
