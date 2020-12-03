package services

import (
	"github.com/Komdosh/go-bookstore-users-api/domain/users"
	"github.com/Komdosh/go-bookstore-users-api/utils/errors"
)

func CreateUser(user users.User) (*users.User, *errors.RestErr) {
	return &user, nil
}

func GetUser() {

}

func FindUser() {

}
