package user

import (
	"errors"
	"regexp"
	"strings"
)

var (
	ErrInvalidUsername = errors.New("username must be 5-32 chars, start with a letter, and contain only a-z, 0-9, or _")
	usernamePattern    = regexp.MustCompile(`^[a-z][a-z0-9_]{4,31}$`)
)

func NormalizeUsername(raw string) string {
	value := strings.TrimSpace(raw)
	value = strings.TrimPrefix(value, "@")
	return strings.ToLower(value)
}

func ValidateUsername(username string) error {
	if !usernamePattern.MatchString(username) {
		return ErrInvalidUsername
	}
	return nil
}
