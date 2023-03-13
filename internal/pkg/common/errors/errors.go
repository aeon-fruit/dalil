package errors

import "errors"

var (
	ErrNotFound    = errors.New("value not found")
	ErrNotModified = errors.New("value not modified")
)
