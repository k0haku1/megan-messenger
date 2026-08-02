package user

import (
	"megan-messenger/internal/model"
)

type MeResponse struct {
	ID                 string `json:"id"`
	Phone              string `json:"phone"`
	Username           string `json:"username,omitempty"`
	AvatarURL          string `json:"avatarUrl,omitempty"`
	HasPassword        bool   `json:"hasPassword"`
	OnboardingComplete bool   `json:"onboardingComplete"`
	UsernameSearchable bool   `json:"usernameSearchable"`
	DMPolicy           string `json:"dmPolicy"`
}

type SearchUsersResponse struct {
	Users []model.UserSearchResult `json:"users"`
}

type UpdateUsernameRequest struct {
	Username string `json:"username" validate:"required,min=5,max=32"`
}

type UpdateUsernameResponse struct {
	Username string `json:"username"`
}

type UpdatePrivacyRequest struct {
	UsernameSearchable bool   `json:"usernameSearchable"`
	DMPolicy           string `json:"dmPolicy" validate:"required,oneof=everyone nobody"`
}

type UpdatePrivacyResponse struct {
	UsernameSearchable bool   `json:"usernameSearchable"`
	DMPolicy           string `json:"dmPolicy"`
}
