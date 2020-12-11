package users

import (
	"github.com/Komdosh/go-bookstore-users-api/datasources/postgresql/users_db"
	"github.com/Komdosh/go-bookstore-users-api/utils/date"
	"github.com/Komdosh/go-bookstore-users-api/utils/errors"
	"github.com/Komdosh/go-bookstore-users-api/utils/postgres_utils"
)

const (
	queryInsertUser     = "INSERT INTO users_db.users(first_name, last_name, email, date_created) VALUES ($1, $2, $3, $4) RETURNING id;"
	querySelectUserById = "SELECT id, first_name, last_name, email, date_created FROM users_db.users WHERE id = $1;"
	queryUpdateUser     = "UPDATE users_db.users SET first_name=$1, last_name=$2, email=$3 WHERE id=$4;"
	queryDeleteUser     = "DELETE FROM users_db.users WHERE id=$1;"
)

func (user *User) Get() *errors.RestErr {
	stmt, err := users_db.Client.Prepare(querySelectUserById)
	if err != nil {
		return errors.NewInternalServerError(err.Error())
	}
	defer stmt.Close()

	if err := stmt.QueryRow(user.Id).Scan(&user.Id, &user.FirstName, &user.LastName, &user.Email, &user.DateCreated); err != nil {
		return postgres_utils.ParseErr(err)
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
	saveErr := stmt.QueryRow(user.FirstName, user.LastName, user.Email, user.DateCreated).Scan(&userId)

	if saveErr != nil {
		return postgres_utils.ParseErr(saveErr)
	}

	user.Id = userId

	return nil
}

func (user *User) Update(newUser User) *errors.RestErr {
	stmt, err := users_db.Client.Prepare(queryUpdateUser)
	if err != nil {
		return errors.NewInternalServerError(err.Error())
	}
	defer stmt.Close()

	_, err = stmt.Exec(user.FirstName, user.LastName, user.Email, user.Id)

	if err != nil {
		return postgres_utils.ParseErr(err)
	}

	return nil
}

func (user *User) Delete() *errors.RestErr {
	stmt, err := users_db.Client.Prepare(queryDeleteUser)
	if err != nil {
		return errors.NewInternalServerError(err.Error())
	}
	defer stmt.Close()

	_, err = stmt.Exec(user.Id)

	if err != nil {
		return postgres_utils.ParseErr(err)
	}

	return nil
}
