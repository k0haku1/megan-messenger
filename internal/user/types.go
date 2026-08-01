package user

import (
	"errors"
)

var (
	ErrInvalidImage = errors.New("invalid image")
	ErrUploadAvatar = errors.New("failed to upload avatar")
)

type MeResponse struct {
	ID                 string `json:"id"`
	Phone              string `json:"phone"`
	Username           string `json:"username,omitempty"`
	AvatarURL          string `json:"avatarUrl,omitempty"`
	HasPassword        bool   `json:"hasPassword"`
	OnboardingComplete bool   `json:"onboardingComplete"`
}
