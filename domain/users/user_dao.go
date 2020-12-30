package users

import (
	"fmt"
	"github.com/Komdosh/go-bookstore-users-api/datasources/postgresql/users_db"
	"github.com/Komdosh/go-bookstore-users-api/utils/errors_utils"
	"github.com/Komdosh/go-bookstore-users-api/utils/postgres_utils"
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
		return errors_utils.NewInternalServerError(err.Error())
	}
	defer stmt.Close()

	if err := stmt.QueryRow(user.Id).Scan(&user.Id, &user.FirstName, &user.LastName, &user.Email, &user.DateCreated, &user.Status); err != nil {
		return postgres_utils.ParseErr(err)
	}

	return nil
}

func (user *User) Save() *errors_utils.RestErr {
	stmt, err := users_db.Client.Prepare(queryInsertUser)
	if err != nil {
		return errors_utils.NewInternalServerError(err.Error())
	}
	defer stmt.Close()

	var userId int64
	saveErr := stmt.QueryRow(user.FirstName, user.LastName, user.Email, user.DateCreated, user.Password, user.Status).Scan(&userId)

	if saveErr != nil {
		return postgres_utils.ParseErr(saveErr)
	}

	user.Id = userId

	return nil
}

func (user *User) Update(newUser User) *errors_utils.RestErr {
	stmt, err := users_db.Client.Prepare(queryUpdateUser)
	if err != nil {
		return errors_utils.NewInternalServerError(err.Error())
	}
	defer stmt.Close()

	_, err = stmt.Exec(user.FirstName, user.LastName, user.Email, user.Id)

	if err != nil {
		return postgres_utils.ParseErr(err)
	}

	return nil
}

func (user *User) Delete() *errors_utils.RestErr {
	stmt, err := users_db.Client.Prepare(queryDeleteUser)
	if err != nil {
		return errors_utils.NewInternalServerError(err.Error())
	}
	defer stmt.Close()

	_, err = stmt.Exec(user.Id)

	if err != nil {
		return postgres_utils.ParseErr(err)
	}

	return nil
}

func (user *User) FindByStatus(status string) (Users, *errors_utils.RestErr) {
	stmt, err := users_db.Client.Prepare(queryFindUserByStatus)
	if err != nil {
		return nil, errors_utils.NewInternalServerError(err.Error())
	}
	defer stmt.Close()

	rows, err := stmt.Query(status)

	if err != nil {
		return nil, errors_utils.NewInternalServerError(err.Error())
	}
	defer rows.Close()

	results := make([]User, 0)

	for rows.Next() {
		var user User
		if err := rows.Scan(&user.Id, &user.FirstName, &user.LastName, &user.Email, &user.DateCreated, &user.Status); err != nil {
			return nil, postgres_utils.ParseErr(err)
		}
		results = append(results, user)
	}

	if len(results) == 0 {
		return nil, errors_utils.NewNotFoundError(fmt.Sprintf("no users matching status %s", status))
	}

	return results, nil
}
