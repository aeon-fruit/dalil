package common

import "errors"

var (
	ErrNotFound    = errors.New("entity not found")
	ErrNotModified = errors.New("entity not modified")
)
