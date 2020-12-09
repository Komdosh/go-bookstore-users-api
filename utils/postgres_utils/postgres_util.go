package postgres_utils

import (
	"github.com/Komdosh/go-bookstore-users-api/utils/errors"
	"github.com/lib/pq"
	"strings"
)

const (
	errorNoRows = "no rows in result set"
)

func ParseErr(err error) *errors.RestErr {
	sqlErr, ok := err.(*pq.Error)
	if !ok {
		if strings.Contains(err.Error(), errorNoRows) {
			return errors.NewNotFoundError("no record matching given id")
		}

		return errors.NewInternalServerError(
			"error parsing database response",
		)
	}
	switch sqlErr.Code {
	case "23505":
		return errors.NewBadRequestError(sqlErr.Detail)
	}

	return errors.NewInternalServerError(sqlErr.Detail)
}
