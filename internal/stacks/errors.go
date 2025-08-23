package stacks

import (
	"errors"

	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
)

var (
	ErrNotFound        = errors.New("stack not found")
	ErrAlreadyArchived = errors.New("stack already archived")
	ErrConflict        = errors.New("slug already exists")
	ErrInvalidInput    = errors.New("invalid input")
)

func isUniqueViolation(err error, constraint string) bool {
	var pgerr *pgconn.PgError
	if !errors.As(err, &pgerr) {
		return false
	}
	if pgerr.Code != pgerrcode.UniqueViolation {
		return false
	}
	if constraint == "" {
		return true
	}
	return pgerr.ConstraintName == constraint
}
