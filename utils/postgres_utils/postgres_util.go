package postgres_utils

import (
	"github.com/Komdosh/go-bookstore-users-api/utils/errors_utils"
	"github.com/lib/pq"
	"strings"
)

const (
	ErrorNoRows = "no rows in result set"
)

func ParseErr(err error) *errors_utils.RestErr {
	sqlErr, ok := err.(*pq.Error)
	if !ok {
		if strings.Contains(err.Error(), ErrorNoRows) {
			return errors_utils.NewNotFoundError("no record matching given id")
		}

		return errors_utils.NewInternalServerError(
			"error parsing database response",
		)
	}
	switch sqlErr.Code {
	case "23505":
		return errors_utils.NewBadRequestError(sqlErr.Detail)
	}

	return errors_utils.NewInternalServerError(sqlErr.Detail)
}
