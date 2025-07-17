package db

import "errors"

var (
	ErrDBConnection = errors.New("database connection error")

	ErrRecordNotFound = errors.New("record not found")
	ErrDBQuery        = errors.New("database query error")
	ErrDuplicateEntry = errors.New("duplicate entry")
	ErrTimeout        = errors.New("database operation timeout")
)
