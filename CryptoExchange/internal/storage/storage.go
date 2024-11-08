package storage

import "errors"

var (
	ErrMaliciousParameter       = errors.New("malicious parameter")
	ErrInvalidSQLCommand        = errors.New("invalid sql command")
	ErrSQLExecFailed            = errors.New("sql command failed")
	ErrConnect                  = errors.New("database connection error")
	ErrIncorrectNumberOfColumns = errors.New("invalid number of columns")

	ErrUserExists = errors.New("user already exists")
)
