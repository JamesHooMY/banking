package user

import "errors"

var (
	ErrUserExisted         = errors.New("user existed")
	ErrUserNotFound        = errors.New("user not found")
)
