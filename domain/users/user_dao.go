package users

import (
	"database/sql"
	"fmt"
	"github.com/Komdosh/go-bookstore-users-api/datasources/postgresql/users_db"
	"github.com/Komdosh/go-bookstore-utils/logger"
	"github.com/Komdosh/go-bookstore-utils/rest_errors"
)

const (
	queryInsertUser                 = "INSERT INTO users_db.users(first_name, last_name, email, date_created, password, status) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id;"
	querySelectUserById             = "SELECT id, first_name, last_name, email, date_created, status FROM users_db.users WHERE id = $1;"
	queryUpdateUser                 = "UPDATE users_db.users SET first_name=$1, last_name=$2, email=$3 WHERE id=$4;"
	queryDeleteUser                 = "DELETE FROM users_db.users WHERE id=$1;"
	queryFindUserByStatus           = "SELECT id, first_name, last_name, email, date_created, status FROM users_db.users WHERE status = $1;"
	queryFindUserByEmailAndPassword = "SELECT id, first_name, last_name, email, date_created, status FROM users_db.users WHERE email = $1 AND password=$2 AND status = $3;"
)

func (user *User) Get() rest_errors.RestErr {
	stmt, err := users_db.Client.Prepare(querySelectUserById)
	if err != nil {
		logger.Error("error when trying to prepare get user statement", err)
		return rest_errors.NewInternalServerError("database error", err)
	}
	defer stmt.Close()

	if err := stmt.QueryRow(user.Id).Scan(&user.Id, &user.FirstName, &user.LastName, &user.Email, &user.DateCreated, &user.Status); err != nil {
		logger.Error("error when trying to get user by id", err)
		return rest_errors.NewInternalServerError("database error", err)
	}

	return nil
}

func (user *User) Save() rest_errors.RestErr {
	stmt, err := users_db.Client.Prepare(queryInsertUser)
	if err != nil {
		logger.Error("error when trying to prepare save user statement", err)
		return rest_errors.NewInternalServerError("database error", err)
	}
	defer stmt.Close()

	var userId int64

	if err := stmt.QueryRow(user.FirstName, user.LastName, user.Email, user.DateCreated, user.Password, user.Status).Scan(&userId); err != nil {
		logger.Error("error when trying to save user", err)
		return rest_errors.NewInternalServerError("database error", err)
	}

	user.Id = userId

	return nil
}

func (user *User) Update(newUser User) rest_errors.RestErr {
	stmt, err := users_db.Client.Prepare(queryUpdateUser)
	if err != nil {
		logger.Error("error when trying to prepare update user statement", err)
		return rest_errors.NewInternalServerError("database error", err)
	}
	defer stmt.Close()

	if _, err = stmt.Exec(user.FirstName, user.LastName, user.Email, user.Id); err != nil {
		logger.Error("error when trying to update user", err)
		return rest_errors.NewInternalServerError("database error", err)
	}

	return nil
}

func (user *User) Delete() rest_errors.RestErr {
	stmt, err := users_db.Client.Prepare(queryDeleteUser)
	if err != nil {
		logger.Error("error when trying to prepare delete user statement", err)
		return rest_errors.NewInternalServerError("database error", err)
	}
	defer stmt.Close()

	if _, err = stmt.Exec(user.Id); err != nil {
		logger.Error("error when trying to delete user", err)
		return rest_errors.NewInternalServerError("database error", err)
	}

	return nil
}

func (user *User) FindByStatus(status string) (Users, rest_errors.RestErr) {
	stmt, err := users_db.Client.Prepare(queryFindUserByStatus)
	if err != nil {
		logger.Error("error when trying to prepare find user by status statement", err)
		return nil, rest_errors.NewInternalServerError("database error", err)
	}
	defer stmt.Close()

	rows, err := stmt.Query(status)

	if err != nil {
		logger.Error("error when trying to find user by status", err)
		return nil, rest_errors.NewInternalServerError("database error", err)
	}
	defer rows.Close()

	results := make([]User, 0)

	for rows.Next() {
		var user User
		if err := rows.Scan(&user.Id, &user.FirstName, &user.LastName, &user.Email, &user.DateCreated, &user.Status); err != nil {
			logger.Error("error scan user by status", err)
			return nil, rest_errors.NewInternalServerError("database error", err)
		}
		results = append(results, user)
	}

	if len(results) == 0 {
		return nil, rest_errors.NewNotFoundError(fmt.Sprintf("no users matching status %s", status))
	}

	return results, nil
}

func (user *User) FindByEmailAndPassword() rest_errors.RestErr {
	stmt, err := users_db.Client.Prepare(queryFindUserByEmailAndPassword)
	if err != nil {
		logger.Error("error when trying to prepare find user by email and password statement", err)
		return rest_errors.NewInternalServerError("database error", err)
	}
	defer stmt.Close()

	if err := stmt.QueryRow(user.Email, user.Password, StatusActive).Scan(&user.Id, &user.FirstName, &user.LastName, &user.Email, &user.DateCreated, &user.Status); err != nil {
		if err == sql.ErrNoRows {
			return rest_errors.NewNotFoundError("invalid user credentials")
		}

		logger.Error("error when trying to get user by id", err)
		return rest_errors.NewInternalServerError("database error", err)
	}

	return nil
}
