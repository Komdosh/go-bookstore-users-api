package services

import (
	"github.com/Komdosh/go-bookstore-users-api/domain/users"
	"github.com/Komdosh/go-bookstore-users-api/utils/crypto_utils"
	"github.com/Komdosh/go-bookstore-users-api/utils/date_utils"
	"github.com/Komdosh/go-bookstore-users-api/utils/errors_utils"
)

var (
	UsersService usersServiceInterface = &usersService{}
)

type usersService struct {
}

type usersServiceInterface interface {
	CreateUser(user users.User) (*users.User, *errors_utils.RestErr)
	GetUser(userId int64) (*users.User, *errors_utils.RestErr)
	UpdateUser(isPartial bool, user users.User) (*users.User, *errors_utils.RestErr)
	DeleteUser(userId int64) *errors_utils.RestErr
	SearchUser(status string) (users.Users, *errors_utils.RestErr)
}

func (s *usersService) CreateUser(user users.User) (*users.User, *errors_utils.RestErr) {
	if err := user.Validate(); err != nil {
		return nil, err
	}

	user.DateCreated = date_utils.GetNowDbFormat()
	user.Status = users.StatusActive
	user.Password = crypto_utils.GetSha1(user.Password)

	if err := user.Save(); err != nil {
		return nil, err
	}

	return &user, nil
}

func (s *usersService) GetUser(userId int64) (*users.User, *errors_utils.RestErr) {
	result := &users.User{Id: userId}

	if err := result.Get(); err != nil {
		return nil, err
	}

	return result, nil
}

func (s *usersService) UpdateUser(isPartial bool, user users.User) (*users.User, *errors_utils.RestErr) {
	current, err := s.GetUser(user.Id)

	if err != nil {
		return nil, err
	}

	if isPartial {
		if user.FirstName != "" {
			current.FirstName = user.FirstName
		}
		if user.LastName != "" {
			current.LastName = user.LastName
		}
		if user.Email != "" {
			current.Email = user.Email
		}
	} else {
		current.FirstName = user.FirstName
		current.LastName = user.LastName
		current.Email = user.Email
	}
	if err := current.Update(user); err != nil {
		return nil, err
	}

	return current, nil
}

func (s *usersService) DeleteUser(userId int64) *errors_utils.RestErr {
	user := &users.User{Id: userId}

	return user.Delete()
}

func (s *usersService) SearchUser(status string) (users.Users, *errors_utils.RestErr) {
	dao := &users.User{}
	return dao.FindByStatus(status)
}
