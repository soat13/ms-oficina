package bun_helper

import (
	"database/sql"
	"errors"
	"strings"
)

var ErrResourceInUse = errors.New("cannot delete resource because it is in use")

func HandleDeleteError(err error) error {
	if err != nil && strings.Contains(err.Error(), "SQLSTATE 23503") {
		return ErrResourceInUse
	}

	return err
}

func IgnoreNoRows(err error) error {
	if errors.Is(err, sql.ErrNoRows) {
		return nil
	}

	return err
}
