package user

import (
	"errors"
)

var (
	ErrInvalidImage        = errors.New("invalid image")
	ErrUploadAvatar        = errors.New("failed to upload avatar")
	ErrUserNotDiscoverable = errors.New("user not found")
	ErrCannotMessageUser   = errors.New("user does not accept messages")
	ErrUsernameTaken       = errors.New("username already taken")
)

const (
	searchMinQueryLength = 2
	searchDefaultLimit   = int32(20)
	searchMaxLimit       = int32(20)
)
