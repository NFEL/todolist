package api_error

import (
	"errors"
	"fmt"
)

var (
	ErrInvalidCredentials = errors.New("invalid username or password")
	ErrUserNotFound       = errors.New("user not found")
	ErrTaskNotFound       = errors.New("task not found")
	ErrUnauthorized       = errors.New("unauthorized")
	ErrTokenExpired       = errors.New("token expired")
	ErrTokenRevoked       = errors.New("token has been revoked")
	ErrInvalidToken       = errors.New("invalid token")
)

func UsernameExists(s string) error {
	return fmt.Errorf("username %s already exists", s)
}
