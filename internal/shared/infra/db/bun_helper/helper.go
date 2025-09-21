package bun_helper

import (
	"database/sql"
	"errors"
)

func IgnoreNoRows(err error) error {
	if errors.Is(err, sql.ErrNoRows) {
		return nil
	}
	return err
}
