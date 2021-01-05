package users

import (
	"fmt"
	"github.com/Komdosh/go-bookstore-users-api/datasources/postgresql/users_db"
	"github.com/Komdosh/go-bookstore-users-api/logger"
	"github.com/Komdosh/go-bookstore-users-api/utils/errors_utils"
)

const (
	queryInsertUser       = "INSERT INTO users_db.users(first_name, last_name, email, date_created, password, status) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id;"
	querySelectUserById   = "SELECT id, first_name, last_name, email, date_created, status FROM users_db.users WHERE id = $1;"
	queryUpdateUser       = "UPDATE users_db.users SET first_name=$1, last_name=$2, email=$3 WHERE id=$4;"
	queryDeleteUser       = "DELETE FROM users_db.users WHERE id=$1;"
	queryFindUserByStatus = "SELECT id, first_name, last_name, email, date_created, status FROM users_db.users WHERE status = $1;"
)

func (user *User) Get() *errors_utils.RestErr {
	stmt, err := users_db.Client.Prepare(querySelectUserById)
	if err != nil {
		logger.Error("error when trying to prepare get user statement", err)
		return errors_utils.NewInternalServerError("database error")
	}
	defer stmt.Close()

	if err := stmt.QueryRow(user.Id).Scan(&user.Id, &user.FirstName, &user.LastName, &user.Email, &user.DateCreated, &user.Status); err != nil {
		logger.Error("error when trying to get user by id", err)
		return errors_utils.NewInternalServerError("database error")
	}

	return nil
}

func (user *User) Save() *errors_utils.RestErr {
	stmt, err := users_db.Client.Prepare(queryInsertUser)
	if err != nil {
		logger.Error("error when trying to prepare save user statement", err)
		return errors_utils.NewInternalServerError("database error")
	}
	defer stmt.Close()

	var userId int64

	if err := stmt.QueryRow(user.FirstName, user.LastName, user.Email, user.DateCreated, user.Password, user.Status).Scan(&userId); err != nil {
		logger.Error("error when trying to save user", err)
		return errors_utils.NewInternalServerError("database error")
	}

	user.Id = userId

	return nil
}

func (user *User) Update(newUser User) *errors_utils.RestErr {
	stmt, err := users_db.Client.Prepare(queryUpdateUser)
	if err != nil {
		logger.Error("error when trying to prepare update user statement", err)
		return errors_utils.NewInternalServerError("database error")
	}
	defer stmt.Close()

	if _, err = stmt.Exec(user.FirstName, user.LastName, user.Email, user.Id); err != nil {
		logger.Error("error when trying to update user", err)
		return errors_utils.NewInternalServerError("database error")
	}

	return nil
}

func (user *User) Delete() *errors_utils.RestErr {
	stmt, err := users_db.Client.Prepare(queryDeleteUser)
	if err != nil {
		logger.Error("error when trying to prepare delete user statement", err)
		return errors_utils.NewInternalServerError("database error")
	}
	defer stmt.Close()

	if _, err = stmt.Exec(user.Id); err != nil {
		logger.Error("error when trying to delete user", err)
		return errors_utils.NewInternalServerError("database error")
	}

	return nil
}

func (user *User) FindByStatus(status string) (Users, *errors_utils.RestErr) {
	stmt, err := users_db.Client.Prepare(queryFindUserByStatus)
	if err != nil {
		logger.Error("error when trying to prepare find user by status statement", err)
		return nil, errors_utils.NewInternalServerError("database error")
	}
	defer stmt.Close()

	rows, err := stmt.Query(status)

	if err != nil {
		logger.Error("error when trying to find user by status", err)
		return nil, errors_utils.NewInternalServerError("database error")
	}
	defer rows.Close()

	results := make([]User, 0)

	for rows.Next() {
		var user User
		if err := rows.Scan(&user.Id, &user.FirstName, &user.LastName, &user.Email, &user.DateCreated, &user.Status); err != nil {
			logger.Error("error scan user by status", err)
			return nil, errors_utils.NewInternalServerError("database error")
		}
		results = append(results, user)
	}

	if len(results) == 0 {
		return nil, errors_utils.NewNotFoundError(fmt.Sprintf("no users matching status %s", status))
	}

	return results, nil
}
